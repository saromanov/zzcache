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

// newShard creates a new shard
func newShard() *shard {
	return &shard{
		l: list.New(),
	}
}

// set provides inserting of the data
func set(key, value []byte) error {
	if len(key) > maxKeySize {
		return errKeyTooLarge
	}

	return nil
}
