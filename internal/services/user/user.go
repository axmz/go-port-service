package user

import (
	"context"

	"github.com/axmz/go-port-service/internal/domain/port"
)

type UserRepository interface {
	Get(ctx context.Context, id string) (*port.Port, error)
}

type Service struct {
	repo UserRepository
}

func New(r UserRepository) *Service {
	return &Service{
		repo: r,
	}
}

func (p *Service) Get(ctx context.Context, id string) (*port.Port, error) {
	return p.repo.Get(ctx, id)
}
