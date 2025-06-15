package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Port struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
}

type PortService interface {
	GetPort(id string) (Port, error)
	GetPortsCount() int
	UploadPorts()
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
	id := r.URL.Query().Get("id")
	log.Println("id", id)
	p, err := h.s.GetPort(id)
	if err == nil {

	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(p)
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
