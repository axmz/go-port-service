package handlers

import (
	"fmt"
	"net/http"

	"github.com/axmz/go-port-service/internal/transport/http/handlers/port"
	"github.com/axmz/go-port-service/internal/transport/http/handlers/webauthn"
)

type Handlers struct {
	Ports    *port.Handlers
	WebAuthn *webauthn.Handlers
}

func New(p port.PortService, wa webauthn.WebAuthnService, session webauthn.SessionManager) *Handlers {
	return &Handlers{
		Ports:    port.New(p),
		WebAuthn: webauthn.New(wa, session),
	}
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
