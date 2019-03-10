package zzcache

import (
	"errors"
	"sync"
)

var errNotInitialized = errors.New("cache is not initialized")

// Cache defines app objects
type Cache struct {
	mu *sync.RWMutex
}

// New creates app
func New(size uint64) *Cache {
	return &Cache{
		mu: &sync.RWMutex{},
	}
}

// Set provides inserting to the cache
func (c *Cache) Set(key, value []byte) error {
	if c == nil {
		return errNotInitialized
	}

	return nil
}
