package zzcache

import (
	"errors"
	"sync"
	"time"

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
	store  Store
}

type Item struct {
	key   string
	value interface{}
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
func (c *Cache) Set(key, value []byte, d time.Duration) error {
	hash := c.hash.Do(key)
	shardID := hash & shardCount
	return c.set(shardID, key, value)
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

	_, ok := c.tree.Delete(string(key))
	if !ok {
		return errNotFound
	}

	defer c.mu.Unlock()
	return nil
}

// inner method for validating of the input data
// and append data to shards
func (c *Cache) set(shardID uint32, key, value []byte) error {
	if c == nil {
		return errNotInitialized
	}
	if err := validateSet(key, value); err != nil {
		return err
	}

	if err := c.shards[shardID].set(key, value); err != nil {
		return err
	}
	c.mu.Lock()
	_, ok := c.tree.Insert(string(key), value)
	if !ok {
		return errNotInserted
	}
	defer c.mu.Unlock()
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

	value, ok := c.tree.Get(string(key))
	if value == nil || !ok {
		return nil, errNotFound
	}
	return value.([]byte), nil
}
