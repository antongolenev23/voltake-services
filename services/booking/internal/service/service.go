package service

import (
	"context"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

type Repository interface {
	GetStations(ctx context.Context) ([]*domain.ChargingStation, error)
}

type Service struct {
	repository Repository
}

func New(
	repository Repository,
) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetStations(
	ctx context.Context,
) ([]*domain.ChargingStation, error) {

	stations, err := s.repository.GetStations(ctx)
	if err != nil {
		return nil, err
	}

	return stations, nil
}
