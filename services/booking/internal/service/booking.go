package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

func (s *Service) CreateBooking(
	ctx context.Context,
	booking domain.Booking,
) (domain.Booking, error) {
	const op = "service.CreateBooking"

	if err := booking.Validate(s.cfg.MinDuration, s.cfg.MaxDuration); err != nil {
		return domain.Booking{}, fmt.Errorf("%s: %w", op, err)
	}

	booking.ReservedUntil = booking.EndTime.Add(
		s.cfg.Buffer,
	)

	created, err := s.repository.CreateBooking(ctx, booking)
	if err != nil {
		return domain.Booking{}, fmt.Errorf("%s: %w", op, err)
	}

	return created, nil
}

func (s *Service) GetBookings(
	ctx context.Context,
	userID uuid.UUID,
	limit, offset int,
) ([]domain.Booking, error) {
	const op = "service.GetBookings"

	bookings, err := s.repository.GetBookings(ctx, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return bookings, nil
}

func (s *Service) GetBooking(
	ctx context.Context,
	userID uuid.UUID,
	bookingID uuid.UUID,
) (domain.BookingDetails, error) {
	const op = "service.GetBooking"

	booking, err := s.repository.GetBooking(ctx, userID, bookingID)

	if err != nil {
		return domain.BookingDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	return booking, nil
}

func (s *Service) CancelBooking(
	ctx context.Context,
	userID uuid.UUID,
	bookingID uuid.UUID,
) (domain.Booking, error) {
	const op = "service.CancelBooking"

	booking, err := s.repository.CancelBooking(ctx, userID, bookingID)

	if err != nil {
		return domain.Booking{}, fmt.Errorf("%s: %w", op, err)
	}

	return booking, nil
}
