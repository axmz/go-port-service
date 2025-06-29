package main

import (
	"log"
	"log/slog"

	"github.com/axmz/go-port-service/internal/app"
	"github.com/axmz/go-port-service/internal/transport/http/server"
	"github.com/axmz/go-port-service/pkg/graceful"
)

func start() error {
	app := app.SetupApp()

	server := server.NewServer(app)

	slog.Info("Application starting", slog.String("env", app.Config.Env))

	go func() {
		server.Run()
	}()

	<-graceful.Shutdown(app.Config.GracefulTimeout, map[string]graceful.Operation{
		"port-database": app.DB.Port.Shutdown,
		"user-database": app.DB.User.Shutdown,
		"http-server":   server.Shutdown,
	})

	slog.Info("Application stopped")

	return nil
}

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}
