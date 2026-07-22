package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	authclient "github.com/antongolenev23/voltake-services/services/booking/internal/auth-client"
	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
	"github.com/antongolenev23/voltake-services/services/booking/internal/http-server/middleware"
)

var ErrRequestBodyTooLarge = errors.New("request body too large")

type StationsProvider interface {
	GetStations(ctx context.Context, limit, offset int) ([]domain.ChargingStation, error)
	GetNearbyStations(
		ctx context.Context,
		lat, lng, radius float64,
		limit, offset int,
	) ([]domain.ChargingStation, error)
	GetStation(ctx context.Context, id uuid.UUID) (domain.ChargingStationDetails, error)
	CreateStation(ctx context.Context, station domain.ChargingStation) (domain.ChargingStation, error)
	UpdateStation(ctx context.Context, station domain.ChargingStation) (domain.ChargingStation, error)
	DeleteStation(ctx context.Context, stationID, ownerID uuid.UUID) error
}

type PortsProvider interface {
	GetPort(ctx context.Context, stationID, portID uuid.UUID) (domain.ChargingPort, error)
	CreatePort(ctx context.Context, port domain.ChargingPort) (domain.ChargingPort, error)
	ActivatePort(ctx context.Context, stationID uuid.UUID, portID uuid.UUID) (domain.ChargingPort, error)
	DeactivatePort(ctx context.Context, stationID uuid.UUID, portID uuid.UUID) (domain.ChargingPort, error)
	DeletePort(ctx context.Context, stationID, portID uuid.UUID) error
}

type BookingProvider interface {
	CreateBooking(ctx context.Context, booking domain.Booking) (domain.Booking, error)
	GetBookings(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Booking, error)
	GetBooking(ctx context.Context, userID, bookingID uuid.UUID) (domain.BookingDetails, error)
	CancelBooking(ctx context.Context, userID, bookingID uuid.UUID) (domain.Booking, error)
}

type Service interface {
	StationsProvider
	PortsProvider
	BookingProvider
}

type ClientInterface interface {
	Register(ctx context.Context, credentials authclient.Credentials) (authclient.AuthResponse, error)
	Login(ctx context.Context, credentials authclient.Credentials) (authclient.AuthResponse, error)

	Close() error
}

type Handler struct {
	service Service
	client  ClientInterface
	Log     *slog.Logger
}

func New(
	service Service,
	client ClientInterface,
	log *slog.Logger,
) *Handler {
	return &Handler{
		service: service,
		client:  client,
		Log:     log,
	}
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any, maxSize int64) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return ErrRequestBodyTooLarge
		}

		return err
	}

	return nil
}

func handleJSONDecodeError(w http.ResponseWriter, log *slog.Logger, err error) {
	if errors.Is(err, ErrRequestBodyTooLarge) {
		log.Info("request body too large")
		http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	log.Info("failed to decode request body",
		slog.String("error", err.Error()),
	)
	http.Error(w, "invalid request body", http.StatusBadRequest)
}

func getUserID(ctx context.Context) (uuid.UUID, error) {
	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	return claims.UserID, nil
}
