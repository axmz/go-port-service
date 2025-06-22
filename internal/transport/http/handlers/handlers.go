package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/axmz/go-port-service/internal/domain/port"
	"github.com/axmz/go-port-service/internal/transport/http/response"
)

type PortService interface {
	GetPortByID(ctx context.Context, id string) (*port.Port, error)
	DeletePortByID(ctx context.Context, id string) (*port.Port, error)
	GetAllPorts(ctx context.Context) ([]*port.Port, error)
	GetPortsCount(ctx context.Context) int
	UploadPort(ctx context.Context, p *port.Port) error
}

type Handlers struct {
	service PortService
}

func NewHTTPHandlers(s PortService) *Handlers {
	return &Handlers{
		service: s,
	}
}

func (h *Handlers) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "./static/index.html")
}

func (h *Handlers) GetAllPorts(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAllPorts(r.Context())
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	res := make([]PortResponse, 0, len(data))
	for _, v := range data {
		res = append(res, h.fromDomainToResponse(v))
	}
	response.OK(w, res)
}

func (h *Handlers) GetPortByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		response.BadRequest(w, "missing id")
		return
	}

	if p, err := h.service.GetPortByID(r.Context(), id); err != nil {
		if errors.Is(err, port.ErrNotFound) {
			response.NotFound(w)
		} else {
			response.InternalServerError(w, err)
		}
		return
	} else {
		response.OK(w, h.fromDomainToResponse(p))
	}
}

func (h *Handlers) GetPortsCount(w http.ResponseWriter, r *http.Request) {
	c := h.service.GetPortsCount(r.Context())
	response.OK(w, c)
}

func readBody(r *http.Request, portCh chan PortRequest, errCh chan error, doneCh chan struct{}) {
	defer close(portCh)
	defer close(errCh)
	defer close(doneCh)

	dec := json.NewDecoder(r.Body)

	if t, err := dec.Token(); err != nil || t != json.Delim('{') {
		errCh <- err
		return
	}

	for dec.More() {
		var id string
		if t, err := dec.Token(); err != nil {
			errCh <- err
			return
		} else {
			id = t.(string)
		}

		var p PortRequest
		if err := dec.Decode(&p); err != nil {
			errCh <- err
			return
		}

		p.ID = id
		portCh <- p
	}

	if _, err := dec.Token(); err != nil {
		errCh <- err
		return
	}

	doneCh <- struct{}{}
}

func (h *Handlers) UploadPorts(w http.ResponseWriter, r *http.Request) {
	const op = "transport.http.handlers.UploadPorts"
	r.Body = http.MaxBytesReader(w, r.Body, 50<<20) // TODO: move to config or middleware

	portCh := make(chan PortRequest)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	go readBody(r, portCh, errCh, doneCh)
	countPorts := 0

	for {
		select {
		case <-r.Context().Done():
			slog.Info("request cancelled", slog.String("op", op))
			return
		case err := <-errCh:
			slog.Info(err.Error(), slog.String("op", op))
			response.Err(w, http.StatusBadRequest, err.Error())
			return
		case p := <-portCh:
			countPorts++
			if portDomain, err := fromRequestToDomain(&p); err != nil {
				response.Err(w, http.StatusBadRequest, err.Error())
				return
			} else if err := h.service.UploadPort(r.Context(), portDomain); err != nil {
				response.BadRequest(w, err.Error())
				return
			}
		case <-doneCh:
			slog.Info("data processed successfully", slog.String("op", op))
			response.OK(w, countPorts)
			return
		}
	}
}

func (h *Handlers) UpdatePort(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		response.BadRequest(w, "missing id")
		return
	}

	if p, err := h.service.GetPortByID(r.Context(), id); err != nil {
		handleError(w, err)
		return
	} else {
		copy, _ := p.Copy()
		// TODO: impl granular update or complete replace
		copy.SetName("TEST")
		if err := h.service.UploadPort(r.Context(), copy); err != nil {
			response.BadRequest(w, err.Error())
			return
		} else {
			response.OK(w, h.fromDomainToResponse(copy))
		}
	}
}

func (h *Handlers) DeletePortByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		response.BadRequest(w, "missing id")
		return
	}

	if p, err := h.service.DeletePortByID(r.Context(), id); err != nil {
		handleError(w, err)
		return
	} else {
		response.OK(w, h.fromDomainToResponse(p))
	}
}

func (h *Handlers) Metrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, metrics!")
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, port.ErrNotFound):
		response.NotFound(w)
	default:
		response.InternalServerError(w, err)
	}
}
