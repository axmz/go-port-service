package handlers

// TODO: split handlers into separate folders

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/axmz/go-port-service/internal/domain/port"
	"github.com/axmz/go-port-service/internal/transport/http/response"
)

type PortService interface {
	GetPort(id string) (*port.Port, error)
	GetPortsCount() int
	UploadPort(*port.Port) error
}

type Handlers struct {
	service PortService
}

func NewHttpHandlers(s PortService) *Handlers {
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

func (h *Handlers) GetPort(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		response.BadRequest(w, "missing id")
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
		response.Ok(w, h.toPortResponse(p))
	}
}

func (h *Handlers) GetPortsCount(w http.ResponseWriter, r *http.Request) {
	c := h.service.GetPortsCount()
	response.Ok(w, c)
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
	r.Body = http.MaxBytesReader(w, r.Body, 50<<20) // TODO: take from config

	portCh := make(chan PortRequest)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	go readBody(r, portCh, errCh, doneCh)
	countPorts := 0

	for {
		select {
		case <-r.Context().Done():
			log.Println("request cancelled")
			return
		case err := <-errCh:
			log.Println(err)
			response.Err(w, http.StatusBadRequest, err.Error())
			return
		case p := <-portCh:
			countPorts++
			if portDomain, err := toDomain(&p); err != nil {
				response.Err(w, http.StatusBadRequest, err.Error())
			} else {
				h.service.UploadPort(portDomain)
			}
		case <-doneCh:
			log.Println("data processed successfully")
			response.Ok(w, countPorts)
			return
		}
	}
}

func (h *Handlers) Metrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, metrics!")
}
