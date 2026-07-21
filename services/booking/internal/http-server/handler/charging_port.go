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

func (h *Handler) GetPorts(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetPorts"

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	stationIDParam := chi.URLParam(r, "stationID")

	stationID, err := uuid.Parse(stationIDParam)
	if err != nil {
		log.Info("invalid station id", slog.String("error", err.Error()))
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	ports, err := h.service.GetPorts(ctx, stationID)
	if err != nil {
		if errors.Is(err, domain.ErrStationNotFound) {
			http.Error(w, "station not found", http.StatusNotFound)
			return
		}

		log.Error("failed to get ports", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.NewPortsResponse(ports)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response",
			slog.String("error", err.Error()),
		)

		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) GetPort(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetPort"

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	stationIDParam := chi.URLParam(r, "stationID")
	portIDParam := chi.URLParam(r, "portID")

	stationID, err := uuid.Parse(stationIDParam)
	if err != nil {
		log.Info("invalid station id", slog.String("error", err.Error()))
		http.Error(w, "invalid station id", http.StatusBadRequest)
		return
	}

	portID, err := uuid.Parse(portIDParam)
	if err != nil {
		log.Info("invalid port id", slog.String("error", err.Error()))
		http.Error(w, "invalid port id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	port, err := h.service.GetPort(ctx, stationID, portID)
	if err != nil {
		if errors.Is(err, domain.ErrPortNotFound) {
			http.Error(w, "port not found", http.StatusNotFound)
			return
		}

		log.Error("failed to get port", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.NewPortResponse(port)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", slog.String("error", err.Error()))
	}
}

func (h *Handler) CreatePort(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreatePort"
	const maxRequestBodySize = 1 * 1024 // 1 KB

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	stationIDParam := chi.URLParam(r, "stationID")

	stationID, err := uuid.Parse(stationIDParam)
	if err != nil {
		log.Info("invalid station id", slog.String("error", err.Error()))
		http.Error(w, "invalid station id", http.StatusBadRequest)
		return
	}

	var req dto.PortRequest

	if err := decodeJSONBody(w, r, &req, maxRequestBodySize); err != nil {
		handleJSONDecodeError(w, log, err)
		return
	}

	port := req.ToCreateDomain(stationID)

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	createdPort, err := h.service.CreatePort(ctx, port)
	if err != nil {
		log.Error("failed to create port", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.NewPortResponse(createdPort)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", slog.String("error", err.Error()))
	}
}

func (h *Handler) ActivatePort(w http.ResponseWriter, r *http.Request) {
	const op = "handler.ActivatePort"
	h.changePortStatus(w, r, true, op)
}

func (h *Handler) DeactivatePort(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeactivatePort"
	h.changePortStatus(w, r, false, op)
}

func (h *Handler) changePortStatus(w http.ResponseWriter, r *http.Request, isActive bool, op string) {
	log := logger.WithRequestContext(r.Context(), h.Log, op)

	stationIDParam := chi.URLParam(r, "stationID")
	portIDParam := chi.URLParam(r, "portID")

	stationID, err := uuid.Parse(stationIDParam)
	if err != nil {
		log.Info("invalid station id", slog.String("error", err.Error()))
		http.Error(w, "invalid station id", http.StatusBadRequest)
		return
	}

	portID, err := uuid.Parse(portIDParam)
	if err != nil {
		log.Info("invalid port id", slog.String("error", err.Error()))
		http.Error(w, "invalid port id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var port domain.ChargingPort

	if isActive {
		port, err = h.service.ActivatePort(ctx, stationID, portID)
	} else {
		port, err = h.service.DeactivatePort(ctx, stationID, portID)
	}

	if err != nil {
		if errors.Is(err, domain.ErrPortNotFound) {
			http.Error(w, "port not found", http.StatusNotFound)
			return
		}

		log.Error("failed to change port status", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.NewPortResponse(port)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", slog.String("error", err.Error()))
	}
}

func (h *Handler) DeletePort(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeletePort"

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	stationIDParam := chi.URLParam(r, "stationID")
	portIDParam := chi.URLParam(r, "portID")

	stationID, err := uuid.Parse(stationIDParam)
	if err != nil {
		log.Info("invalid station id", slog.String("error", err.Error()))
		http.Error(w, "invalid station id", http.StatusBadRequest)
		return
	}

	portID, err := uuid.Parse(portIDParam)
	if err != nil {
		log.Info("invalid port id", slog.String("error", err.Error()))
		http.Error(w, "invalid port id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err = h.service.DeletePort(ctx, stationID, portID)
	if err != nil {
		if errors.Is(err, domain.ErrPortNotFound) {
			http.Error(w, "port not found", http.StatusNotFound)
			return
		}

		log.Error("failed to delete port", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
