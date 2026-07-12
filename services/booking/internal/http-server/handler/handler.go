package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	authclient "github.com/antongolenev23/voltake-services/services/booking/internal/auth-client"
)

type Booking interface {
}

type ClientInterface interface {
	Register(ctx context.Context, credentials authclient.Credentials) (authclient.AuthResponse, error)
	Login(ctx context.Context, credentials authclient.Credentials) (authclient.AuthResponse, error)

	Close() error
}

type Handler struct {
	booking Booking
	client ClientInterface
	log *slog.Logger
}

func New(
	booking Booking,
	client ClientInterface,
	log *slog.Logger,
) *Handler {
	return &Handler{
		booking: booking,
		client: client,
		log: log,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req authclient.Credentials

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.client.Register(r.Context(), req)
	if err != nil {
		h.log.Error("failed to register", slog.String("error", err.Error()))
		http.Error(w, "failed to register", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req authclient.Credentials

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.client.Login(r.Context(), req)
	if err != nil {
		http.Error(w, "failed to login", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
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
