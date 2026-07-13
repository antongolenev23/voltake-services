package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	authclient "github.com/antongolenev23/voltake-services/services/booking/internal/auth-client"
)

var ErrRequestBodyTooLarge = errors.New("request body too large")

type Booking interface {
}

type ClientInterface interface {
	Register(ctx context.Context, credentials authclient.Credentials) (authclient.AuthResponse, error)
	Login(ctx context.Context, credentials authclient.Credentials) (authclient.AuthResponse, error)

	Close() error
}

type Handler struct {
	booking Booking
	client  ClientInterface
	Log     *slog.Logger
}

func New(
	booking Booking,
	client ClientInterface,
	log *slog.Logger,
) *Handler {
	return &Handler{
		booking: booking,
		client:  client,
		Log:     log,
	}
}

//  Users

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	// TODO
}

//  Stations

func (h *Handler) GetStations(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) GetStation(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) CreateStation(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) UpdateStation(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) DeleteStation(w http.ResponseWriter, r *http.Request) {
	// TODO
}

//  Ports

func (h *Handler) GetPorts(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) GetPort(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) CreatePort(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) UpdatePort(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) DeletePort(w http.ResponseWriter, r *http.Request) {
	// TODO
}

//  Bookings

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	// TODO
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
		log.Debug("request body too large")
		http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	log.Debug("failed to decode request body",
		slog.String("error", err.Error()),
	)
	http.Error(w, "invalid request body", http.StatusBadRequest)
}
