package main

import (
	"context"
	"log"
	"time"

	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/transport/http"
	"github.com/axmz/go-port-service/pkg/graceful"
	"github.com/axmz/go-port-service/pkg/inmem"

	repo "github.com/axmz/go-port-service/internal/repository/port"
	serv "github.com/axmz/go-port-service/internal/services/port"
)

func run() error {
	cfg := config.LoadConfig()

	d := inmem.NewInMemoryDB()

	r := repo.NewPortRepository(d)

	s := serv.NewPortService(r)

	srv := http.StartServer(cfg, s)

	wait := graceful.Shutdown(
		2*time.Second,
		map[string]func(ctx context.Context) error{
			"database": func(ctx context.Context) error {
				return d.Shutdown()
			},
			"http-server": func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			},
		})

	<-wait

	log.Println("Application stopped")

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
