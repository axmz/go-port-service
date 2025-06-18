package port

import (
	"github.com/axmz/go-port-service/internal/domain/port"
)

type InMem interface {
	Get(key string) (any, bool)
	Put(key string, value any)
	Len() int
}

type PortRepository struct {
	db InMem
}

func NewPortRepository(db InMem) *PortRepository {
	return &PortRepository{
		db: db,
	}
}

func (r PortRepository) GetPort(id string) (*port.Port, error) {
	portDb, exists := r.db.Get(id)
	if !exists {
		return nil, port.ErrNotFound
	}

	p, err := fromRepositoryToDomain(portDb.(*Port))
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r PortRepository) GetPortsCount() int {
	return r.db.Len()
}

func (r PortRepository) UploadPort(p *port.Port) error {
	portRepo, err := fromDomainToRepository(p)
	if err != nil {
		return err
	}
	r.db.Put(portRepo.ID, portRepo)
	return nil
}
