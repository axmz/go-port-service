package router

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/axmz/go-port-service/internal/transport/http/handlers"
	"github.com/axmz/go-port-service/internal/transport/http/middleware"
)

func Router(h *handlers.Handlers, gqlsrv *handler.Server) http.Handler {
	fs := http.FileServer(http.Dir("../../static"))
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", h.HomePage)
	mux.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", gqlsrv)
	mux.HandleFunc("/metrics", h.Metrics)

	mux.HandleFunc("POST /api/ports", h.Ports.Upload)
	mux.HandleFunc("GET /api/ports", h.Ports.GetAll)
	mux.HandleFunc("GET /api/ports/{id}", h.Ports.Get)
	mux.HandleFunc("GET /api/ports/count", h.Ports.Count)
	mux.HandleFunc("PUT /api/ports/{id}", h.Ports.UpdatePort)
	mux.HandleFunc("DELETE /api/ports/{id}", h.Ports.Delete)

	return middleware.Recoverer(middleware.RequestID(middleware.Logger(mux)))
}
