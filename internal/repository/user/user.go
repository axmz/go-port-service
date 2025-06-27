package user

import (
	"context"

	"github.com/axmz/go-port-service/internal/domain/user"
)

type InMem[T any] interface {
	Get(ctx context.Context, key string) (T, bool)
	Put(ctx context.Context, key string, value T)
	Delete(ctx context.Context, key string) (T, bool)
}

type Repository struct {
	db InMem[*user.User]
}

func New(db InMem[*user.User]) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) Get(ctx context.Context, id string) (*user.User, error) {
	u, exists := r.db.Get(ctx, id)
	if !exists {
		return nil, user.ErrNotFound
	}

	return u, nil
}

func (r Repository) Put(ctx context.Context, u *user.User) (*user.User, error) {
	r.db.Put(ctx, string(u.ID), u)
	return u, nil
}

func (r Repository) Delete(ctx context.Context, id string) (*user.User, error) {
	u, exists := r.db.Delete(ctx, id)
	if !exists {
		return nil, user.ErrNotFound
	}

	return u, nil
}
