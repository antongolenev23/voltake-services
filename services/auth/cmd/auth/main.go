package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/antongolenev23/voltake-services/pkg/logger"
	"github.com/antongolenev23/voltake-services/services/auth/internal/app"
	"github.com/antongolenev23/voltake-services/services/auth/internal/config"
)

func main() {
	cfg := config.MustLoad()

	log := logger.MustInit(cfg.Env)
	log = log.With(slog.String("service", "auth"))

	log.Info("starting application", slog.String("env", cfg.Env))
	log.Debug("debug messages enabled")

	a := app.New(cfg, log)
	go a.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sig := <-stop

	log.Info("stopping application", slog.String("signal", sig.String()))

	a.Stop()

	log.Info("application stopped")
}
