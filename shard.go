package zzcache

import (
	"container/list"
	"errors"
)

const maxKeySize = 65535

var errKeyTooLarge = errors.New("key is too large")

// shard provides implementation of the shard for cache
type shard struct {
	l *list.List
}

// entry represents d-s for inserting to linked list
type entry struct {
	key   []byte
	value []byte
}

// newShard creates a new shard
func newShard() *shard {
	return &shard{
		l: list.New(),
	}
}

// set provides inserting of the data
func (s *shard) set(key, value []byte) error {
	if len(key) > maxKeySize {
		return errKeyTooLarge
	}

	s.l.PushFront(&entry{key, value})
	return nil
}
