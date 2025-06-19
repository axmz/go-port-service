package port

import (
	"github.com/axmz/go-port-service/internal/domain/port"
)

type PortRepo interface {
	GetPortByID(id string) (*port.Port, error)
	GetAllPorts() ([]*port.Port, error)
	GetPortsCount() int
	UploadPort(*port.Port) error
	DeletePortByID(id string) (*port.Port, error)
}

type PortService struct {
	repo PortRepo
}

func NewPortService(r PortRepo) *PortService {
	return &PortService{
		repo: r,
	}
}

func (p *PortService) GetPortByID(id string) (*port.Port, error) {
	return p.repo.GetPortByID(id)
}

func (p *PortService) GetAllPorts() ([]*port.Port, error) {
	return p.repo.GetAllPorts()
}

func (p *PortService) GetPortsCount() int {
	return p.repo.GetPortsCount()
}

func (p *PortService) UploadPort(port *port.Port) error {
	return p.repo.UploadPort(port)
}

func (p *PortService) DeletePortByID(id string) (*port.Port, error) {
	return p.repo.DeletePortByID(id)
}
