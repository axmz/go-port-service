package graph

import (
	"github.com/axmz/go-port-service/graph/model"
	"github.com/axmz/go-port-service/internal/domain/port"
)

func convertToGraphQLPort(p *port.Port) *model.Port {
	return &model.Port{
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
}
