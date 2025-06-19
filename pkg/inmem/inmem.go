package inmem

import (
	"context"
	"sync"
)

type InMemoryDB[T any] struct {
	data map[string]T
	mu   sync.RWMutex
}

func NewInMemoryDB[T any]() *InMemoryDB[T] {
	return &InMemoryDB[T]{
		data: make(map[string]T),
	}
}

func (db *InMemoryDB[T]) Get(key string) (T, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, ok := db.data[key]
	return val, ok
}

func (db *InMemoryDB[T]) GetAll() []T {
	res := make([]T, 0, len(db.data))

	db.mu.RLock()
	defer db.mu.RUnlock()

	for _, v := range db.data {
		res = append(res, v)
	}

	return res
}

func (db *InMemoryDB[T]) Put(key string, value T) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}

func (db *InMemoryDB[T]) Delete(key string) (T, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	temp, ok := db.data[key]
	delete(db.data, key)
	return temp, ok
}

func (db *InMemoryDB[T]) Len() int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return len(db.data)
}

func (db *InMemoryDB[T]) Shutdown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
