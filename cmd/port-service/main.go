package main

import (
	"encoding/gob"
	"log"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/domain/user"
	"github.com/axmz/go-port-service/internal/logger"
	"github.com/axmz/go-port-service/pkg/graceful"
	"github.com/axmz/go-port-service/pkg/inmem"
	"github.com/go-webauthn/webauthn/webauthn"

	portRepository "github.com/axmz/go-port-service/internal/repository/port"
	userRepository "github.com/axmz/go-port-service/internal/repository/user"
	portServices "github.com/axmz/go-port-service/internal/services/port"
	webauthnServices "github.com/axmz/go-port-service/internal/services/webauthn"
	"github.com/axmz/go-port-service/internal/transport/http"
)

func run() error {
	cfg := config.MustLoad()

	logger.Setup(cfg.Env)

	slog.Info("Application starting", slog.String("env", cfg.Env))

	// DB
	portDB := inmem.New[*portRepository.Port]()
	userDB := inmem.New[*user.User]()

	// Repositories
	portRepo := portRepository.New(portDB)
	userRepo := userRepository.New(userDB)

	// Services
	portSvc := portServices.New(portRepo)
	webauthnSvc := webauthnServices.New(cfg, userRepo)

	gob.Register(webauthn.SessionData{})
	sessionManager := scs.New()

	// Server
	srv := http.Start(cfg, portSvc, webauthnSvc, sessionManager)

	<-graceful.Shutdown(cfg.GracefulTimeout, map[string]graceful.Operation{
		"port-database": portDB.Shutdown,
		"user-database": userDB.Shutdown,
		"http-server":   srv.Shutdown,
	})

	slog.Info("Application stopped")

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
