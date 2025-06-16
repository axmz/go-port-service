package port

import (
	"github.com/axmz/go-port-service/internal/domain/port"
	"github.com/axmz/go-port-service/internal/transport/http"
)

type PortRepo interface {
	GetPort(id string) (http.PortResponse, error)
	GetPortsCount() int
	UploadPorts()
}

type PortService struct {
	repo PortRepo
}

func NewPortService(r PortRepo) *PortService {
	return &PortService{
		repo: r,
	}
}

func (p *PortService) GetPort(id string) (port.Port, error) {
	return port.Port{}, nil
}

func (p *PortService) GetPortsCount() int {

	return 10
}

func (p *PortService) UploadPorts() {

}
