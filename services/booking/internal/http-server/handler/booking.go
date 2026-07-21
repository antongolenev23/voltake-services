package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/pkg/logger"
	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
	"github.com/antongolenev23/voltake-services/services/booking/internal/http-server/dto"
)

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateBooking"
	const maxRequestBodySize = 1 * 1024 // 1 KB

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	userID, err := getUserID(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	portID, err := uuid.Parse(chi.URLParam(r, "portID"))
	if err != nil {
		http.Error(w, "invalid port id", http.StatusBadRequest)
		return
	}

	var req dto.BookingRequest
	if err := decodeJSONBody(w, r, &req, maxRequestBodySize); err != nil {
		handleJSONDecodeError(w, log, err)
		return
	}

	booking := req.ToDomain(userID, portID)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	created, err := h.service.CreateBooking(ctx, booking)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBookingConflict), errors.Is(err, domain.ErrPortUnavailable):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, domain.ErrPortNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, domain.ErrBookingInPast), errors.Is(err, domain.ErrBookingTooLong):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Error("failed to create booking", slog.String("error", err.Error()))
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(dto.NewBookingResponse(created)); err != nil {
		log.Error("failed to encode response", slog.String("error", err.Error()))
	}
}

func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetBookings"

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	userID, err := getUserID(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	limit, offset := parseLimitAndOffset(r.URL.Query())

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	bookings, err := h.service.GetBookings(ctx, userID, limit, offset)

	if err != nil {
		log.Error("failed to get bookings", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.NewBookingListResponse(bookings)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", slog.String("error", err.Error()))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetBooking(
	w http.ResponseWriter,
	r *http.Request,
) {
	const op = "handler.GetBooking"

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	userID, err := getUserID(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	bookingID, err := uuid.Parse(chi.URLParam(r, "bookingID"))
	if err != nil {
		http.Error(w, "invalid booking id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	booking, err := h.service.GetBooking(ctx, userID, bookingID)

	if err != nil {
		if errors.Is(err, domain.ErrBookingNotFound) {
			http.Error(w, "booking not found", http.StatusNotFound)
			return
		}

		log.Error("failed to get booking", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.NewBookingDetailsResponse(booking)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response", slog.String("error", err.Error()))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CancelBooking(
	w http.ResponseWriter,
	r *http.Request,
) {
	const op = "handler.CancelBooking"

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	userID, err := getUserID(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	bookingID, err := uuid.Parse(chi.URLParam(r, "bookingID"))

	if err != nil {
		http.Error(w, "invalid booking id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	booking, err := h.service.CancelBooking(ctx, userID, bookingID)

	if err != nil {
		if errors.Is(err, domain.ErrBookingNotFound) {
			http.Error(w, "booking not found", http.StatusNotFound)
			return
		}

		log.Error("failed to cancel booking", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	resp := dto.NewBookingResponse(booking)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", slog.String("error", err.Error()))
	}
}
