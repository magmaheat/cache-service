package app

import (
	"github.com/magmaheat/cache-service/config"
	"github.com/magmaheat/cache-service/pkg/postgres"
	log "github.com/sirupsen/logrus"
)

func Run() {
	cfg := config.New()

	setupLogger(cfg.Level)

	log.Info("Initializing postgres...")
	pg := postgres.New(cfg.PG.URL)

}
