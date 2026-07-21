package postgres

import "context"

func (p *Postgres) CompleteExpiredBookings(ctx context.Context) (int64, error) {

	query := `
		UPDATE bookings
		SET status = 'completed'
		WHERE status = 'booked'
		AND end_time < NOW()
	`

	tag, err := p.db.Exec(ctx, query)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}
