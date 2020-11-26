package linkedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedMap(t *testing.T) {
	lm := NewLinkedMap(4)

	assert.True(t, lm.Empty())
	lm.PushBack(2, "b")
	lm.PushBack(1, "a")

	value, found := lm.Get(2)
	assert.Equal(t, found, true)
	assert.Equal(t, value, "b")

	value, found = lm.Get(5)
	assert.Equal(t, found, false)
	assert.Equal(t, value, nil)

	lm.Remove(2)

	value, found = lm.Get(2)
	assert.Equal(t, found, false)
	assert.Equal(t, value, nil)

	lm.PushBack(3, "c")
	lm.PushBack(3, "c")
	assert.Equal(t, lm.Size(), 2)

	lm.PushBack(4, "d")
	lm.PushBack(5, "e")

	value, found = lm.Get(1)
	assert.Equal(t, found, true)
	assert.Equal(t, value, "a")
	lm.PushBack(6, "f")

	value, found = lm.Get(1)
	assert.Equal(t, found, false)
	assert.Equal(t, value, nil)
}
