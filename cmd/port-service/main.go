package main

import (
	"log"
	"log/slog"

	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/logger"
	"github.com/axmz/go-port-service/internal/transport/http"
	"github.com/axmz/go-port-service/pkg/graceful"
	"github.com/axmz/go-port-service/pkg/inmem"

	repo "github.com/axmz/go-port-service/internal/repository/port"
	serv "github.com/axmz/go-port-service/internal/services/port"
)

func run() error {
	cfg := config.MustLoad()

	logger.Setup(cfg.Env)

	slog.Info("Application starting", slog.String("env", cfg.Env))

	d := inmem.NewInMemoryDB[*repo.Port]()

	r := repo.NewPortRepository(d)

	s := serv.NewPortService(r)

	srv := http.StartServer(cfg, s)

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
