package port

import (
	"github.com/axmz/go-port-service/internal/domain/port"
)

type PortRepo interface {
	GetPortById(id string) (*port.Port, error)
	GetAllPorts() ([]string, error)
	GetPortsCount() int
	UploadPort(*port.Port) error
	DeletePortById(id string) (*port.Port, error)
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
	return p.repo.GetPortById(id)
}

func (p *PortService) GetAllPorts() ([]string, error) {
	return p.repo.GetAllPorts()
}

func (p *PortService) GetPortsCount() int {
	return p.repo.GetPortsCount()
}

func (p *PortService) UploadPort(port *port.Port) error {
	return p.repo.UploadPort(port)
}

func (p *PortService) DeletePortById(id string) (*port.Port, error) {
	return p.repo.DeletePortById(id)
}
