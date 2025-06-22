package port

import (
	"context"

	"github.com/axmz/go-port-service/internal/domain/port"
)

type InMem[T any] interface {
	Get(ctx context.Context, key string) (T, bool)
	GetAll(ctx context.Context) []T
	Put(ctx context.Context, key string, value T)
	Delete(ctx context.Context, key string) (T, bool)
	Len(ctx context.Context) int
}

type PortRepository struct {
	db InMem[*Port]
}

func New(db InMem[*Port]) *PortRepository {
	return &PortRepository{
		db: db,
	}
}

func (r PortRepository) GetPortByID(ctx context.Context, id string) (*port.Port, error) {
	portDb, exists := r.db.Get(ctx, id)
	if !exists {
		return nil, port.ErrNotFound
	}

	p, err := fromRepositoryToDomain(portDb)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r PortRepository) GetAllPorts(ctx context.Context) ([]*port.Port, error) {
	arr := r.db.GetAll(ctx)
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

func (r PortRepository) GetPortsCount(ctx context.Context) int {
	return r.db.Len(ctx)
}

func (r PortRepository) UploadPort(ctx context.Context, p *port.Port) error {
	portRepo, err := fromDomainToRepository(p)
	if err != nil {
		return err
	}
	r.db.Put(ctx, portRepo.ID, portRepo)
	return nil
}

func (r PortRepository) DeletePortByID(ctx context.Context, id string) (*port.Port, error) {
	portDb, exists := r.db.Delete(ctx, id)
	if !exists {
		return nil, port.ErrNotFound
	}

	p, err := fromRepositoryToDomain(portDb)
	if err != nil {
		return nil, err
	}

	return p, nil
}
