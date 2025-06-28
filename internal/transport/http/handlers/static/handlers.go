package static

import (
	"fmt"
	"net/http"

	"github.com/axmz/go-port-service/internal/renderer"
)

type Handlers struct {
	TemplateRenderer *renderer.TemplateRenderer
}

func New(r *renderer.TemplateRenderer) *Handlers {
	return &Handlers{
		TemplateRenderer: r,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.TemplateRenderer.Render(w, "home.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handlers) Private(w http.ResponseWriter, r *http.Request) {
	err := h.TemplateRenderer.Render(w, "private.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handlers) Metrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, metrics!")
}
