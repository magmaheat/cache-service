package app

import (
	"fmt"
	"github.com/magmaheat/cache-service/config"
	"github.com/magmaheat/cache-service/intarnal/http"
	"github.com/magmaheat/cache-service/intarnal/repo"
	"github.com/magmaheat/cache-service/intarnal/service"
	"github.com/magmaheat/cache-service/pkg/httpserver"
	"github.com/magmaheat/cache-service/pkg/postgres"
	"github.com/magmaheat/cache-service/pkg/redis"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	cfg := config.New()

	setupLogger(cfg.Level)

	log.Info("Initializing postgres...")
	pg := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.MaxPoolSize))

	log.Info("Initializing redis...")
	rd := redis.New(cfg.RD.URL)

	log.Info("Initializing repositories...")
	repositories := repo.New(pg, rd)

	log.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos:    repositories,
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}
	services := service.New(deps)

	log.Info("Initializing handlers and routers...")
	handlers := http.Init(services)

	log.Info("Starting http server...")
	log.Debugf("Server port: %s", cfg.Port)
	httpServer := httpserver.New(handlers, httpserver.Port(cfg.HTTP.Port))

	log.Info("Configuration graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %v", err))
	}

	log.Info("Shutting down...")
	err := httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %v", err))
	}
}
