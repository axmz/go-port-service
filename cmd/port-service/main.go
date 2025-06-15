package main

import (
	"log"
	"net/http"
	"time"

	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/gracefulshutdown"
	hh "github.com/axmz/go-port-service/internal/transport/http"
)

func run() error {
	cfg := config.LoadConfig()

	h := hh.NewHttpServer(nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HomePage)
	mux.HandleFunc("/metrics", h.Metrics)
	mux.HandleFunc("/port", h.GetPort)
	mux.HandleFunc("/count", h.GetPortsCount)
	mux.HandleFunc("POST /ports", h.UploadPorts)

	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
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
