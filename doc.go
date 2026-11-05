// Package zzcache provides a sharded in-memory key/value cache
// with per-entry TTL and pluggable stores.
//
// Entries are distributed across shards by a CRC32 hash of the key;
// each shard has its own lock, so operations on different keys do not
// contend. The store behind each shard is pluggable: the default is a
// Go map, "radix" switches to a radix tree (useful for prefix-shaped
// keys).
//
// Basic usage:
//
//	c, err := zzcache.New(512, "")
//	if err != nil {
//		return err
//	}
//	if err := c.Set([]byte("key"), []byte("value"), time.Minute); err != nil {
//		return err
//	}
//	value, err := c.Get([]byte("key"))
//	if errors.Is(err, zzcache.ErrExpired) {
//		// the entry existed but its TTL elapsed and it was evicted
//	}
//
// A zero duration passed to Set means the entry never expires.
// Expired entries are evicted on read; a background janitor can evict
// them without reads:
//
//	stop := c.StartCleanup(time.Minute)
//	defer stop()
package zzcache
