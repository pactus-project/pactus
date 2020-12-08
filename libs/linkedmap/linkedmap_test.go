package linkedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedMap(t *testing.T) {

	t.Run("Should adds items", func(t *testing.T) {
		lm := NewLinkedMap(4)

		lm.PushBack(2, "b")
		lm.PushBack(1, "a")

		value, found := lm.Get(2)
		assert.Equal(t, found, true)
		assert.Equal(t, value, "b")

		value, found = lm.Get(5)
		assert.Equal(t, found, false)
		assert.Equal(t, value, nil)
	})

	t.Run("Should removes item", func(t *testing.T) {
		lm := NewLinkedMap(4)

		lm.PushBack(2, "b")
		lm.PushBack(1, "a")
		lm.Remove(2)

		value, found := lm.Get(2)
		assert.Equal(t, found, false)
		assert.Equal(t, value, nil)
	})

	t.Run("Should updates value", func(t *testing.T) {
		lm := NewLinkedMap(4)

		lm.PushBack(0xa, "a")
		lm.PushBack(0xa, "A")

		value, found := lm.Get(0xa)
		assert.Equal(t, found, true)
		assert.Equal(t, value, "A")
	})

	t.Run("Should prunces oldest item", func(t *testing.T) {
		lm := NewLinkedMap(4)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")
		lm.PushBack(4, "d")

		value, found := lm.Get(1)
		assert.Equal(t, found, true)
		assert.Equal(t, value, "a")

		lm.PushBack(5, "e")

		value, found = lm.Get(1)
		assert.Equal(t, found, false)
		assert.Equal(t, value, nil)
	})

	t.Run("Should prunces by changing capacity", func(t *testing.T) {
		lm := NewLinkedMap(4)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")
		lm.PushBack(4, "d")

		lm.SetCapacity(6)

		value, found := lm.Get(2)
		assert.Equal(t, found, true)
		assert.Equal(t, value, "b")

		lm.SetCapacity(2)
		assert.True(t, lm.Full())

		value, found = lm.Get(2)
		assert.Equal(t, found, false)
		assert.Equal(t, value, nil)
	})

	t.Run("Should returns first", func(t *testing.T) {
		lm := NewLinkedMap(3)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")
		lm.PushBack(4, "d")

		k, v := lm.First()
		assert.Equal(t, k, 2)
		assert.Equal(t, v, "b")
	})

	t.Run("Should returns last", func(t *testing.T) {
		lm := NewLinkedMap(3)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")
		lm.PushBack(4, "d")

		k, v := lm.Last()
		assert.Equal(t, k, 4)
		assert.Equal(t, v, "d")
	})

}
