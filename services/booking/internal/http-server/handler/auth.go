package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/antongolenev23/voltake-services/pkg/logger"
	authclient "github.com/antongolenev23/voltake-services/services/booking/internal/auth-client"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "handler.Register"
	const maxRequestBodySize = 1 * 1024 // 1 KB

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	var req authclient.Credentials
	if err := decodeJSONBody(w, r, &req, maxRequestBodySize); err != nil {
		handleJSONDecodeError(w, log, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := h.client.Register(ctx, req)
	if err != nil {
		handleRegisterError(err, log, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response",
			slog.String("error", err.Error()),
		)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handler.Login"
	const maxRequestBodySize = 1 * 1024 // 1 KB

	log := logger.WithRequestContext(r.Context(), h.Log, op)

	var req authclient.Credentials
	if err := decodeJSONBody(w, r, &req, maxRequestBodySize); err != nil {
		handleJSONDecodeError(w, log, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := h.client.Login(ctx, req)
	if err != nil {
		handleLoginError(err, log, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response",
			slog.String("error", err.Error()),
		)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func handleRegisterError(err error, log *slog.Logger, w http.ResponseWriter) {
	st, ok := status.FromError(err)
	if !ok {
		log.Error("failed to register", slog.String("error", err.Error()))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	switch st.Code() {
	case codes.InvalidArgument:
		log.Info("invalid register request", slog.String("error", st.Message()))
		http.Error(w, "invalid register request", http.StatusBadRequest)

	case codes.AlreadyExists:
		log.Info("user already exists", slog.String("error", st.Message()))
		http.Error(w, "user already exists", http.StatusConflict)

	case codes.DeadlineExceeded:
		log.Error("auth service timeout",
			slog.String("error", st.Message()),
		)
		http.Error(w, "service unavailable", http.StatusGatewayTimeout)

	default:
		log.Error("failed to register", slog.String("error", st.Message()))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func handleLoginError(err error, log *slog.Logger, w http.ResponseWriter) {
	st, ok := status.FromError(err)
	if !ok {
		log.Error("failed to login", slog.String("error", err.Error()))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	switch st.Code() {
	case codes.Unauthenticated:
		log.Info("invalid credentials", slog.String("error", st.Message()))
		http.Error(w, "invalid credentials", http.StatusUnauthorized)

	case codes.DeadlineExceeded:
		log.Error("auth service timeout",
			slog.String("error", st.Message()),
		)
		http.Error(w, "service unavailable", http.StatusGatewayTimeout)

	default:
		log.Error("failed to login", slog.String("error", st.Message()))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
