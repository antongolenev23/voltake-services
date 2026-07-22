package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

func (p *Postgres) GetStations(
	ctx context.Context,
	limit, offset int,
) ([]domain.ChargingStation, error) {
	const op = "postgres.GetStations"

	const query = `
		SELECT
			id,
			name,
			address,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,
			created_at
		FROM charging_stations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := p.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	stations := make([]domain.ChargingStation, 0)

	for rows.Next() {
		var station domain.ChargingStation

		if err := rows.Scan(
			&station.ID,
			&station.Name,
			&station.Address,
			&station.Latitude,
			&station.Longitude,
			&station.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		stations = append(stations, station)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return stations, nil
}

func (p *Postgres) GetNearbyStations(
	ctx context.Context,
	lat, lng, radius float64,
	limit, offset int,
) ([]domain.ChargingStation, error) {
	const op = "postgres.GetNearbyStations"

	const query = `
		SELECT
			id,
			name,
			address,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,
			created_at
		FROM charging_stations
		WHERE ST_DWithin(
			location,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
			$3 * 1000
		)
		ORDER BY location <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
		LIMIT $4 OFFSET $5
	`

	rows, err := p.db.Query(ctx, query, lng, lat, radius, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	stations := make([]domain.ChargingStation, 0)

	for rows.Next() {
		var station domain.ChargingStation

		if err := rows.Scan(
			&station.ID,
			&station.Name,
			&station.Address,
			&station.Latitude,
			&station.Longitude,
			&station.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		stations = append(stations, station)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return stations, nil
}

func (p *Postgres) GetStation(
	ctx context.Context,
	id uuid.UUID,
) (domain.ChargingStation, error) {
	const op = "postgres.GetStation"

	const query = `
		SELECT
			id,
			name,
			address,
			latitude,
			longitude,
			created_at
		FROM charging_stations
		WHERE id = $1
	`

	var station domain.ChargingStation

	err := p.db.QueryRow(ctx, query, id).Scan(
		&station.ID,
		&station.Name,
		&station.Address,
		&station.Latitude,
		&station.Longitude,
		&station.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ChargingStation{}, fmt.Errorf("%s: %w", op, domain.ErrStationNotFound)
		}

		return domain.ChargingStation{}, fmt.Errorf("%s: %w", op, err)
	}

	return station, nil
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
