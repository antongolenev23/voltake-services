package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"

	"github.com/antongolenev23/voltake-services/pkg/logger"
	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
	"github.com/antongolenev23/voltake-services/services/booking/internal/http-server/dto"
)

func (h *Handler) GetStations(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetStations"
	log := logger.WithRequestContext(r.Context(), h.Log, op)

	limit, offset := parseLimitAndOffset(r.URL.Query())

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	stations, err := h.service.GetStations(ctx, limit, offset)
	if err != nil {
		log.Error("failed to get charging stations", slog.String("error", err.Error()))
		http.Error(w, "failed to get charging stations", http.StatusInternalServerError)
		return
	}

	resp := dto.NewStationsResponse(stations)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("failed to encode response",
			slog.String("error", err.Error()),
		)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetStation(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetStation"
	log := logger.WithRequestContext(r.Context(), h.Log, op)

	stationID := chi.URLParam(r, "stationID")

	id, err := uuid.Parse(stationID)
	if err != nil {
		log.Info("invalid station id", slog.String("error", err.Error()))
		http.Error(w, "invalid station id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	station, err := h.service.GetStation(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrStationNotFound) {
			http.Error(w, "station not found", http.StatusNotFound)
			return
		}

		log.Error("failed to get station", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.NewStationResponse(station)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response",
			slog.String("error", err.Error()),
		)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (h *Handler) CreateStation(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateStation"
	const maxRequestBodySize = 1 * 1024 // 1 KB

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	userID, err := getUserID(r.Context())
	if err != nil {
		log.Info("failed to get user id", slog.String("error", err.Error()))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.StationRequest
	if err := decodeJSONBody(w, r, &req, maxRequestBodySize); err != nil {
		handleJSONDecodeError(w, log, err)
		return
	}

	station := req.ToCreateDomain(userID)

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	createdStation, err := h.service.CreateStation(ctx, station)
	if err != nil {
		log.Error("failed to create station", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.NewStationResponse(createdStation)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error(
			"failed to encode response",
			slog.String("error", err.Error()),
		)
	}
}

func (h *Handler) UpdateStation(
	w http.ResponseWriter,
	r *http.Request,
) {
	const op = "handler.UpdateStation"
	const maxRequestBodySize = 1 * 1024 // 1 KB

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	stationIDParam := chi.URLParam(r, "stationID")

	stationID, err := uuid.Parse(stationIDParam)
	if err != nil {
		log.Info("invalid station id", slog.String("error", err.Error()))
		http.Error(w, "invalid station id", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r.Context())
	if err != nil {
		log.Info("failed to get user id", slog.String("error", err.Error()))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.StationRequest

	if err := decodeJSONBody(w, r, &req, maxRequestBodySize); err != nil {
		handleJSONDecodeError(w, log, err)
		return
	}

	station := req.ToUpdateDomain(stationID, userID)

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	updatedStation, err := h.service.UpdateStation(ctx, station)
	if err != nil {
		if errors.Is(err, domain.ErrStationNotFound) {
			http.Error(w, "station not found", http.StatusNotFound)
			return
		}

		log.Error("failed to update station", slog.String("error", err.Error()))

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.NewStationResponse(updatedStation)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", slog.String("error", err.Error()))
	}
}

func (h *Handler) DeleteStation(
	w http.ResponseWriter,
	r *http.Request,
) {
	const op = "handler.DeleteStation"

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	stationIDParam := chi.URLParam(r, "stationID")

	stationID, err := uuid.Parse(stationIDParam)
	if err != nil {
		log.Info("invalid station id", slog.String("error", err.Error()))
		http.Error(w, "invalid station id", http.StatusBadRequest)
		return
	}

	ownerID, err := getUserID(r.Context())
	if err != nil {
		log.Info("failed to get user id", slog.String("error", err.Error()))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err = h.service.DeleteStation(ctx, stationID, ownerID)

	if err != nil {
		if errors.Is(err, domain.ErrStationNotFound) {
			http.Error(w, "station not found", http.StatusNotFound)
			return
		}

		log.Error("failed to delete station", slog.String("error", err.Error()))

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseLimitAndOffset(values url.Values) (int, int) {
	offset, err := strconv.Atoi(values.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(values.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	return limit, offset
}
