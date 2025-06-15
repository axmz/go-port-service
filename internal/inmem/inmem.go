package inmem

import "sync"

type InMemoryDB struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewInMemoryDB initializes the DB.
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[string]string),
	}
}

// Get returns the value for a key, or false if not found.
func (db *InMemoryDB) Get(key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, ok := db.data[key]
	return val, ok
}

// Put sets the value for a key.
func (db *InMemoryDB) Put(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}
