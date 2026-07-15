package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

func (s *Service) GetPorts(
	ctx context.Context,
	stationID uuid.UUID,
) ([]domain.ChargingPort, error) {
	const op = "service.GetPorts"

	ports, err := s.repository.GetPorts(ctx, stationID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ports, nil
}

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

func (s *Service) UpdatePort(
	ctx context.Context,
	port domain.ChargingPort,
) (domain.ChargingPort, error) {
	const op = "service.UpdatePort"

	updatedPort, err := s.repository.UpdatePort(ctx, port)

	if err != nil {
		return domain.ChargingPort{}, fmt.Errorf(
			"%s: %w",
			op,
			err,
		)
	}

	return updatedPort, nil
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
