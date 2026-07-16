package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/antongolenev23/voltake-services/pkg/logger"
	"github.com/antongolenev23/voltake-services/services/booking/internal/http-server/dto"
	"github.com/antongolenev23/voltake-services/services/booking/internal/http-server/middleware"
)

func (h *Handler) GetMe(
	w http.ResponseWriter,
	r *http.Request,
) {
	const op = "handler.GetMe"

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	claims, err := middleware.GetClaims(r.Context())

	if err != nil {
		log.Info("failed to get claims", slog.String("error", err.Error()))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	resp := dto.MeResponse{
		Email:   claims.Email,
		IsAdmin: claims.IsAdmin,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error(
			"failed to encode response",
			slog.String("error", err.Error()),
		)
	}
}
