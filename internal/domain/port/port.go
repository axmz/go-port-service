package port

import (
	"errors"
	"fmt"
	"slices"
)

var (
	ErrNotFound   = errors.New("port not found")
	ErrValidation = errors.New("validation error")
	ErrRequired   = fmt.Errorf("%w: value cannot be empty", ErrValidation)
)

type Port struct {
	id          string
	name        string
	code        string
	city        string
	country     string
	alias       []string
	regions     []string
	coordinates []float64
	province    string
	timezone    string
	unlocs      []string
}

func NewPort(
	id, name, code, city, country string,
	alias, regions []string,
	coords []float64,
	province, tz string,
	unlocs []string) (*Port, error) {

	if err := validate(id, name, city, country); err != nil {
		// TEST remove
		if errors.Is(err, ErrValidation) {
			fmt.Println("validation err")
		}
		if errors.Is(err, ErrRequired) {
			fmt.Println("required err")
		}
		return nil, err
	}

	return &Port{
		id:          id,
		name:        name,
		code:        code,
		city:        city,
		country:     country,
		alias:       alias,
		regions:     regions,
		coordinates: coords,
		province:    province,
		timezone:    tz,
		unlocs:      unlocs,
	}, nil
}

func (p *Port) ID() string {
	return p.id
}

func (p *Port) Name() string {
	return p.name
}

func (p *Port) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: port name is required", ErrRequired)
	}
	p.name = name
	return nil
}

func (p *Port) Code() string {
	return p.code
}

func (p *Port) City() string {
	return p.city
}

func (p *Port) Country() string {
	return p.country
}

func (p *Port) Alias() []string {
	return p.alias
}

func (p *Port) Regions() []string {
	return p.regions
}

func (p *Port) Coordinates() []float64 {
	return p.coordinates
}

func (p *Port) Province() string {
	return p.province
}

func (p *Port) Timezone() string {
	return p.timezone
}

func (p *Port) Unlocs() []string {
	return p.unlocs
}

func (p *Port) Copy() (*Port, error) {
	return NewPort(
		p.ID(),
		p.Name(),
		p.Code(),
		p.City(),
		p.Country(),
		slices.Clone(p.Alias()),
		slices.Clone(p.Regions()),
		slices.Clone(p.Coordinates()),
		p.Province(),
		p.Timezone(),
		slices.Clone(p.Unlocs()),
	)
}
