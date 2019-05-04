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
	errNoKey          = errors.New("key is not defined")
	errNoValue        = errors.New("value is not defined")
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
	return c.set(key, value)
}

// Get provides getting data from the cache
func (c *Cache) Get(key []byte) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c == nil {
		return nil, errNotInitialized
	}

	value, ok := c.tree.Get(string(key))
	if value == nil || !ok {
		return nil, errNotFound
	}
	return value.([]byte), nil
}

// Delete provides deletetign data from the cache
func (c *Cache) Delete(key []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.tree.Delete(string(key))
	if !ok {
		return errNotFound
	}

	return nil
}

func (c *Cache) set(key, value []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c == nil {
		return errNotInitialized
	}
	if err := validateSet(key, value); err != nil {
		return err
	}

	_, ok := c.tree.Insert(string(key), value)
	if !ok {
		return errNotInserted
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
