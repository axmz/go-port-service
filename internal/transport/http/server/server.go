package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/axmz/go-port-service/internal/app"
	"github.com/axmz/go-port-service/internal/transport/http/middleware"
)

type Server struct {
	Router *http.Server
}

func NewServer(app *app.App) *Server {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.HandleFunc("/", app.Handlers.Page.Home)
	mux.Handle("/private", middleware.LoggedInMiddleware(app.Services.SessionManager, http.HandlerFunc(app.Handlers.Page.Private)))
	mux.Handle("/public", http.HandlerFunc(app.Handlers.Page.Private))
	mux.HandleFunc("/metrics", app.Handlers.Page.Metrics)

	mux.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", app.Handlers.GraphQLQuery)

	mux.HandleFunc("POST /api/ports", app.Handlers.Ports.Upload)
	mux.HandleFunc("GET /api/ports", app.Handlers.Ports.GetAll)
	mux.HandleFunc("GET /api/ports/{id}", app.Handlers.Ports.Get)
	mux.HandleFunc("GET /api/ports/count", app.Handlers.Ports.Count)
	mux.HandleFunc("PUT /api/ports/{id}", app.Handlers.Ports.UpdatePort)
	mux.HandleFunc("DELETE /api/ports/{id}", app.Handlers.Ports.Delete)

	mux.HandleFunc("POST /api/webauth/register/begin", app.Handlers.WebAuthn.BeginRegistration)
	mux.HandleFunc("POST /api/webauth/register/finish", app.Handlers.WebAuthn.FinishRegistration)
	mux.HandleFunc("POST /api/webauth/login/begin", app.Handlers.WebAuthn.BeginLogin)
	mux.HandleFunc("POST /api/webauth/login/finish", app.Handlers.WebAuthn.FinishLogin)
	mux.HandleFunc("POST /api/webauth/logout", app.Handlers.WebAuthn.Logout)

	handler :=
		middleware.Recoverer(
			app.Services.SessionManager.LoadAndSave(
				middleware.RequestID(
					middleware.Logger(mux))))

	r := &http.Server{
		Handler:      handler,
		Addr:         app.Config.HTTPServer.Port,
		IdleTimeout:  app.Config.HTTPServer.IdleTimeout,
		ReadTimeout:  app.Config.HTTPServer.ReadTimeout,
		WriteTimeout: app.Config.HTTPServer.WriteTimeout,
	}

	return &Server{
		Router: r,
	}
}

func (s *Server) Run() {
	slog.Info(fmt.Sprintf("Starting server on %s", s.Router.Addr), slog.String("op", "main.Server.Run"))
	if err := s.Router.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Router.Shutdown(ctx)
}
