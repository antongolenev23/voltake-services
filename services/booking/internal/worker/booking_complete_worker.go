package worker

import (
	"context"
	"log/slog"
	"time"
)

type BookingRepository interface {
	CompleteExpiredBookings(ctx context.Context) (int64, error)
}

type BookingComplete struct {
	repo BookingRepository
}

func NewBookingWorker(repo BookingRepository) *BookingComplete {
	return &BookingComplete{
		repo: repo,
	}
}

func (w *BookingComplete) Start(ctx context.Context, log *slog.Logger) {

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	log.Info("booking complete worker started")

	for {
		select {

		case <-ticker.C:

			count, err := w.repo.CompleteExpiredBookings(ctx)
			if err != nil {
				log.Error(
					"complete expired bookings failed",
					"error",
					err,
				)

				continue
			}

			if count > 0 {
				log.Info("bookings completed", "count", count)
			}

		case <-ctx.Done():

			log.Info("booking complete worker stopped")

			return
		}
	}
}
