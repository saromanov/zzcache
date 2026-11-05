package zzcache

import (
	"encoding/binary"
	"sync/atomic"
	"time"
)

// RangingStore provides iteration over stored entries.
// Stock stores implement it; a custom store without Range
// supports Count, Keys, Flush and the cleanup janitor as no-ops.
type RangingStore interface {
	// Range calls fn for each stored entry until fn returns false.
	// fn must not modify the store.
	Range(fn func(key string, value []byte) bool)
}

// Stats provides the accumulated cache statistics.
type Stats struct {
	Hits   uint64
	Misses uint64
}

// Count returns the number of stored entries, including expired
// entries that were not evicted yet.
// Returns 0 for stores that do not support iteration.
func (c *Cache) Count() int {
	if c == nil {
		return 0
	}
	n := 0
	for _, s := range c.shards {
		n += s.count()
	}
	return n
}

// Keys returns the stored keys. The order is not defined.
// Expired entries that were not evicted yet are included.
// Returns nil for stores that do not support iteration.
func (c *Cache) Keys() [][]byte {
	if c == nil {
		return nil
	}
	var out [][]byte
	for _, s := range c.shards {
		out = append(out, s.keys()...)
	}
	return out
}

// Flush removes all stored entries.
// No-op for stores that do not support iteration.
func (c *Cache) Flush() {
	if c == nil {
		return
	}
	for _, s := range c.shards {
		s.flush()
	}
}

// Stats returns the accumulated hits and misses.
func (c *Cache) Stats() Stats {
	if c == nil {
		return Stats{}
	}
	var st Stats
	for _, s := range c.shards {
		st.Hits += s.hits()
		st.Misses += s.misses()
	}
	return st
}

// StartCleanup starts a background goroutine that periodically
// evicts expired entries. Returns a stop function that terminates
// the goroutine. Calling it more than once is not supported.
func (c *Cache) StartCleanup(interval time.Duration) (stop func()) {
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.cleanup()
			case <-done:
				return
			}
		}
	}()
	return func() { close(done) }
}

func (c *Cache) cleanup() {
	for _, s := range c.shards {
		s.cleanupExpired()
	}
}

func (s *shard) count() int {
	rs, ok := s.store.(RangingStore)
	if !ok {
		return 0
	}
	n := 0
	s.mu.RLock()
	rs.Range(func(key string, value []byte) bool {
		n++
		return true
	})
	s.mu.RUnlock()
	return n
}

func (s *shard) keys() [][]byte {
	rs, ok := s.store.(RangingStore)
	if !ok {
		return nil
	}
	var out [][]byte
	s.mu.RLock()
	rs.Range(func(key string, value []byte) bool {
		out = append(out, []byte(key))
		return true
	})
	s.mu.RUnlock()
	return out
}

func (s *shard) flush() {
	rs, ok := s.store.(RangingStore)
	if !ok {
		return
	}
	keys := make([]string, 0, 16)
	s.mu.RLock()
	rs.Range(func(key string, value []byte) bool {
		keys = append(keys, key)
		return true
	})
	s.mu.RUnlock()
	for _, key := range keys {
		s.mu.Lock()
		_ = s.store.Delete(key)
		s.mu.Unlock()
	}
}

func (s *shard) cleanupExpired() int {
	rs, ok := s.store.(RangingStore)
	if !ok {
		return 0
	}
	now := time.Now().UTC().UnixNano()
	expired := make([]string, 0)
	s.mu.RLock()
	rs.Range(func(key string, value []byte) bool {
		if len(value) >= ttlHeaderSize {
			ttl := int64(binary.BigEndian.Uint64(value[:ttlHeaderSize]))
			if ttl != 0 && now > ttl {
				expired = append(expired, key)
			}
		}
		return true
	})
	s.mu.RUnlock()
	for _, key := range expired {
		s.mu.Lock()
		_ = s.store.Delete(key)
		s.mu.Unlock()
	}
	return len(expired)
}

func (s *shard) hits() uint64 {
	return atomic.LoadUint64(&s.hitCount)
}

func (s *shard) misses() uint64 {
	return atomic.LoadUint64(&s.missCount)
}
