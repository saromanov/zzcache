package zzcache

import (
	"sync"

	"github.com/armon/go-radix"
)

// Store provides a pluggable backend for a single shard.
// Each shard owns its own Store instance; implementations
// do not need to be safe for concurrent use by themselves,
// but all stock implementations are.
type Store interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type Map struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewMap() Store {
	return &Map{
		data: map[string][]byte{},
	}
}

func (m *Map) Get(key string) ([]byte, error) {
	m.mu.RLock()
	v, ok := m.data[key]
	m.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func (m *Map) Set(key string, value []byte) error {
	m.mu.Lock()
	m.data[key] = value
	m.mu.Unlock()
	return nil
}

func (m *Map) Delete(key string) error {
	m.mu.Lock()
	delete(m.data, key)
	m.mu.Unlock()
	return nil
}

// Range calls fn for each entry under a read lock.
// fn must not call back into the same Map.
func (m *Map) Range(fn func(key string, value []byte) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !fn(k, v) {
			return
		}
	}
}

type Radix struct {
	tree *radix.Tree
}

func NewRadix() Store {
	return &Radix{
		tree: radix.New(),
	}
}

func (r *Radix) Get(key string) ([]byte, error) {
	res, ok := r.tree.Get(key)
	if !ok {
		return nil, ErrNotFound
	}
	return res.([]byte), nil
}

func (r *Radix) Set(key string, value []byte) error {
	r.tree.Insert(key, value)
	return nil
}

func (r *Radix) Delete(key string) error {
	// Deleting a missing key is not an error.
	r.tree.Delete(key)
	return nil
}

// Range calls fn for each entry until fn returns false.
// fn must not call back into the same Radix.
//
// Note: go-radix Walk uses inverted semantics (true stops the
// walk), so the callback is inverted here.
func (r *Radix) Range(fn func(key string, value []byte) bool) {
	r.tree.Walk(func(k string, v interface{}) bool {
		return !fn(k, v.([]byte))
	})
}
