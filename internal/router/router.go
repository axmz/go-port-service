package router

import (
	"net/http"

	"github.com/axmz/go-port-service/internal/middleware"
	transport "github.com/axmz/go-port-service/internal/transport/http"
)

func Router(h *transport.HttpServer) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HomePage)
	mux.HandleFunc("/metrics", h.Metrics)
	mux.HandleFunc("/port", h.GetPort)
	mux.HandleFunc("/count", h.GetPortsCount)
	// TODO: add get ports
	mux.HandleFunc("POST /ports", h.UploadPorts)
	return middleware.LoggingMiddleware(mux)
}
