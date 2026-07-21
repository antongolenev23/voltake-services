package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	authclient "github.com/antongolenev23/voltake-services/services/booking/internal/auth-client"
	"github.com/antongolenev23/voltake-services/services/booking/internal/config"
	"github.com/antongolenev23/voltake-services/services/booking/internal/http-server/handler"
	"github.com/antongolenev23/voltake-services/services/booking/internal/http-server/router"
	"github.com/antongolenev23/voltake-services/services/booking/internal/repository/postgres"
	"github.com/antongolenev23/voltake-services/services/booking/internal/service"
	"github.com/antongolenev23/voltake-services/services/booking/internal/worker"
)

type App struct {
	cfg *config.Config
	log *slog.Logger

	db         *pgxpool.Pool
	httpServer *http.Server
	authClient *authclient.Client

	bookingCompleteWorker *worker.BookingComplete
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

	authClient, err := authclient.New(cfg.AuthService.Address)
	if err != nil {
		log.Error("failed to connect auth service", "error", err)
		os.Exit(1)
	}

	repository := postgres.New(pgxpool)
	booking := service.New(repository)
	handlerHTTP := handler.New(booking, authClient, log)

	bookingCompleteWorker := worker.NewBookingWorker(repository)

	r := router.New(handlerHTTP, &cfg.JWT)

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.RequestReadTimeout,
		WriteTimeout: cfg.HTTPServer.ResponseWriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return &App{
		cfg:                   cfg,
		log:                   log,
		db:                    pgxpool,
		httpServer:            server,
		authClient:            authClient,
		bookingCompleteWorker: bookingCompleteWorker,
	}

}

func (a *App) Run() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	a.runServer()

	a.bookingCompleteWorker.Start(ctx, a.log)

	defer stop()

	<-ctx.Done()
	a.log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	a.shutdown(ctx)
}

func (a *App) runServer() {
	go func() {
		a.log.Info("starting server",
			slog.String("address", a.cfg.HTTPServer.Address),
		)

		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Error("server stopped", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()
}

func (a *App) shutdown(ctx context.Context) {
	var err error

	if err = a.httpServer.Shutdown(ctx); err != nil {
		a.log.Warn("shutdown error", slog.String("error", err.Error()))
	} else {
		a.log.Info("http server stopped gracefully")
	}

	if err := a.authClient.Close(); err != nil {
		a.log.Warn("failed to close auth client", "error", err)
	} else {
		a.log.Info("auth client connection closed gracefully")
	}

	a.db.Close()
	a.log.Info("database connection pool closed gracefully")
}
