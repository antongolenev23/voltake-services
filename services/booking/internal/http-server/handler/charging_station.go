package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/antongolenev23/voltake-services/pkg/logger"
	"github.com/antongolenev23/voltake-services/services/booking/internal/domain"
)

type StationsProvider interface {
	GetStations(ctx context.Context) ([]*domain.ChargingStation, error)
}

func (h *Handler) GetStations(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetStations"

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	log.Debug("starting get stations request processing")

	stations, err := h.service.GetStations(ctx)
	if err != nil {
		log.Error("failed to get charging stations", slog.String("error", err.Error()))
		http.Error(w, "failed to get charging stations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(stations)
	if err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
		return
	}
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
