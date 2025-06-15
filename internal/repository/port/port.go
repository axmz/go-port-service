package port

import "github.com/axmz/go-port-service/internal/transport/http"

type InMem interface {
	Get(key string) (string, bool)
	Put(key, value string)
}

type PortRepository struct {
	db InMem
}

func NewPortRepository(db InMem) *PortRepository {
	return &PortRepository{
		db: db,
	}
}

func (r PortRepository) GetPort(id string) (http.Port, error) {
	if _, exists := r.db.Get(id); exists {
		return http.Port{}, nil
	}
	return http.Port{}, nil
}

func (r PortRepository) GetPortsCount() int {
	return 5
}

func (r PortRepository) UploadPorts() {
}
