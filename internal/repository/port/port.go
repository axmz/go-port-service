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
	return r.db.Len()
}

func toRepository(p *port.Port) (*Port, error) {
	return &Port{
		ID:          p.ID(),
		Name:        p.Name(),
		Code:        p.Code(),
		City:        p.City(),
		Country:     p.Country(),
		Alias:       append([]string(nil), p.Alias()...),
		Regions:     append([]string(nil), p.Regions()...),
		Coordinates: append([]float64(nil), p.Coordinates()...),
		Province:    p.Province(),
		Timezone:    p.Timezone(),
		Unlocs:      append([]string(nil), p.Unlocs()...),
	}, nil
}

func (r PortRepository) UploadPort(p *port.Port) error {
	portRepo, err := toRepository(p)
	if err != nil {
		return err
	}
	r.db.Put(portRepo.ID, portRepo)
	return nil
}
