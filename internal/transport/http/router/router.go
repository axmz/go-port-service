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
	mux.HandleFunc("/port", h.GetPort)
	mux.HandleFunc("/count", h.GetPortsCount)
	mux.HandleFunc("GET /ports", h.GetPortsCount)
	mux.HandleFunc("POST /ports", h.UploadPorts)
	return middleware.LoggingMiddleware(mux)
}
