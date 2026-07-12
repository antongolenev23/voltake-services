package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/antongolenev23/voltake-services/pkg/types"
)

func RequestID() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req any,
		reply any,
		cc *grpc.ClientConn,
		invoke grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		requestID, _ := ctx.Value(types.RequestIDKey).(string)

		if requestID != "" {
			ctx = metadata.AppendToOutgoingContext(
				ctx,
				"x-request-id",
				requestID,
			)
		}

		return invoke(ctx, method, req, reply, cc, opts...)
	}
}
