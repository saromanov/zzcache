package zzcache

import (
	"errors"
	"sync"
)

const shardCount = 512

var errNotInitialized = errors.New("cache is not initialized")

// Cache defines app objects
type Cache struct {
	mu     *sync.RWMutex
	hash   Hasher
	shards [shardCount]shard
}

// New creates app
func New(size uint64) *Cache {
	return &Cache{
		mu:   &sync.RWMutex{},
		hash: new(CRC32),
	}
}

// Set provides inserting to the cache
func (c *Cache) Set(key, value []byte) error {
	if c == nil {
		return errNotInitialized
	}

	return nil
}
