package logger

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/antongolenev23/voltake-services/pkg/types"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func MustInit(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log.Fatalf("can not initialize logger, incorrect env: %s", env)
	}

	return logger
}

func LoggerWithRequestID(log *slog.Logger, ctx context.Context) *slog.Logger {
	requestID, _ := ctx.Value(types.RequestIDKey).(string)

	return log.With(
		slog.String("request_id", requestID),
	)
}
