package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/antongolenev23/voltake-services/services/auth/internal/config"
	grpcapi "github.com/antongolenev23/voltake-services/services/auth/internal/grpc"
	"github.com/antongolenev23/voltake-services/services/auth/internal/storage/postgres"
	"github.com/antongolenev23/voltake-services/services/auth/internal/usecase"
)

type App struct {
	cfg *config.Config
	log *slog.Logger

	gRPCServer *grpc.Server
}

func New(cfg *config.Config, log *slog.Logger) *App {
	const op = "app.New"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pgxpool, err := postgres.NewPgxpool(ctx, &cfg.Storage)
	if err != nil {
		log.Error("failed to init repo", "error", err)
		os.Exit(1)
	}

	storage := postgres.New(pgxpool)
	usecaseAuth := usecase.New(&cfg.JWT, storage)
	gRPCServer := grpc.NewServer()
	grpcapi.Register(gRPCServer, usecaseAuth, log)

	hs := health.NewServer()
	healthpb.RegisterHealthServer(gRPCServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	return &App{
		cfg:        cfg,
		log:        log,
		gRPCServer: gRPCServer,
	}
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const op = "app.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.cfg.GRPC.Port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.GRPC.Port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "app.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
