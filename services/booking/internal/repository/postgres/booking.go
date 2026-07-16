package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

func (p *Postgres) CreateBooking(ctx context.Context, booking domain.Booking) (domain.Booking, error) {
	const op = "postgres.CreateBooking"

	const query = `
		INSERT INTO bookings (
			user_id,
			port_id,
			start_time,
			end_time
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			user_id,
			port_id,
			start_time,
			end_time,
			status,
			created_at,
			updated_at
	`

	var result domain.Booking

	err := p.db.QueryRow(
		ctx,
		query,
		booking.UserID,
		booking.PortID,
		booking.StartTime,
		booking.EndTime,
	).Scan(
		&result.ID,
		&result.UserID,
		&result.PortID,
		&result.StartTime,
		&result.EndTime,
		&result.Status,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23P01":
				return domain.Booking{}, fmt.Errorf("%s: %w", op, domain.ErrBookingConflict)
			case "23503":
				return domain.Booking{}, fmt.Errorf("%s: %w", op, domain.ErrPortNotFound)
			}
		}

		return domain.Booking{}, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (p *Postgres) GetBookings(
	ctx context.Context,
	userID uuid.UUID,
	limit, offset int,
) ([]domain.Booking, error) {
	const op = "postgres.GetBookings"

	const query = `
		SELECT
			id,
			user_id,
			port_id,
			start_time,
			end_time,
			status,
			created_at,
			updated_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY start_time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := p.db.Query(ctx, query, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	bookings := make([]domain.Booking, 0)

	for rows.Next() {
		var booking domain.Booking

		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.PortID,
			&booking.StartTime,
			&booking.EndTime,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return bookings, nil
}

func (p *Postgres) GetBooking(
	ctx context.Context,
	userID uuid.UUID,
	bookingID uuid.UUID,
) (domain.BookingDetails, error) {
	const op = "postgres.GetBooking"

	const query = `
		SELECT
			b.id,
			b.user_id,
			b.port_id,
			b.start_time,
			b.end_time,
			b.status,
			b.created_at,
			b.updated_at,

			cp.id,
			cp.station_id,
			cp.connector_type,
			cp.power_kw,
			cp.is_active,
			cp.created_at,

			cs.id,
			cs.name,
			cs.address,
			cs.latitude,
			cs.longitude,
			cs.created_at

		FROM bookings b

		JOIN charging_ports cp
			ON cp.id = b.port_id

		JOIN charging_stations cs
			ON cs.id = cp.station_id

		WHERE b.id = $1
		AND b.user_id = $2
	`

	var booking domain.BookingDetails

	err := p.db.QueryRow(
		ctx,
		query,
		bookingID,
		userID,
	).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.PortID,
		&booking.StartTime,
		&booking.EndTime,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,

		&booking.Port.ID,
		&booking.Port.StationID,
		&booking.Port.ConnectorType,
		&booking.Port.PowerKW,
		&booking.Port.IsActive,
		&booking.Port.CreatedAt,

		&booking.Station.ID,
		&booking.Station.Name,
		&booking.Station.Address,
		&booking.Station.Latitude,
		&booking.Station.Longitude,
		&booking.Station.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.BookingDetails{}, fmt.Errorf("%s: %w", op, domain.ErrBookingNotFound)
		}

		return domain.BookingDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	return booking, nil
}

func (p *Postgres) CancelBooking(
	ctx context.Context,
	userID uuid.UUID,
	bookingID uuid.UUID,
) (domain.Booking, error) {
	const op = "postgres.CancelBooking"

	const query = `
		UPDATE bookings
		SET
			status = 'cancelled',
			updated_at = NOW()
		WHERE id = $1
		AND user_id = $2
		AND status != 'completed'
		RETURNING
			id,
			user_id,
			port_id,
			start_time,
			end_time,
			status,
			created_at,
			updated_at
	`

	var booking domain.Booking

	err := p.db.QueryRow(
		ctx,
		query,
		bookingID,
		userID,
	).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.PortID,
		&booking.StartTime,
		&booking.EndTime,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Booking{}, fmt.Errorf("%s: %w", op, domain.ErrBookingNotFound)
		}

		return domain.Booking{}, fmt.Errorf("%s: %w", op, err)
	}

	return booking, nil
}
