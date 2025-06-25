package handlers

import (
	"fmt"
	"net/http"

	"github.com/axmz/go-port-service/internal/transport/http/handlers/port"
)

type Handlers struct {
	Ports *port.Handlers
}

func New(p port.PortService) *Handlers {
	return &Handlers{
		Ports: port.New(p),
	}
}

func (h *Handlers) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "./static/index.html")
}

func (h *Handlers) Metrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, metrics!")
}
