package handlers

import (
	"errors"

	"github.com/axmz/go-port-service/internal/domain/port"
)

func (h *Handlers) toPortResponse(p *port.Port) PortResponse {
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

func toDomain(p *PortRequest) (*port.Port, error) {
	if p == nil {
		return nil, errors.New("store port is nil")
	}
	return port.NewPort(
		p.ID,
		p.Name,
		p.Code,
		p.City,
		p.Country,
		append([]string(nil), p.Alias...),
		append([]string(nil), p.Regions...),
		append([]float64(nil), p.Coordinates...),
		p.Province,
		p.Timezone,
		append([]string(nil), p.Unlocs...),
	)
}
