package zzcache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const shardCount = 512

var (
	errNotInitialized = errors.New("cache is not initialized")
	errNotInserted    = errors.New("unable to insert data")
	errNotFound       = errors.New("element is not found")
	errNoKey          = errors.New("key is not defined")
	errNoValue        = errors.New("value is not defined")
)

// Cache defines app objects
type Cache struct {
	mu     *sync.RWMutex
	hash   Hasher
	shards []*shard
	size   uint32
}

// New creates app
func New(size uint32, storeType string) (*Cache, error) {
	if size == 0 {
		return nil, fmt.Errorf("size of shards is not deifned")
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

// Set provides inserting to the cache
func (c *Cache) Set(key, value []byte, d time.Duration) error {
	hash := c.hash.Do(key)
	shardID := hash & c.size
	return c.set(shardID, key, value, d)
}

// Get provides getting data from the cache
func (c *Cache) Get(key []byte) ([]byte, error) {
	hash := c.hash.Do(key)
	shardID := hash & shardCount
	return c.get(shardID, key)
}

// Delete provides deletetign data from the cache
func (c *Cache) Delete(key []byte) error {
	c.mu.Lock()

	err := c.store.Delete(string(key))
	if err != nil {
		return err
	}

	defer c.mu.Unlock()
	return nil
}

// inner method for validating of the input data
// and append data to shards
func (c *Cache) set(shardID uint32, key, value []byte, d time.Duration) error {
	if c == nil {
		return errNotInitialized
	}
	if err := validateSet(key, value); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.shards[shardID].set(key, value); err != nil {
		return err
	}
	return nil
}

// validateSet provides validating of the input data
// before inserting to the cache
func validateSet(key, value []byte) error {
	if len(key) == 0 {
		return errNoKey
	}

	if len(value) == 0 {
		return errNoValue
	}

	return nil
}

func (c *Cache) get(shardID uint32, key []byte) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c == nil {
		return nil, errNotInitialized
	}
	return value, nil
}
