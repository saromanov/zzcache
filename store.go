package zzcache

import (
	"fmt"

	"github.com/armon/go-radix"
)

type Store interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type Map struct {
	data map[string][]byte
}

func NewMap() Store {
	return &Map{
		data: map[string][]byte{},
	}
}

func (m *Map) Get(key string) ([]byte, error) {
	v, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key: %s is not found", key)
	}
	return v, nil
}

func (m *Map) Set(key string, value []byte) error {
	m.data[key] = value
	return nil
}

func (m *Map) Delete(key string) error {
	delete(m.data, key)
	return nil
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
		return nil, fmt.Errorf("unable to get value")
	}
	return res.([]byte), nil
}

func (r *Radix) Set(key string, value []byte) error {
	r.tree.Insert(key, value)
	return nil
}

func (r *Radix) Delete(key string) error {
	_, ok := r.tree.Delete(key)
	if !ok {
		return fmt.Errorf("unable to delete key: %s", key)
	}
	return nil
}
