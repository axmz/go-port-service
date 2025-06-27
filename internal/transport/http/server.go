package http

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/alexedwards/scs/v2"
	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/services/port"
	graphql "github.com/axmz/go-port-service/internal/transport/graphql"
	"github.com/axmz/go-port-service/internal/transport/http/handlers"
	"github.com/axmz/go-port-service/internal/transport/http/handlers/webauthn"
	"github.com/axmz/go-port-service/internal/transport/http/router"
	"github.com/vektah/gqlparser/v2/ast"
)

func Start(
	cfg *config.Config,
	portSvc *port.Service,
	webauthnSvc webauthn.WebAuthnService,
	session webauthn.SessionManager,
) *http.Server {
	const op = "transport.http.server.Start"
	h := handlers.New(portSvc, webauthnSvc, session)

	gqlsrv := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{
		PortService: portSvc,
	}}))
	gqlsrv.AddTransport(transport.Options{})
	gqlsrv.AddTransport(transport.GET{})
	gqlsrv.AddTransport(transport.POST{})
	gqlsrv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	gqlsrv.Use(extension.Introspection{})
	gqlsrv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	mux := router.Router(h, gqlsrv, (session).(*scs.SessionManager))

	srv := &http.Server{
		Handler:      mux,
		Addr:         cfg.HTTPServer.Port,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
	}

	go func() {
		slog.Info(fmt.Sprintf("Starting server on %s", cfg.HTTPServer.Port), slog.String("op", op))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	return srv
}
