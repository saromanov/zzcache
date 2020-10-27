package zzcache

import (
	"bytes"
	"encoding/gob"
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
func (s *shard) set(key, value []byte, d time.Duration) error {
	if len(key) > maxKeySize {
		return errKeyTooLarge
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(&entry{
		value: value,
		ttl:   time.Now().UTC().Add(d),
	}); err != nil {
		return err
	}
	err := s.store.Set(string(key), buf.Bytes())
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
