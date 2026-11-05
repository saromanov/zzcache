package zzcache

import (
	"bytes"
	"encoding/gob"
	"time"
)

const maxKeySize = 65535

// shard provides implementation of the shard for cache
type shard struct {
	store Store
}

// entry represents d-s for inserting to linked list
type entry struct {
	Value []byte
	TTL   time.Time
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
		return ErrKeyTooLarge
	}

	// A zero duration means no expiration: the TTL is stored
	// as the zero time, which is never in the past.
	var ttl time.Time
	if d > 0 {
		ttl = time.Now().UTC().Add(d)
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(&entry{
		Value: value,
		TTL:   ttl,
	}); err != nil {
		return err
	}
	return s.store.Set(string(key), buf.Bytes())
}

func (s *shard) get(key []byte) ([]byte, error) {
	value, err := s.store.Get(string(key))
	if err != nil {
		return nil, ErrNotFound
	}
	if value == nil {
		return nil, ErrNotFound
	}

	w := entry{}

	//Use default gob decoder
	reader := bytes.NewReader(value)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&w); err != nil {
		return nil, err
	}
	// The zero TTL means the entry never expires.
	if !w.TTL.IsZero() && time.Now().UTC().After(w.TTL) {
		// Evict the expired entry so it does not leak memory.
		if delErr := s.store.Delete(string(key)); delErr != nil {
			return nil, ErrExpired
		}
		return nil, ErrExpired
	}
	return w.Value, nil
}

func (s *shard) del(key []byte) error {
	return s.store.Delete(string(key))
}
