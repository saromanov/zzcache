package zzcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	c, err := New(10, "")
	assert.NoError(t, err)
	assert.NoError(t, c.Set([]byte("key"), []byte("value"), 2*time.Second))

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
	assert.NoError(t, c.Set([]byte("key"), []byte("value"), 2*time.Second))
	d, err := c.Get([]byte("key"))
	assert.NoError(t, err)
	assert.Equal(t, "value", string(d))
	time.Sleep(2 * time.Second)
	_, err = c.Get([]byte("key"))
	assert.Error(t, err)
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
