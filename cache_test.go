package zzcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	c := New(10, "")
	assert.NoError(t, c.Set([]byte("key"), []byte("value"), 2*time.Second))
}
