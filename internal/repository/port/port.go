package port

import (
	"errors"
	"time"

	"github.com/axmz/go-port-service/internal/domain/port"
)

type Port struct {
	ID          string
	Name        string
	Code        string
	City        string
	Country     string
	Alias       []string
	Regions     []string
	Coordinates []float64
	Province    string
	Timezone    string
	Unlocs      []string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type InMem interface {
	Get(key string) (any, bool)
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

func toDomain(p *Port) (*port.Port, error) {
	if p == nil {
		return nil, errors.New("store port is nil")
	}
	return port.NewPort(
		p.ID,
		p.Name,
		p.Code,
		p.City,
		p.Country,
		append([]string(nil), p.Alias...),
		append([]string(nil), p.Regions...),
		append([]float64(nil), p.Coordinates...),
		p.Province,
		p.Timezone,
		append([]string(nil), p.Unlocs...),
	)
}

func (r PortRepository) GetPort(id string) (*port.Port, error) {
	portDb, exists := r.db.Get(id)
	if !exists {
		return nil, port.ErrNotFound
	}
	p, err := toDomain(portDb.(*Port))
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r PortRepository) GetPortsCount() int {
	return 5
}

func (r PortRepository) UploadPort(p *port.Port) error {
	return nil
}
