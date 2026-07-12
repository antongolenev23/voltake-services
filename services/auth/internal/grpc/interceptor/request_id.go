package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/antongolenev23/voltake-services/pkg/types"
)

func RequestID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		requestID := ""

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("x-request-id")

			if len(values) > 0 {
				requestID = values[0]
			}
		}

		ctx = context.WithValue(
			ctx,
			types.RequestIDKey,
			requestID,
		)

		return handler(ctx, req)
	}
}
