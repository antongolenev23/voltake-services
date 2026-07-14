package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

func (s *Service) GetStations(
	ctx context.Context,
	limit, offset int,
) ([]*domain.ChargingStation, error) {
	const op = "service.GetStations"

	stations, err := s.repository.GetStations(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return stations, nil
}

func (s *Service) GetStation(
	ctx context.Context,
	id uuid.UUID,
) (*domain.ChargingStation, error) {
	const op = "service.GetStation"

	station, err := s.repository.GetStation(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return station, nil
}

func (s *Service) CreateStation(
	ctx context.Context,
	station *domain.ChargingStation,
) (*domain.ChargingStation, error) {
	const op = "service.CreateStation"

	if err := station.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	createdStation, err := s.repository.CreateStation(ctx, station)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return createdStation, nil
}

func (s *Service) UpdateStation(
	ctx context.Context,
	station *domain.ChargingStation,
) (*domain.ChargingStation, error) {
	const op = "service.UpdateStation"

	if err := station.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	updatedStation, err := s.repository.UpdateStation(ctx, station)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return updatedStation, nil
}

func (s *Service) DeleteStation(
	ctx context.Context,
	id uuid.UUID,
	ownerID uuid.UUID,
) error {
	const op = "service.DeleteStation"

	err := s.repository.DeleteStation(ctx, id, ownerID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
