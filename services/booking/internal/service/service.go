package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

type StationsRepository interface {
	GetStations(ctx context.Context, limit, offset int) ([]domain.ChargingStation, error)
	GetNearbyStations(
		ctx context.Context,
		lat, lng, radius float64,
		limit, offset int,
	) ([]domain.ChargingStation, error)
	GetStation(ctx context.Context, id uuid.UUID) (domain.ChargingStation, error)
	CreateStation(ctx context.Context, station domain.ChargingStation) (domain.ChargingStation, error)
	UpdateStation(ctx context.Context, station domain.ChargingStation) (domain.ChargingStation, error)
	DeleteStation(ctx context.Context, stationID, ownerID uuid.UUID) error
}

type PortsRepository interface {
	GetPorts(ctx context.Context, stationID uuid.UUID) ([]domain.ChargingPort, error)
	GetPort(ctx context.Context, stationID uuid.UUID, portID uuid.UUID) (domain.ChargingPort, error)
	CreatePort(ctx context.Context, port domain.ChargingPort) (domain.ChargingPort, error)
	SetPortActive(ctx context.Context, stationID uuid.UUID, portID uuid.UUID, isActive bool) (domain.ChargingPort, error)
	DeletePort(ctx context.Context, stationID uuid.UUID, portID uuid.UUID) error
}

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking domain.Booking) (domain.Booking, error)
	GetBookings(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Booking, error)
	GetBooking(ctx context.Context, userID, bookingID uuid.UUID) (domain.BookingDetails, error)
	CancelBooking(ctx context.Context, userID, bookingID uuid.UUID) (domain.Booking, error)
}

type Repository interface {
	StationsRepository
	PortsRepository
	BookingRepository
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
