package port

import (
	"context"

	"github.com/axmz/go-port-service/internal/domain/port"
)

type PortRepo interface {
	GetPortByID(ctx context.Context, id string) (*port.Port, error)
	GetAllPorts(ctx context.Context) ([]*port.Port, error)
	GetPortsCount(ctx context.Context) int
	UploadPort(ctx context.Context, p *port.Port) error
	DeletePortByID(ctx context.Context, id string) (*port.Port, error)
}

type PortService struct {
	repo PortRepo
}

func New(r PortRepo) *PortService {
	return &PortService{
		repo: r,
	}
}

func (p *PortService) GetPortByID(ctx context.Context, id string) (*port.Port, error) {
	return p.repo.GetPortByID(ctx, id)
}

func (p *PortService) GetAllPorts(ctx context.Context) ([]*port.Port, error) {
	return p.repo.GetAllPorts(ctx)
}

func (p *PortService) GetPortsCount(ctx context.Context) int {
	return p.repo.GetPortsCount(ctx)
}

func (p *PortService) UploadPort(ctx context.Context, port *port.Port) error {
	return p.repo.UploadPort(ctx, port)
}

func (p *PortService) DeletePortByID(ctx context.Context, id string) (*port.Port, error) {
	return p.repo.DeletePortByID(ctx, id)
}
