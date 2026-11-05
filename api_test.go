package zzcache

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setBatch(t *testing.T, c *Cache, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		key := []byte("cnt-" + string(rune('a'+i%26)) + "-" + itoa(i))
		if err := c.Set(key, []byte("v"), time.Hour); err != nil {
			t.Fatal(err)
		}
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	digits := []byte{}
	for i > 0 {
		digits = append([]byte{byte('0' + i%10)}, digits...)
		i /= 10
	}
	return string(digits)
}

func TestCount(t *testing.T) {
	for _, storeType := range []string{"", "radix"} {
		c, err := New(8, storeType)
		assert.NoError(t, err)
		assert.Equal(t, 0, c.Count())
		setBatch(t, c, 50)
		assert.Equal(t, 50, c.Count())
	}
}

func TestFlush(t *testing.T) {
	for _, storeType := range []string{"", "radix"} {
		c, err := New(8, storeType)
		assert.NoError(t, err)
		setBatch(t, c, 20)
		assert.Equal(t, 20, c.Count())
		c.Flush()
		assert.Equal(t, 0, c.Count())
		_, err = c.Get([]byte("cnt-a-0"))
		assert.True(t, errors.Is(err, ErrNotFound))
		// The cache must stay usable after a flush.
		assert.NoError(t, c.Set([]byte("k"), []byte("v"), time.Second))
	}
}

func TestKeys(t *testing.T) {
	c, err := New(8, "")
	assert.NoError(t, err)
	setBatch(t, c, 10)
	keys := c.Keys()
	assert.Len(t, keys, 10)
	seen := map[string]bool{}
	for _, k := range keys {
		seen[string(k)] = true
	}
	assert.Len(t, seen, 10)
}

func TestStats(t *testing.T) {
	c, err := New(4, "")
	assert.NoError(t, err)
	assert.NoError(t, c.Set([]byte("hit"), []byte("v"), time.Second))
	_, _ = c.Get([]byte("hit"))  // hit
	_, _ = c.Get([]byte("hit"))  // hit
	_, _ = c.Get([]byte("miss")) // miss
	st := c.Stats()
	assert.Equal(t, uint64(2), st.Hits)
	assert.Equal(t, uint64(1), st.Misses)
}

// The janitor must evict expired entries on its own.
func TestStartCleanup(t *testing.T) {
	c, err := New(4, "")
	assert.NoError(t, err)
	assert.NoError(t, c.Set([]byte("gone"), []byte("v"), 10*time.Millisecond))
	stop := c.StartCleanup(5 * time.Millisecond)
	defer stop()
	time.Sleep(60 * time.Millisecond)
	assert.Equal(t, 0, c.Count())
	_, err = c.Get([]byte("gone"))
	assert.True(t, errors.Is(err, ErrNotFound))
}

// A cache without a running janitor must keep entries until read.
func TestNoCleanupWithoutStart(t *testing.T) {
	c, err := New(4, "")
	assert.NoError(t, err)
	assert.NoError(t, c.Set([]byte("stay"), []byte("v"), 10*time.Millisecond))
	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, 1, c.Count())
}

func TestAPIOnNilCache(t *testing.T) {
	var c *Cache
	assert.Equal(t, 0, c.Count())
	assert.Nil(t, c.Keys())
	assert.NotPanics(t, func() { c.Flush() })
	assert.Equal(t, Stats{}, c.Stats())
	assert.NotPanics(t, func() {
		stop := c.StartCleanup(time.Millisecond)
		stop()
	})
}
