package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/antongolenev23/voltake-services/pkg/logger"
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "handler.Register"
	log := logger.LoggerWithRequestID(h.Log, r.Context())
	log = log.With(slog.String("operation", op))

	log.Info("starting register request processing")

	var req authclient.Credentials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("failed to decode request body", slog.String("error", err.Error()))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.client.Register(r.Context(), req)
	if err != nil {
		sendRegisterErrorResponse(err, log, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to register", slog.String("error", err.Error()))
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

func sendRegisterErrorResponse(err error, log *slog.Logger, w http.ResponseWriter) {
	st, ok := status.FromError(err)
	if !ok {
		log.Error("failed to register", slog.String("error", err.Error()))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	switch st.Code() {
	case codes.InvalidArgument:
		log.Info("invalid register request", slog.String("error", st.Message()))
		http.Error(w, st.Message(), http.StatusBadRequest)

	case codes.AlreadyExists:
		log.Info("user already exists", slog.String("error", st.Message()))
		http.Error(w, st.Message(), http.StatusConflict)

	default:
		log.Error("failed to register", slog.String("error", st.Message()))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
