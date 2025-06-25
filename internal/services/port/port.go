package port

import (
	"context"

	"github.com/axmz/go-port-service/internal/domain/port"
)

type PortRepository interface {
	Get(ctx context.Context, id string) (*port.Port, error)
	GetAll(ctx context.Context) ([]*port.Port, error)
	Count(ctx context.Context) int
	Upload(ctx context.Context, p *port.Port) error
	Delete(ctx context.Context, id string) (*port.Port, error)
}

type Service struct {
	port PortRepository
}

func New(r PortRepository) *Service {
	return &Service{
		port: r,
	}
}

func (p *Service) Get(ctx context.Context, id string) (*port.Port, error) {
	return p.port.Get(ctx, id)
}

func (p *Service) GetAll(ctx context.Context) ([]*port.Port, error) {
	return p.port.GetAll(ctx)
}

func (p *Service) Count(ctx context.Context) int {
	return p.port.Count(ctx)
}

func (p *Service) Upload(ctx context.Context, port *port.Port) error {
	return p.port.Upload(ctx, port)
}

func (p *Service) Delete(ctx context.Context, id string) (*port.Port, error) {
	return p.port.Delete(ctx, id)
}
