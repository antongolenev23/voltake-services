package main

import (
	"github.com/antongolenev23/voltake-services/pkg/logger"
	"github.com/antongolenev23/voltake-services/services/booking/internal/app"
	"github.com/antongolenev23/voltake-services/services/booking/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log := logger.MustInit(cfg.Env)

	app := app.New(cfg, log)
	app.Run()
}
