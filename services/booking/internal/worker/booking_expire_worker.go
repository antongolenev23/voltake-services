package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/antongolenev23/voltake-services/services/booking/internal/config"
)

type BookingRepository interface {
	ExpireUnclaimedBookings(ctx context.Context, timeout time.Duration) (int64, error)
}

type BookingExpire struct {
	repo BookingRepository
	cfg  *config.Config
}

func NewBookingExpireWorker(repo BookingRepository, cfg *config.Config) *BookingExpire {
	return &BookingExpire{
		repo: repo,
		cfg:  cfg,
	}
}

func (w *BookingExpire) Start(ctx context.Context, log *slog.Logger) {
	const op = "worker.BookingExpireWorker"

	log = log.With(slog.String("operation", op))

	ticker := time.NewTicker(w.cfg.Worker.BookingExpire.Interval)
	defer ticker.Stop()

	log.Info("booking expire worker started")

	for {
		select {

		case <-ticker.C:

			count, err := w.repo.ExpireUnclaimedBookings(
				ctx, w.cfg.DomainRules.Booking.CheckInTimeout,
			)

			if err != nil {
				log.Error("expire unclaimed bookings failed", "error", err)
				continue
			}

			if count > 0 {
				log.Info("bookings expired", "count", count)
			}

		case <-ctx.Done():

			log.Info("booking expire worker stopped")

			return
		}
	}
}
