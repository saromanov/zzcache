package zzcache

import (
	"encoding/binary"
	"sync"
	"time"
)

const maxKeySize = 65535

// ttlHeaderSize is the size of the 8-byte big-endian TTL header
// stored in front of every value.
const ttlHeaderSize = 8

// shard provides implementation of the shard for cache.
// Each shard has its own lock, so operations on different shards
// do not contend with each other.
type shard struct {
	mu    sync.RWMutex
	store Store
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

	// A zero duration means no expiration: the TTL header is zero.
	var ttl int64
	if d > 0 {
		ttl = time.Now().UTC().Add(d).UnixNano()
	}

	// buf is a fresh allocation, so the stored value does not alias
	// the caller's slice.
	buf := make([]byte, ttlHeaderSize+len(value))
	binary.BigEndian.PutUint64(buf[:ttlHeaderSize], uint64(ttl))
	copy(buf[ttlHeaderSize:], value)

	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store.Set(string(key), buf)
}

func (s *shard) get(key []byte) ([]byte, error) {
	s.mu.RLock()
	buf, err := s.store.Get(string(key))
	s.mu.RUnlock()
	if err != nil {
		return nil, ErrNotFound
	}
	if len(buf) < ttlHeaderSize {
		return nil, ErrCorrupted
	}

	ttl := int64(binary.BigEndian.Uint64(buf[:ttlHeaderSize]))
	// A zero TTL means the entry never expires.
	if ttl != 0 && time.Now().UTC().UnixNano() > ttl {
		// Evict the expired entry so it does not leak memory.
		s.mu.Lock()
		_ = s.store.Delete(string(key))
		s.mu.Unlock()
		return nil, ErrExpired
	}

	// Copy the value out: the returned slice must not alias
	// the stored buffer.
	value := buf[ttlHeaderSize:]
	out := make([]byte, len(value))
	copy(out, value)
	return out, nil
}

func (s *shard) del(key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store.Delete(string(key))
}
