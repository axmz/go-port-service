package main

import (
	"log"
	"log/slog"

	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/logger"
	"github.com/axmz/go-port-service/pkg/graceful"
	"github.com/axmz/go-port-service/pkg/inmem"

	repository "github.com/axmz/go-port-service/internal/repository/port"
	services "github.com/axmz/go-port-service/internal/services/port"
	httpServer "github.com/axmz/go-port-service/internal/transport/http"
)

func run() error {
	cfg := config.MustLoad()

	logger.Setup(cfg.Env)

	slog.Info("Application starting", slog.String("env", cfg.Env))

	d := inmem.New[*repository.Port]()

	r := repository.New(d)

	s := services.New(r)

	srv := httpServer.Start(cfg, s)

	<-graceful.Shutdown(cfg.GracefulTimeout, map[string]graceful.Operation{
		"database":    d.Shutdown,
		"http-server": srv.Shutdown,
	})

	slog.Info("Application stopped")

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
