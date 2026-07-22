package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

func (s *Service) GetPort(
	ctx context.Context,
	stationID uuid.UUID,
	portID uuid.UUID,
) (domain.ChargingPort, error) {
	const op = "service.GetPort"

	port, err := s.repository.GetPort(
		ctx,
		stationID,
		portID,
	)

	if err != nil {
		return domain.ChargingPort{}, fmt.Errorf("%s: %w", op, err)
	}

	return port, nil
}

func (s *Service) GetPortAvailability(
	ctx context.Context,
	portID uuid.UUID,
	date time.Time,
) ([]domain.TimeRange, error) {
	const op = "service.GetPortAvailability"

	loc := time.UTC

	dayStart := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0,
		loc,
	)

	dayEnd := dayStart.Add(24 * time.Hour)

	bookings, err := s.repository.GetBookedIntervals(
		ctx,
		portID,
		dayStart,
		dayEnd,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for i := range bookings {
		bookings[i].Start = bookings[i].Start.Add(-s.cfg.Buffer)
	}

	return calculateAvailableSlots(
		dayStart,
		dayEnd,
		bookings,
		s.cfg.MinDuration,
	), nil
}

func (s *Service) CreatePort(
	ctx context.Context,
	port domain.ChargingPort,
) (domain.ChargingPort, error) {
	const op = "service.CreatePort"

	createdPort, err := s.repository.CreatePort(ctx, port)
	if err != nil {
		return domain.ChargingPort{}, fmt.Errorf("%s: %w", op, err)
	}

	return createdPort, nil
}

func (s *Service) ActivatePort(
	ctx context.Context,
	stationID uuid.UUID,
	portID uuid.UUID,
) (domain.ChargingPort, error) {
	const op = "service.ActivatePort"

	port, err := s.repository.SetPortActive(ctx, stationID, portID, true)

	if err != nil {
		return domain.ChargingPort{}, fmt.Errorf("%s: %w", op, err)
	}

	return port, nil
}

func (s *Service) DeactivatePort(
	ctx context.Context,
	stationID uuid.UUID,
	portID uuid.UUID,
) (domain.ChargingPort, error) {
	const op = "service.DeactivatePort"

	port, err := s.repository.SetPortActive(ctx, stationID, portID, false)

	if err != nil {
		return domain.ChargingPort{}, fmt.Errorf("%s: %w", op, err)
	}

	return port, nil
}

func (s *Service) DeletePort(
	ctx context.Context,
	stationID uuid.UUID,
	portID uuid.UUID,
) error {
	const op = "service.DeletePort"

	err := s.repository.DeletePort(
		ctx,
		stationID,
		portID,
	)

	if err != nil {
		return fmt.Errorf(
			"%s: %w",
			op,
			err,
		)
	}

	return nil
}

func calculateAvailableSlots(
	start time.Time,
	end time.Time,
	busy []domain.TimeRange,
	minDuration time.Duration,
) []domain.TimeRange {
	slots := make([]domain.TimeRange, 0)

	cursor := start

	for _, item := range busy {
		if cursor.Before(item.Start) {
			slot := domain.TimeRange{
				Start: cursor,
				End:   item.Start,
			}

			if slot.End.Sub(slot.Start) >= minDuration {
				slots = append(slots, slot)
			}
		}

		if item.End.After(cursor) {
			cursor = item.End
		}
	}

	if cursor.Before(end) {
		slot := domain.TimeRange{
			Start: cursor,
			End:   end,
		}

		if slot.End.Sub(slot.Start) >= minDuration {
			slots = append(slots, slot)
		}
	}

	return slots
}
