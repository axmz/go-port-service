package main

import (
	"log"

	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/transport/http"
	"github.com/axmz/go-port-service/pkg/graceful"
	"github.com/axmz/go-port-service/pkg/inmem"

	repo "github.com/axmz/go-port-service/internal/repository/port"
	serv "github.com/axmz/go-port-service/internal/services/port"
)

func run() error {
	cfg := config.LoadConfig()

	d := inmem.NewInMemoryDB[*repo.Port]()

	r := repo.NewPortRepository(d)

	s := serv.NewPortService(r)

	srv := http.StartServer(cfg, s)

	<-graceful.Shutdown(cfg.GracefulTimeout, map[string]graceful.Operation{
		"database":    d.Shutdown,
		"http-server": srv.Shutdown,
	})

	log.Println("Application stopped")

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
