package router

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alexedwards/scs/v2"
	"github.com/axmz/go-port-service/internal/transport/http/handlers"
	"github.com/axmz/go-port-service/internal/transport/http/middleware"
)

func Router(h *handlers.Handlers, gqlsrv *handler.Server, session *scs.SessionManager) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.HandleFunc("/", h.HomePage)
	privateHandler := middleware.LoggedInMiddleware(session, http.HandlerFunc(h.PrivatePage))
	mux.Handle("/private", privateHandler)
	// mux.HandleFunc("/private", h.PrivatePage)
	mux.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", gqlsrv)
	mux.HandleFunc("/metrics", h.Metrics)

	mux.HandleFunc("POST /api/ports", h.Ports.Upload)
	mux.HandleFunc("GET /api/ports", h.Ports.GetAll)
	mux.HandleFunc("GET /api/ports/{id}", h.Ports.Get)
	mux.HandleFunc("GET /api/ports/count", h.Ports.Count)
	mux.HandleFunc("PUT /api/ports/{id}", h.Ports.UpdatePort)
	mux.HandleFunc("DELETE /api/ports/{id}", h.Ports.Delete)

	mux.HandleFunc("POST /api/webauth/register/begin", h.WebAuthn.BeginRegistration)
	mux.HandleFunc("POST /api/webauth/register/finish", h.WebAuthn.FinishRegistration)
	mux.HandleFunc("POST /api/webauth/login/begin", h.WebAuthn.BeginLogin)
	mux.HandleFunc("POST /api/webauth/login/finish", h.WebAuthn.FinishLogin)

	return middleware.Recoverer(session.LoadAndSave(middleware.RequestID(middleware.Logger(mux))))
}
