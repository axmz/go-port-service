package static

import (
	"fmt"
	"net/http"
)

type Handlers struct {
}

func New() *Handlers {
	return &Handlers{}
}

func (h *Handlers) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "./static/index.html")
	} else {
		http.NotFound(w, r)
		return
	}
}

func (h *Handlers) PrivatePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/private/index.html")
}

func (h *Handlers) Metrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, metrics!")
}
