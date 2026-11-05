package zzcache

import (
	"errors"
	"sync"
	"time"
)

// Sentinel errors returned by the cache. Use errors.Is to check for them.
var (
	ErrNotInitialized = errors.New("cache is not initialized")
	ErrNotFound       = errors.New("element is not found")
	ErrExpired        = errors.New("data has expired")
	ErrNoKey          = errors.New("key is not defined")
	ErrNoValue        = errors.New("value is not defined")
	ErrKeyTooLarge    = errors.New("key is too large")
	ErrNoShards       = errors.New("number of shards is not defined")
)

// Cache defines app objects
type Cache struct {
	mu     *sync.RWMutex
	hash   Hasher
	shards []*shard
	size   uint32
}

// New creates app. If d passed to Set is zero, the entry never expires.
func New(size uint32, storeType string) (*Cache, error) {
	if size == 0 {
		return nil, ErrNoShards
	}
	store := NewMap()
	if storeType == "radix" {
		store = NewRadix()
	}
	return &Cache{
		mu:     &sync.RWMutex{},
		hash:   new(CRC32),
		shards: initShards(size, store),
		size:   size,
	}, nil
}

func initShards(size uint32, st Store) []*shard {
	s := make([]*shard, size)
	for i := uint32(0); i < size; i++ {
		s[i] = newShard(st)
	}
	return s
}

// Set provides inserting to the cache. A zero duration means no expiration.
func (c *Cache) Set(key, value []byte, d time.Duration) error {
	if c == nil {
		return ErrNotInitialized
	}
	return c.set(c.getShard(key), key, value, d)
}

// Get provides getting data from the cache
func (c *Cache) Get(key []byte) ([]byte, error) {
	if c == nil {
		return nil, ErrNotInitialized
	}
	return c.get(c.getShard(key), key)
}

// Delete provides deleting data from the cache.
// Deleting a missing key is not an error.
func (c *Cache) Delete(key []byte) error {
	if c == nil {
		return ErrNotInitialized
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.shards[c.getShard(key)].del(key); err != nil {
		return err
	}
	return nil
}

// inner method for validating of the input data
// and append data to shards
func (c *Cache) set(shardID uint32, key, value []byte, d time.Duration) error {
	if err := validateSet(key, value); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.shards[shardID].set(key, value, d)
}

func (c *Cache) get(shardID uint32, key []byte) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.shards[shardID].get(key)
}

func (c *Cache) getShard(key []byte) uint32 {
	hash := c.hash.Do(key)
	return hash % c.size
}

// validateSet provides validating of the input data
// before inserting to the cache
func validateSet(key, value []byte) error {
	if len(key) == 0 {
		return ErrNoKey
	}

	if len(value) == 0 {
		return ErrNoValue
	}

	return nil
}
