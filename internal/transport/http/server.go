package http

import (
	"log"
	"net/http"

	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/services/port"
	"github.com/axmz/go-port-service/internal/transport/http/handlers"
	"github.com/axmz/go-port-service/internal/transport/http/router"
)

func StartServer(cfg *config.Config, s *port.PortService) *http.Server {
	h := handlers.NewHttpHandlers(s)
	mux := router.Router(h)

	srv := &http.Server{
		Handler:      mux,
		Addr:         cfg.Port,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go func() {
		log.Printf("Starting server on %s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	return srv
}
