package zzcache

import (
	"errors"
	"sync"

	"github.com/armon/go-radix"
)

const shardCount = 512

var errNotInitialized = errors.New("cache is not initialized")

// Cache defines app objects
type Cache struct {
	mu     *sync.RWMutex
	hash   Hasher
	shards [shardCount]shard
	tree   *radix.Tree
}

// New creates app
func New(size uint64) *Cache {
	return &Cache{
		mu:   &sync.RWMutex{},
		hash: new(CRC32),
		tree: radix.New(),
	}
}

// Set provides inserting to the cache
func (c *Cache) Set(key, value []byte) error {
	if c == nil {
		return errNotInitialized
	}

	return nil
}
