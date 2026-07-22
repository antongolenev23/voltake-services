package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

func (p *Postgres) GetStations(
	ctx context.Context,
	filter domain.StationFilter,
) ([]domain.ChargingStation, error) {
	const op = "postgres.GetStations"

	query, args := buildStationsQuery(filter)

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	stations := make([]domain.ChargingStation, 0)

	for rows.Next() {
		var station domain.ChargingStation

		err := rows.Scan(
			&station.ID,
			&station.Name,
			&station.Address,
			&station.Latitude,
			&station.Longitude,
			&station.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}

		stations = append(stations, station)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows: %w", op, err)
	}

	return stations, nil
}

func buildStationsQuery(
	filter domain.StationFilter,
) (string, []any) {
	var query strings.Builder

	args := make([]any, 0)
	where := make([]string, 0)

	arg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}

	addPortFilters(&where, arg, filter)
	addGeoFilter(&where, arg, filter)

	order := "cs.created_at DESC"

	if filter.Geo != nil {
		order = fmt.Sprintf(`
			cs.location <-> ST_SetSRID(
				ST_MakePoint(%s,%s),
				4326
			)::geography
		`,
			arg(filter.Geo.Lng),
			arg(filter.Geo.Lat),
		)
	}

	query.WriteString(`
		SELECT
			cs.id,
			cs.name,
			cs.address,
			ST_Y(cs.location::geometry),
			ST_X(cs.location::geometry),
			cs.created_at
		FROM charging_stations cs
	`)

	if len(where) > 0 {
		query.WriteString("\nWHERE ")
		query.WriteString(strings.Join(where, " AND "))
	}

	query.WriteString("\nORDER BY ")
	query.WriteString(order)

	query.WriteString(fmt.Sprintf(
		" LIMIT %s OFFSET %s",
		arg(filter.Limit),
		arg(filter.Offset),
	))

	return query.String(), args
}

func addPortFilters(
	where *[]string,
	arg func(any) string,
	filter domain.StationFilter,
) {
	if filter.ConnectorType == nil && filter.MinPowerKW == nil {
		return
	}

	conditions := make([]string, 0, 2)

	if filter.ConnectorType != nil {
		conditions = append(
			conditions,
			fmt.Sprintf(
				"cp.connector_type = %s",
				arg(*filter.ConnectorType),
			),
		)
	}

	if filter.MinPowerKW != nil {
		conditions = append(
			conditions,
			fmt.Sprintf(
				"cp.power_kw >= %s",
				arg(*filter.MinPowerKW),
			),
		)
	}

	*where = append(*where, fmt.Sprintf(`
		EXISTS (
			SELECT 1
			FROM charging_ports cp
			WHERE cp.station_id = cs.id
			AND %s
		)
	`, strings.Join(conditions, " AND ")))
}

func addGeoFilter(
	where *[]string,
	arg func(any) string,
	filter domain.StationFilter,
) {
	if filter.Geo == nil {
		return
	}

	*where = append(*where, fmt.Sprintf(`
		ST_DWithin(
			cs.location,
			ST_SetSRID(
				ST_MakePoint(%s,%s),
				4326
			)::geography,
			%s * 1000
		)
	`,
		arg(filter.Geo.Lng),
		arg(filter.Geo.Lat),
		arg(filter.Geo.Radius),
	))
}

