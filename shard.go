package zzcache

import "errors"

const maxKeySize = 65535

var errKeyTooLarge = errors.New("key is too large")

// shard provides implementation of the shard for cache
type shard struct {

}

// set provides inserting of the data
func set(key, value []byte) error {
	if len(key) > maxKeySize {
		return errKeyTooLarge
	}

	return nil
}
