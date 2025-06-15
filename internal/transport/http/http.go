package http

import (
	"fmt"
	"net/http"
)

type PortService interface {
}

type HttpServer struct {
	s PortService
}

func NewHttpServer(s PortService) *HttpServer {
	return &HttpServer{
		s: s,
	}
}

func (h *HttpServer) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, "Hello, world!")
}

func (h *HttpServer) GetPort(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GetPort")
}

func (h *HttpServer) GetPortsCount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GetPortsCount")
}

func (h *HttpServer) UploadPorts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "UploadPorts")
}

func (h *HttpServer) Metrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, metrics!")
}
