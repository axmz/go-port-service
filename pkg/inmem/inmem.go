package inmem

import (
	"context"
	"sync"
)

type InMemoryDB struct {
	data map[string]any
	mu   sync.RWMutex
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[string]any),
	}
}

func (db *InMemoryDB) Get(key string) (any, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, ok := db.data[key]
	return val, ok
}

func (db *InMemoryDB) GetAll() []string {
	res := make([]string, 0, len(db.data))

	db.mu.RLock()
	defer db.mu.RUnlock()

	for k := range db.data {
		res = append(res, k)
	}

	return res
}

func (db *InMemoryDB) Put(key string, value any) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}

func (db *InMemoryDB) Delete(key string) (any, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	temp, ok := db.data[key]
	delete(db.data, key)
	return temp, ok
}

func (db *InMemoryDB) Len() int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return len(db.data)
}

func (db *InMemoryDB) Shutdown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
