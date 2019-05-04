package zzcache

import (
	"errors"
	"sync"

	"github.com/armon/go-radix"
)

const shardCount = 512

var (
	errNotInitialized = errors.New("cache is not initialized")
	errNotInserted    = errors.New("unable to insert data")
	errNotFound       = errors.New("element is not found")
)

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

	_, ok := c.tree.Insert(string(key), value)
	if !ok {
		return errNotInserted
	}
	return nil
}

// Get provides getting data from the cache
func (c *Cache) Get(key []byte) error {
	if c == nil {
		return errNotInitialized
	}

	value, ok := c.tree.Get(string(key))
	if value == nil || !ok {
		return errNotFound
	}
	return nil
}
