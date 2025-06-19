package port

import (
	"github.com/axmz/go-port-service/internal/domain/port"
)

type InMem[T any] interface {
	Get(key string) (T, bool)
	GetAll() []T
	Put(key string, value T)
	Delete(key string) (T, bool)
	Len() int
}

type PortRepository struct {
	db InMem[*Port]
}

func NewPortRepository(db InMem[*Port]) *PortRepository {
	return &PortRepository{
		db: db,
	}
}

func (r PortRepository) GetPortByID(id string) (*port.Port, error) {
	portDb, exists := r.db.Get(id)
	if !exists {
		return nil, port.ErrNotFound
	}

	p, err := fromRepositoryToDomain(portDb)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r PortRepository) GetAllPorts() ([]*port.Port, error) {
	arr := r.db.GetAll()
	res := make([]*port.Port, 0, len(arr))

	for _, v := range arr {
		p, err := fromRepositoryToDomain(v)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
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

func (r PortRepository) DeletePortByID(id string) (*port.Port, error) {
	portDb, exists := r.db.Delete(id)
	if !exists {
		return nil, port.ErrNotFound
	}

	p, err := fromRepositoryToDomain(portDb)
	if err != nil {
		return nil, err
	}

	return p, nil
}
