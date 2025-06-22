package port

import (
	"errors"

	"slices"

	"github.com/axmz/go-port-service/internal/domain/port"
)

func fromDomainToRepository(p *port.Port) (*Port, error) {
	return &Port{
		ID:          p.ID(),
		Name:        p.Name(),
		Code:        p.Code(),
		City:        p.City(),
		Country:     p.Country(),
		Alias:       append([]string(nil), p.Alias()...),
		Regions:     append([]string(nil), p.Regions()...),
		Coordinates: append([]float64(nil), p.Coordinates()...),
		Province:    p.Province(),
		Timezone:    p.Timezone(),
		Unlocs:      append([]string(nil), p.Unlocs()...),
	}, nil
}

func fromRepositoryToDomain(p *Port) (*port.Port, error) {
	if p == nil {
		return nil, errors.New("store port is nil")
	}
	return port.New(
		p.ID,
		p.Name,
		p.Code,
		p.City,
		p.Country,
		slices.Clone(p.Alias),
		slices.Clone(p.Regions),
		slices.Clone(p.Coordinates),
		p.Province,
		p.Timezone,
		slices.Clone(p.Unlocs),
	)
}
