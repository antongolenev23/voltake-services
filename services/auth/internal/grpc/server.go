package grpc

import (
	"context"
	"log/slog"

	authv1 "github.com/antongolenev23/voltake-protos/gen/go/auth/v1"
	"google.golang.org/grpc"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
	) (token string, err error)

	Register(ctx context.Context,
		email string,
		password string,
	) (token string, err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
	log  *slog.Logger
}

func Register(gRPC *grpc.Server, auth Auth, log *slog.Logger) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth, log: log})
}