func (p *Postgres) GetStation(
	ctx context.Context,
	id uuid.UUID,
) (domain.ChargingStationDetails, error) {
	const op = "postgres.GetStation"

	const stationQuery = `
		SELECT
			id,
			name,
			address,
			ST_Y(location::geometry),
			ST_X(location::geometry),
			created_at
		FROM charging_stations
		WHERE id = $1
	`

	var station domain.ChargingStation

	err := p.db.QueryRow(ctx, stationQuery, id).Scan(
		&station.ID,
		&station.Name,
		&station.Address,
		&station.Latitude,
		&station.Longitude,
		&station.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ChargingStationDetails{}, fmt.Errorf(
				"%s: %w",
				op,
				domain.ErrStationNotFound,
			)
		}

		return domain.ChargingStationDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	const portsQuery = `
		SELECT
			id,
			station_id,
			connector_type,
			power_kw,
			is_active,
			created_at
		FROM charging_ports
		WHERE station_id = $1
		ORDER BY created_at
	`

	rows, err := p.db.Query(ctx, portsQuery, id)
	if err != nil {
		return domain.ChargingStationDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	ports := make([]domain.ChargingPort, 0)

	for rows.Next() {
		var port domain.ChargingPort

		err := rows.Scan(
			&port.ID,
			&port.StationID,
			&port.ConnectorType,
			&port.PowerKW,
			&port.IsActive,
			&port.CreatedAt,
		)

		if err != nil {
			return domain.ChargingStationDetails{}, fmt.Errorf("%s: %w", op, err)
		}

		ports = append(ports, port)
	}

	if err := rows.Err(); err != nil {
		return domain.ChargingStationDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	return domain.ChargingStationDetails{
		ChargingStation: station,
		Ports:           ports,
	}, nil
}

func (p *Postgres) CreateStation(
	ctx context.Context,
	station domain.ChargingStation,
) (domain.ChargingStation, error) {
	const op = "postgres.CreateStation"

	const query = `
		INSERT INTO charging_stations (
			owner_id,
			name,
			address,
			latitude,
			longitude
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			owner_id,
			name,
			address,
			latitude,
			longitude,
			created_at
	`

	var createdStation domain.ChargingStation

	err := p.db.QueryRow(
		ctx,
		query,
		station.OwnerID,
		station.Name,
		station.Address,
		station.Latitude,
		station.Longitude,
	).Scan(
		&createdStation.ID,
		&createdStation.OwnerID,
		&createdStation.Name,
		&createdStation.Address,
		&createdStation.Latitude,
		&createdStation.Longitude,
		&createdStation.CreatedAt,
	)

	if err != nil {
		return domain.ChargingStation{}, fmt.Errorf(
			"%s: %w",
			op,
			err,
		)
	}

	return createdStation, nil
}

func (p *Postgres) UpdateStation(
	ctx context.Context,
	station domain.ChargingStation,
) (domain.ChargingStation, error) {
	const op = "postgres.UpdateStation"

	const query = `
		UPDATE charging_stations
		SET
			name = $1,
			address = $2,
			latitude = $3,
			longitude = $4
		WHERE id = $5
		AND owner_id = $6
		RETURNING
			id,
			owner_id,
			name,
			address,
			latitude,
			longitude,
			created_at
	`

	var updatedStation domain.ChargingStation

	err := p.db.QueryRow(
		ctx,
		query,
		station.Name,
		station.Address,
		station.Latitude,
		station.Longitude,
		station.ID,
		station.OwnerID,
	).Scan(
		&updatedStation.ID,
		&updatedStation.OwnerID,
		&updatedStation.Name,
		&updatedStation.Address,
		&updatedStation.Latitude,
		&updatedStation.Longitude,
		&updatedStation.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ChargingStation{}, fmt.Errorf(
				"%s: %w",
				op,
				domain.ErrStationNotFound,
			)
		}

		return domain.ChargingStation{}, fmt.Errorf("%s: %w", op, err)
	}

	return updatedStation, nil
}

func (p *Postgres) DeleteStation(
	ctx context.Context,
	id uuid.UUID,
	ownerID uuid.UUID,
) error {
	const op = "postgres.DeleteStation"

	const query = `
		DELETE FROM charging_stations
		WHERE id = $1
		AND owner_id = $2
	`

	result, err := p.db.Exec(ctx, query, id, ownerID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, domain.ErrStationNotFound)
	}

	return nil
}
