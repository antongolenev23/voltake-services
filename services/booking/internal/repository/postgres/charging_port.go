package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

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

func (p *Postgres) SetPortActive(
	ctx context.Context,
	stationID uuid.UUID,
	portID uuid.UUID,
	isActive bool,
) (domain.ChargingPort, error) {
	const op = "postgres.SetPortActive"

	const query = `
		UPDATE charging_ports
		SET
			is_active = $1
		WHERE id = $2
		AND station_id = $3
		RETURNING
			id,
			station_id,
			connector_type,
			power_kw,
			is_active,
			created_at
	`

	var port domain.ChargingPort

	err := p.db.QueryRow(
		ctx,
		query,
		isActive,
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

		return domain.ChargingPort{}, fmt.Errorf(
			"%s: %w",
			op,
			err,
		)
	}

	return port, nil
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
