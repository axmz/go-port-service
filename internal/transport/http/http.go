package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/axmz/go-port-service/internal/domain/port"
	"github.com/axmz/go-port-service/internal/transport/http/response"
)

type PortRequest = PortResponse
type PortResponse struct {
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
	GetPort(id string) (port.Port, error)
	GetPortsCount() int
	UploadPorts()
}

type HttpServer struct {
	service PortService
}

func NewHttpServer(s PortService) *HttpServer {
	return &HttpServer{
		service: s,
	}
}

func (h *HttpServer) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, "Hello, world!")
}

func (h *HttpServer) toPortResponse(p port.Port) PortResponse {
	r := PortResponse{
		ID:          p.ID(),
		Name:        p.Name(),
		City:        p.City(),
		Country:     p.Country(),
		Alias:       p.Alias(),
		Regions:     p.Regions(),
		Coordinates: p.Coordinates(),
		Province:    p.Province(),
		Timezone:    p.Timezone(),
		Unlocs:      p.Unlocs(),
	}

	return r
}

func (h *HttpServer) GetPort(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		response.MissingID(w)
		return
	}

	if p, err := h.service.GetPort(id); err != nil {
		if errors.Is(err, port.ErrNotFound) {
			response.NotFound(w)
		} else {
			response.InternalServerError(w, err)
		}
		return
	} else {
		response.JSONOK(w, h.toPortResponse(p))
	}
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
