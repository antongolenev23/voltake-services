package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

func (p *Postgres) GetPorts(
	ctx context.Context,
	stationID uuid.UUID,
) ([]domain.ChargingPort, error) {
	const op = "postgres.GetPorts"

	const query = `
		SELECT
			cp.id,
			cp.station_id,
			cp.connector_type,
			cp.power_kw,
			cp.is_active,
			cp.created_at
		FROM charging_stations cs
		LEFT JOIN charging_ports cp
			ON cp.station_id = cs.id
		WHERE cs.id = $1
		ORDER BY cp.created_at;
	`

	rows, err := p.db.Query(ctx, query, stationID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	ports := make([]domain.ChargingPort, 0)

	foundStation := false

	for rows.Next() {
		foundStation = true

		var (
			portID        *uuid.UUID
			portStationID *uuid.UUID
			connectorType *string
			powerKW       *int
			isActive      *bool
			createdAt     *time.Time
		)

		err := rows.Scan(
			&portID,
			&portStationID,
			&connectorType,
			&powerKW,
			&isActive,
			&createdAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		// станция есть, но портов нет
		if portID == nil {
			continue
		}

		ports = append(ports, domain.ChargingPort{
			ID:            *portID,
			StationID:     *portStationID,
			ConnectorType: *connectorType,
			PowerKW:       *powerKW,
			IsActive:      *isActive,
			CreatedAt:     *createdAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !foundStation {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrStationNotFound)
	}

	return ports, nil
}

func (p *Postgres) GetPort(
	ctx context.Context,
	stationID uuid.UUID,
	portID uuid.UUID,
) (domain.ChargingPort, error) {
	const op = "postgres.GetPort"

	const query = `
		SELECT
			id,
			station_id,
			connector_type,
			power_kw,
			is_active,
			created_at
		FROM charging_ports
		WHERE id = $1
		AND station_id = $2
	`

	var port domain.ChargingPort

	err := p.db.QueryRow(
		ctx,
		query,
		portID,
		stationID,
	).Scan(
		&port.ID,
		&port.StationID,
		&port.ConnectorType,
		&port.PowerKW,
		&port.IsActive,
		&port.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ChargingPort{}, fmt.Errorf(
				"%s: %w",
				op,
				domain.ErrPortNotFound,
			)
		}

		return domain.ChargingPort{}, fmt.Errorf("%s: %w", op, err)
	}

	return port, nil
}

func (p *Postgres) CreatePort(
	ctx context.Context,
	port domain.ChargingPort,
) (domain.ChargingPort, error) {
	const op = "postgres.CreatePort"

	const query = `
		INSERT INTO charging_ports (
			station_id,
			connector_type,
			power_kw,
			is_active
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			station_id,
			connector_type,
			power_kw,
			is_active,
			created_at
	`

	var createdPort domain.ChargingPort

	err := p.db.QueryRow(
		ctx,
		query,
		port.StationID,
		port.ConnectorType,
		port.PowerKW,
		port.IsActive,
	).Scan(
		&createdPort.ID,
		&createdPort.StationID,
		&createdPort.ConnectorType,
		&createdPort.PowerKW,
		&createdPort.IsActive,
		&createdPort.CreatedAt,
	)

	if err != nil {
		return domain.ChargingPort{}, fmt.Errorf("%s: %w", op, err)
	}

	return createdPort, nil
}

func (p *Postgres) UpdatePort(
	ctx context.Context,
	port domain.ChargingPort,
) (domain.ChargingPort, error) {
	const op = "postgres.UpdatePort"

	const query = `
		UPDATE charging_ports
		SET
			connector_type = $1,
			power_kw = $2,
			is_active = $3
		WHERE id = $4
		AND station_id = $5
		RETURNING
			id,
			station_id,
			connector_type,
			power_kw,
			is_active,
			created_at
	`

	var updatedPort domain.ChargingPort

	err := p.db.QueryRow(
		ctx,
		query,
		port.ConnectorType,
		port.PowerKW,
		port.IsActive,
		port.ID,
		port.StationID,
	).Scan(
		&updatedPort.ID,
		&updatedPort.StationID,
		&updatedPort.ConnectorType,
		&updatedPort.PowerKW,
		&updatedPort.IsActive,
		&updatedPort.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ChargingPort{}, fmt.Errorf(
				"%s: %w",
				op,
				domain.ErrPortNotFound,
			)
		}

		return domain.ChargingPort{}, fmt.Errorf(
			"%s: %w",
			op,
			err,
		)
	}

	return updatedPort, nil
}

func (p *Postgres) DeletePort(
	ctx context.Context,
	stationID uuid.UUID,
	portID uuid.UUID,
) error {
	const op = "postgres.DeletePort"

	const query = `
		DELETE FROM charging_ports
		WHERE id = $1
		AND station_id = $2
		RETURNING id
	`

	var deletedID uuid.UUID

	err := p.db.QueryRow(
		ctx,
		query,
		portID,
		stationID,
	).Scan(&deletedID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf(
				"%s: %w",
				op,
				domain.ErrPortNotFound,
			)
		}

		return fmt.Errorf(
			"%s: %w",
			op,
			err,
		)
	}

	return nil
}
