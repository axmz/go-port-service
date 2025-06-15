package gracefulshutdown

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Start(srv *http.Server) chan struct{} {
	shutdown := make(chan struct{})

	go func() {
		// Listen for SIGINT (Ctrl+C) or SIGTERM
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// Shutdown with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(shutdown)
	}()

	return shutdown
}
