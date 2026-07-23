package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/antongolenev23/voltake-services/services/booking/internal/config"
)

type BookingCompleteRepository interface {
	CompleteActiveBookings(ctx context.Context) (int64, error)
}

type BookingComplete struct {
	repo BookingCompleteRepository
	cfg  *config.BookingCompleteWorker
}

func NewBookingCompleteWorker(
	repo BookingCompleteRepository,
	cfg *config.BookingCompleteWorker,
) *BookingComplete {
	return &BookingComplete{
		repo: repo,
		cfg:  cfg,
	}
}

func (w *BookingComplete) Start(ctx context.Context, log *slog.Logger) {

	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	log.Info("booking complete worker started")

	for {
		select {

		case <-ticker.C:

			count, err := w.repo.CompleteActiveBookings(ctx)
			if err != nil {
				log.Error(
					"complete active bookings failed",
					"error",
					err,
				)

				continue
			}

			if count > 0 {
				log.Info(
					"bookings completed",
					"count",
					count,
				)
			}

		case <-ctx.Done():

			log.Info("booking complete worker stopped")

			return
		}
	}
}
