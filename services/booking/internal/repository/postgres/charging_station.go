package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

func (p *Postgres) GetStations(
	ctx context.Context,
) ([]*domain.ChargingStation, error) {

	const query = `
		SELECT
			id,
			name,
			address,
			latitude,
			longitude,
			created_at
		FROM charging_stations
		ORDER BY created_at DESC
	`

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stations []*domain.ChargingStation

	for rows.Next() {

		var (
			id        uuid.UUID
			name      string
			address   string
			latitude  *float64
			longitude *float64
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
			return nil, err
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
		return nil, err
	}

	return stations, nil
}
