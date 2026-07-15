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

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/antongolenev23/voltake-services/services/auth/internal/config"
	grpcapi "github.com/antongolenev23/voltake-services/services/auth/internal/grpc"
	"github.com/antongolenev23/voltake-services/services/auth/internal/grpc/interceptor"
	"github.com/antongolenev23/voltake-services/services/auth/internal/repository/postgres"
	"github.com/antongolenev23/voltake-services/services/auth/internal/service"
)

type App struct {
	cfg *config.Config
	log *slog.Logger

	db         *pgxpool.Pool
	gRPCServer *grpc.Server
}

func New(cfg *config.Config, log *slog.Logger) *App {
	const op = "app.New"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pgxpool, err := postgres.NewPgxpool(ctx, &cfg.Repository)
	if err != nil {
		log.Error("failed to init repo", "error", err)
		os.Exit(1)
	}

	repository := postgres.New(pgxpool)
	serviceAuth := service.New(&cfg.JWT, repository)
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
			interceptor.RequestID(),
			interceptor.Logging(log),
		),
	)
	grpcapi.Register(gRPCServer, serviceAuth, log)

	hs := health.NewServer()
	healthpb.RegisterHealthServer(gRPCServer, hs)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	return &App{
		cfg:        cfg,
		log:        log,
		db:         pgxpool,
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

	a.log.With(slog.String("op", op))

	a.log.Info("stopping gRPC server")

	done := make(chan struct{})

	go func() {
		a.gRPCServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		a.log.Info("gRPC server stopped gracefully")
	case <-time.After(10 * time.Second):
		a.log.Warn("gRPC graceful shutdown timeout exceeded, forcing stop")
		a.gRPCServer.Stop()
	}

	a.db.Close()
	a.log.Info("database connection pool closed gracefully")
}
