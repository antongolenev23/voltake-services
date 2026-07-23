package postgres

import (
	"context"
	"time"
)

func (p *Postgres) ExpireUnclaimedBookings(
	ctx context.Context,
	timeout time.Duration,
) (int64, error) {

	query := `
		UPDATE bookings
		SET status = 'expired'
		WHERE status = 'booked'
		AND start_time + $1 < NOW()
	`

	tag, err := p.db.Exec(
		ctx,
		query,
		timeout,
	)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func (p *Postgres) CompleteActiveBookings(
	ctx context.Context,
) (int64, error) {

	query := `
		UPDATE bookings
		SET status = 'completed'
		WHERE status = 'active'
		AND reserved_until < NOW()
	`

	tag, err := p.db.Exec(ctx, query)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}
