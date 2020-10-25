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
}

func TestGet(t *testing.T) {
	c, err := New(10, "")
	assert.NoError(t, err)
	assert.NoError(t, c.Set([]byte("key"), []byte("value"), 2*time.Second))
	d, err := c.Get([]byte("key"))
	assert.NoError(t, err)
	assert.Equal(t, "value", string(d))
}
