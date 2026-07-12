package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/antongolenev23/voltake-services/pkg/types"
)

func Logging(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		start := time.Now()

		resp, err := handler(ctx, req)

		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return resp, err
		}

		requestID, _ := ctx.Value(types.RequestIDKey).(string)

		fields := []any{
			slog.String("grpc_method", info.FullMethod),
			slog.Duration("duration", time.Since(start)),
			slog.String("request_id", requestID),
		}

		if err != nil {
			st, _ := status.FromError(err)

			fields = append(
				fields,
				slog.String("grpc_code", st.Code().String()),
				slog.String("error", st.Message()),
			)
		}

		log.Info("grpc request completed", fields...)

		return resp, err
	}
}
