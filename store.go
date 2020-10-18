package zzcache

import "fmt"

type Store interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
}

type Map struct {
	data map[string][]byte
}

func (m *Map) Get(key string) ([]byte, error) {
	v, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key: %s is not found", key)
	}
	return v, nil
}

func (m *Map) Set(key string, value []byte) ([]byte, error) {
	m.data[key] = value
	return nil, nil
}
