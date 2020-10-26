package zzcache

import (
	"errors"
	"time"
)

const maxKeySize = 65535

var errKeyTooLarge = errors.New("key is too large")

// shard provides implementation of the shard for cache
type shard struct {
	store Store
}

// entry represents d-s for inserting to linked list
type entry struct {
	value []byte
	ttl   time.Time
}

// newShard creates a new shard
func newShard(s Store) *shard {
	return &shard{
		store: s,
	}
}

// set provides inserting of the data
func (s *shard) set(key, value []byte) error {
	if len(key) > maxKeySize {
		return errKeyTooLarge
	}

	err := s.store.Set(string(key), value)
	if err != nil {
		return err
	}
	return nil
}

func (s *shard) get(key []byte) ([]byte, error) {
	value, err := s.store.Get(string(key))
	if value == nil {
		return nil, errNotFound
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (s *shard) del(key []byte) error {
	return s.store.Delete(string(key))
}
