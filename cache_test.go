package zzcache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	c := New(10, "")
	assert.NoError(t, c.Set("key", "value"))
}
