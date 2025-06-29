package app

import (
	"encoding/gob"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/axmz/go-port-service/internal/config"
	"github.com/axmz/go-port-service/internal/domain/user"
	"github.com/axmz/go-port-service/internal/logger"
	"github.com/axmz/go-port-service/internal/renderer"
	"github.com/axmz/go-port-service/pkg/inmem"

	portRepository "github.com/axmz/go-port-service/internal/repository/port"
	userRepository "github.com/axmz/go-port-service/internal/repository/user"

	portServices "github.com/axmz/go-port-service/internal/services/port"
	webAuthnServices "github.com/axmz/go-port-service/internal/services/webauthn"

	gqlHandler "github.com/axmz/go-port-service/internal/transport/graphql/handler"
	portHandlers "github.com/axmz/go-port-service/internal/transport/http/handlers/port"
	staticHandlers "github.com/axmz/go-port-service/internal/transport/http/handlers/static"
	webAuthnHandlers "github.com/axmz/go-port-service/internal/transport/http/handlers/webauthn"
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
