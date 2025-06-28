package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alexedwards/scs/v2"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/domain/user"
	"github.com/axmz/go-port-service/internal/logger"
	"github.com/axmz/go-port-service/internal/renderer"
	"github.com/axmz/go-port-service/pkg/graceful"
	"github.com/axmz/go-port-service/pkg/inmem"

	portRepository "github.com/axmz/go-port-service/internal/repository/port"
	userRepository "github.com/axmz/go-port-service/internal/repository/user"

	portServices "github.com/axmz/go-port-service/internal/services/port"
	webAuthnServices "github.com/axmz/go-port-service/internal/services/webauthn"

	gqlHandler "github.com/axmz/go-port-service/internal/transport/graphql/handler"
	portHandlers "github.com/axmz/go-port-service/internal/transport/http/handlers/port"
	staticHandlers "github.com/axmz/go-port-service/internal/transport/http/handlers/static"
	webAuthnHandlers "github.com/axmz/go-port-service/internal/transport/http/handlers/webauthn"

	"github.com/axmz/go-port-service/internal/transport/http/middleware"
)

type App struct {
	Config *config.Config
	Log    *slog.Logger
	DB     struct {
		Port *inmem.InMemoryDB[*portRepository.Port]
		User *inmem.InMemoryDB[*user.User]
	}
	Repos struct {
		Port *portRepository.Repository
		User *userRepository.Repository
	}
	Services struct {
		Port           *portServices.Service
		WebAuthn       *webAuthnServices.Service
		SessionManager *scs.SessionManager
	}
	Handlers struct {
		Page         *staticHandlers.Handlers
		Ports        *portHandlers.Handlers
		WebAuthn     *webAuthnHandlers.Handlers
		GraphQLQuery *gqlHandler.GraphQLHandler
	}
	TemplateRenderer *renderer.TemplateRenderer
}

func SetupApp() *App {
	app := &App{}

	// Config
	app.Config = config.MustLoad()

	// Logger
	app.Log = logger.Setup(app.Config.Env) // TODO: use logger in the app instead of slog?

	// DB
	app.DB.Port = inmem.New[*portRepository.Port]()
	app.DB.User = inmem.New[*user.User]()

	// Repositories
	app.Repos.Port = portRepository.New(app.DB.Port)
	app.Repos.User = userRepository.New(app.DB.User)

	// Services
	app.Services.Port = portServices.New(app.Repos.Port)
	app.Services.WebAuthn = webAuthnServices.New(app.Config, app.Repos.User)
	app.Services.SessionManager = scs.New()
	gob.Register(webauthn.SessionData{})

	// Renderer
	app.TemplateRenderer = renderer.NewTemplateRenderer()

	// Handlers
	app.Handlers.Page = staticHandlers.New(app.TemplateRenderer)
	app.Handlers.Ports = portHandlers.New(app.Services.Port)
	app.Handlers.WebAuthn = webAuthnHandlers.New(app.Services.WebAuthn, app.Services.SessionManager)
	app.Handlers.GraphQLQuery = gqlHandler.InitGql(app.Services.Port)

	return app
}

type Server struct {
	router *http.Server
}

func NewServer(app *App) *Server {
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
		router: r,
	}
}

func (s *Server) Run() {
	slog.Info(fmt.Sprintf("Starting server on %s", s.router.Addr), slog.String("op", "main.Server.Run"))
	if err := s.router.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

func start() error {
	app := SetupApp()
	server := NewServer(app)
	slog.Info("Application starting", slog.String("env", app.Config.Env))

	go func() {
		server.Run()
	}()

	<-graceful.Shutdown(app.Config.GracefulTimeout, map[string]graceful.Operation{
		"port-database": app.DB.Port.Shutdown,
		"user-database": app.DB.User.Shutdown,
		"http-server":   server.router.Shutdown,
	})

	slog.Info("Application stopped")

	return nil
}

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}
