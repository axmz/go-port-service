package router

import (
	"net/http"

	"github.com/axmz/go-port-service/internal/transport/http/handlers"
	"github.com/axmz/go-port-service/internal/transport/http/middleware"
)

func Router(h *handlers.Handlers) http.Handler {
	fs := http.FileServer(http.Dir("../../static"))
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", h.HomePage)
	mux.HandleFunc("/metrics", h.Metrics)

	// TODO: give same name in postman
	mux.HandleFunc("POST /api/ports", h.UploadPorts)
	mux.HandleFunc("GET /api/ports", h.GetAllPorts)
	mux.HandleFunc("GET /api/ports/{id}", h.GetPortById)
	mux.HandleFunc("GET /api/ports/count", h.GetPortsCount)
	mux.HandleFunc("PUT /api/ports/{id}", h.UpdatePort)
	mux.HandleFunc("DELETE /api/ports/{id}", h.DeletePortById)

	return middleware.LoggingMiddleware(mux)
}
