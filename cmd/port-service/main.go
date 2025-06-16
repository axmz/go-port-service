package main

import (
	"log"
	"net/http"

	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/gracefulshutdown"
	"github.com/axmz/go-port-service/internal/router"

	db "github.com/axmz/go-port-service/internal/inmem"
	repository "github.com/axmz/go-port-service/internal/repository/port"
	services "github.com/axmz/go-port-service/internal/services/port"
	transport "github.com/axmz/go-port-service/internal/transport/http"
)

func run() error {
	cfg := config.LoadConfig()

	d := db.NewInMemoryDB()

	r := repository.NewPortRepository(d)

	s := services.NewPortService(r)

	h := transport.NewHttpServer(s)

	mux := router.Router(h)

	srv := &http.Server{
		Handler:      mux,
		Addr:         cfg.Port,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	shutdown := gracefulshutdown.Start(srv)

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-shutdown
	log.Println("Server stopped")
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
