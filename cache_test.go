package zzcache

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// fakeHasher returns the preconfigured hash value.
type fakeHasher struct {
	v uint32
}

func (f fakeHasher) Do([]byte) uint32 {
	return f.v
}

func TestSet(t *testing.T) {
	c, err := New(10, "")
	assert.NoError(t, err)
	assert.NoError(t, c.Set([]byte("key"), []byte("value"), 2*time.Second))
	assert.Error(t, c.Set([]byte(""), []byte("value"), 2*time.Second))
	assert.Error(t, c.Set([]byte("key"), []byte(""), 2*time.Second))

	_, err = New(0, "")
	assert.Error(t, err)

	c, err = New(10, "radix")
	assert.NoError(t, err)
	assert.NoError(t, c.Set([]byte("key"), []byte("value"), 2*time.Second))

	_, err = New(0, "")
	assert.Error(t, err)
}

func TestGet(t *testing.T) {
	c, err := New(10, "")
	assert.NoError(t, err)
	assert.NoError(t, c.Set([]byte("key"), []byte("value"), 20*time.Millisecond))
	d, err := c.Get([]byte("key"))
	assert.NoError(t, err)
	assert.Equal(t, "value", string(d))
	time.Sleep(30 * time.Millisecond)
	_, err = c.Get([]byte("key"))
	assert.Error(t, err)

	c, err = New(10, "radix")
	assert.NoError(t, err)
	assert.NoError(t, c.Set([]byte("key"), []byte("value"), 20*time.Millisecond))
	d, err = c.Get([]byte("key"))
	assert.NoError(t, err)
	assert.Equal(t, "value", string(d))
}

func TestDelete(t *testing.T) {
	c, err := New(5, "")
	assert.NoError(t, err)
	assert.NoError(t, c.Set([]byte("key"), []byte("value"), 2*time.Second))
	assert.NoError(t, c.Delete([]byte("key")))
	_, err = c.Get([]byte("key"))
	assert.Error(t, err)

	assert.NoError(t, c.Delete([]byte("va")))
}

// Deleting a missing key must not error out and must not
// leave the cache locked (it deadlocked before the fix).
func TestDeleteMissingKeyNoDeadlock(t *testing.T) {
	c, err := New(5, "")
	assert.NoError(t, err)
	assert.NoError(t, c.Delete([]byte("missing")))
	// If the mutex is stuck, this call hangs forever.
	assert.NoError(t, c.Set([]byte("key"), []byte("value"), 2*time.Second))
	d, err := c.Get([]byte("key"))
	assert.NoError(t, err)
	assert.Equal(t, "value", string(d))

	c2, err := New(5, "radix")
	assert.NoError(t, err)
	assert.NoError(t, c2.Delete([]byte("missing")))
}

// A zero duration must mean "no expiration".
func TestNoTTL(t *testing.T) {
	c, err := New(5, "")
	assert.NoError(t, err)
	assert.NoError(t, c.Set([]byte("key"), []byte("value"), 0))
	time.Sleep(30 * time.Millisecond)
	d, err := c.Get([]byte("key"))
	assert.NoError(t, err)
	assert.Equal(t, "value", string(d))

	c2, err := New(5, "radix")
	assert.NoError(t, err)
	assert.NoError(t, c2.Set([]byte("key"), []byte("value"), 0))
	time.Sleep(30 * time.Millisecond)
	d, err = c2.Get([]byte("key"))
	assert.NoError(t, err)
	assert.Equal(t, "value", string(d))
}

// Expired entries must be evicted, not just reported.
func TestExpiredEntryEvicted(t *testing.T) {
	c, err := New(5, "")
	assert.NoError(t, err)
	assert.NoError(t, c.Set([]byte("key"), []byte("value"), 10*time.Millisecond))
	time.Sleep(30 * time.Millisecond)
	_, err = c.Get([]byte("key"))
	assert.True(t, errors.Is(err, ErrExpired))

	// The entry must be gone from the underlying store.
	sh := c.shards[c.getShard([]byte("key"))]
	_, err = sh.store.Get("key")
	assert.Error(t, err)
}

// Expired vs missing must be distinguishable.
func TestSentinelErrors(t *testing.T) {
	c, err := New(5, "")
	assert.NoError(t, err)
	_, err = c.Get([]byte("nope"))
	assert.True(t, errors.Is(err, ErrNotFound))
	_, err = c.Get(nil)
	// A nil key hashes fine, but the entry does not exist.
	assert.True(t, errors.Is(err, ErrNotFound))

	err = c.Set([]byte(""), []byte("v"), time.Second)
	assert.True(t, errors.Is(err, ErrNoKey))
	err = c.Set([]byte("k"), []byte(""), time.Second)
	assert.True(t, errors.Is(err, ErrNoValue))

	big := make([]byte, maxKeySize+1)
	err = c.Set(big, []byte("v"), time.Second)
	assert.True(t, errors.Is(err, ErrKeyTooLarge))
}

// getShard must never return an index outside the shard slice.
// With the old hash&size masking and a non-power-of-two size,
// index == size was reachable and panicked.
func TestShardIndexBounds(t *testing.T) {
	const size = 10
	for h := uint32(0); h < 32; h++ {
		c, err := New(size, "")
		assert.NoError(t, err)
		c.hash = fakeHasher{v: h}
		key := []byte("key")
		assert.NoError(t, c.Set(key, []byte("value"), time.Second))
		d, err := c.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, "value", string(d))
	}
}

// All shards must stay in use with the modulo mapping.
func TestShardsAllReachable(t *testing.T) {
	const size = 4
	c, err := New(size, "")
	assert.NoError(t, err)
	used := map[uint32]bool{}
	for h := uint32(0); h < 16; h++ {
		c.hash = fakeHasher{v: h}
		used[c.getShard([]byte("k"))] = true
	}
	assert.Len(t, used, size)
}

func TestNotInitialized(t *testing.T) {
	var c *Cache
	err := c.Set([]byte("k"), []byte("v"), time.Second)
	assert.True(t, errors.Is(err, ErrNotInitialized))
	_, err = c.Get([]byte("k"))
	assert.True(t, errors.Is(err, ErrNotInitialized))
	err = c.Delete([]byte("k"))
	assert.True(t, errors.Is(err, ErrNotInitialized))
}
