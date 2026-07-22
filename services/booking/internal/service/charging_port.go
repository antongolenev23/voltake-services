package service

import (
	"context"
	"fmt"

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
