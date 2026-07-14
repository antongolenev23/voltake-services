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

func (p *Postgres) GetStations(
	ctx context.Context,
	limit, offset int,
) ([]*domain.ChargingStation, error) {
	const op = "postgres.GetStations"

	const query = `
		SELECT
			id,
			name,
			address,
			latitude,
			longitude,
			created_at
		FROM charging_stations
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := p.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var stations []*domain.ChargingStation

	for rows.Next() {

		var (
			id        uuid.UUID
			name      string
			address   string
			latitude  float64
			longitude float64
			createdAt time.Time
		)

		err := rows.Scan(
			&id,
			&name,
			&address,
			&latitude,
			&longitude,
			&createdAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		stations = append(
			stations,
			&domain.ChargingStation{
				ID:        id,
				Name:      name,
				Address:   address,
				Latitude:  latitude,
				Longitude: longitude,
				CreatedAt: createdAt,
			},
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return stations, nil
}

func (p *Postgres) GetStation(
	ctx context.Context,
	id uuid.UUID,
) (*domain.ChargingStation, error) {
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
			return nil, fmt.Errorf("%s: %w", op, domain.ErrStationNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &station, nil
}

func (p *Postgres) CreateStation(
	ctx context.Context,
	station *domain.ChargingStation,
) (*domain.ChargingStation, error) {
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

	err := p.db.QueryRow(
		ctx,
		query,
		station.OwnerID,
		station.Name,
		station.Address,
		station.Latitude,
		station.Longitude,
	).Scan(
		&station.ID,
		&station.OwnerID,
		&station.Name,
		&station.Address,
		&station.Latitude,
		&station.Longitude,
		&station.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"%s: %w",
			op,
			err,
		)
	}

	return station, nil
}

func (p *Postgres) UpdateStation(
	ctx context.Context,
	station *domain.ChargingStation,
) (*domain.ChargingStation, error) {
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
		&station.ID,
		&station.OwnerID,
		&station.Name,
		&station.Address,
		&station.Latitude,
		&station.Longitude,
		&station.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrStationNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return station, nil
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
