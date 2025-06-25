package port

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/axmz/go-port-service/internal/domain/port"
	"github.com/axmz/go-port-service/internal/transport/http/response"
)

type PortService interface {
	Get(ctx context.Context, id string) (*port.Port, error)
	Delete(ctx context.Context, id string) (*port.Port, error)
	GetAll(ctx context.Context) ([]*port.Port, error)
	Count(ctx context.Context) int
	Upload(ctx context.Context, p *port.Port) error
}

type Handlers struct {
	port PortService
}

func New(s PortService) *Handlers {
	return &Handlers{
		port: s,
	}
}

func (h *Handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.port.GetAll(r.Context())
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	res := make([]Response, 0, len(data))
	for _, v := range data {
		res = append(res, h.fromDomainToResponse(v))
	}
	response.OK(w, res)
}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		response.BadRequest(w, "missing id")
		return
	}

	if p, err := h.port.Get(r.Context(), id); err != nil {
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

func (h *Handlers) Count(w http.ResponseWriter, r *http.Request) {
	c := h.port.Count(r.Context())
	response.OK(w, c)
}

func readBody(r *http.Request, portCh chan Request, errCh chan error, doneCh chan struct{}) {
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

		var p Request
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

func (h *Handlers) Upload(w http.ResponseWriter, r *http.Request) {
	const op = "transport.http.handlers.Upload"
	r.Body = http.MaxBytesReader(w, r.Body, 50<<20) // TODO: move to config or middleware

	portCh := make(chan Request)
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
			} else if err := h.port.Upload(r.Context(), portDomain); err != nil {
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

	if p, err := h.port.Get(r.Context(), id); err != nil {
		handleError(w, err)
		return
	} else {
		copy, _ := p.Copy()
		// TODO: impl granular update or complete replace
		copy.SetName("TEST")
		if err := h.port.Upload(r.Context(), copy); err != nil {
			response.BadRequest(w, err.Error())
			return
		} else {
			response.OK(w, h.fromDomainToResponse(copy))
		}
	}
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		response.BadRequest(w, "missing id")
		return
	}

	if p, err := h.port.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	} else {
		response.OK(w, h.fromDomainToResponse(p))
	}
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, port.ErrNotFound):
		response.NotFound(w)
	default:
		response.InternalServerError(w, err)
	}
}
