package port

import (
	"github.com/axmz/go-port-service/internal/domain/port"
)

type PortRepo interface {
	GetPort(id string) (*port.Port, error)
	GetPortsCount() int
	UploadPorts() error
}

type PortService struct {
	repo PortRepo
}

func NewPortService(r PortRepo) *PortService {
	return &PortService{
		repo: r,
	}
}

func (p *PortService) GetPort(id string) (*port.Port, error) {
	return p.repo.GetPort(id)
}

func (p *PortService) GetPortsCount() int {
	return p.repo.GetPortsCount()
}

func (p *PortService) UploadPorts() error {
	return p.repo.UploadPorts()
}
