package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

type StationsRepository interface {
	GetStations(ctx context.Context, limit, offset int) ([]domain.ChargingStation, error)
	GetStation(ctx context.Context, id uuid.UUID) (domain.ChargingStation, error)
	CreateStation(ctx context.Context, station domain.ChargingStation) (domain.ChargingStation, error)
	UpdateStation(ctx context.Context, station domain.ChargingStation) (domain.ChargingStation, error)
	DeleteStation(ctx context.Context, stationID, ownerID uuid.UUID) error
}

type PortsRepository interface {
	GetPorts(ctx context.Context, stationID uuid.UUID) ([]domain.ChargingPort, error)
	GetPort(ctx context.Context, stationID uuid.UUID, portID uuid.UUID) (domain.ChargingPort, error)
	CreatePort(ctx context.Context, port domain.ChargingPort) (domain.ChargingPort, error)
	UpdatePort(ctx context.Context, port domain.ChargingPort) (domain.ChargingPort, error)
	DeletePort(ctx context.Context, stationID uuid.UUID, portID uuid.UUID) error
}

type Repository interface {
	StationsRepository
	PortsRepository
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
