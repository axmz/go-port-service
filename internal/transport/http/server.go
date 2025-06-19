package http

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/services/port"
	graphql "github.com/axmz/go-port-service/internal/transport/graphql"
	"github.com/axmz/go-port-service/internal/transport/http/handlers"
	"github.com/axmz/go-port-service/internal/transport/http/router"
	"github.com/vektah/gqlparser/v2/ast"
)

func StartServer(cfg *config.Config, s *port.PortService) *http.Server {
	h := handlers.NewHTTPHandlers(s)

	gqlsrv := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{
		PortService: s,
	}}))
	gqlsrv.AddTransport(transport.Options{})
	gqlsrv.AddTransport(transport.GET{})
	gqlsrv.AddTransport(transport.POST{})
	gqlsrv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	gqlsrv.Use(extension.Introspection{})
	gqlsrv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	mux := router.Router(h, gqlsrv)

	srv := &http.Server{
		Handler:      mux,
		Addr:         cfg.HTTPServer.Port,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
	}

	go func() {
		log.Printf("Starting server on %s\n", cfg.HTTPServer.Port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	return srv
}
