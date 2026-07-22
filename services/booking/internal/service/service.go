package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/services/booking/internal/config"
	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

type StationsRepository interface {
	GetStations(ctx context.Context, filter domain.StationFilter) ([]domain.ChargingStation, error)
	GetStation(ctx context.Context, id uuid.UUID) (domain.ChargingStationDetails, error)
	CreateStation(ctx context.Context, station domain.ChargingStation) (domain.ChargingStation, error)
	UpdateStation(ctx context.Context, station domain.ChargingStation) (domain.ChargingStation, error)
	DeleteStation(ctx context.Context, stationID, ownerID uuid.UUID) error
}

type PortsRepository interface {
	GetPort(ctx context.Context, stationID uuid.UUID, portID uuid.UUID) (domain.ChargingPort, error)
	GetBookedIntervals(ctx context.Context, portID uuid.UUID, from time.Time, to time.Time) ([]domain.TimeRange, error)
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
	cfg        *config.BookingConfig
}

func New(
	repository Repository,
	cfg *config.BookingConfig,
) *Service {
	return &Service{
		repository: repository,
		cfg:        cfg,
	}
}
