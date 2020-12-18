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

		v, ok := lm.Get(2)
		assert.Equal(t, ok, true)
		assert.Equal(t, v, "b")

		v, ok = lm.Get(5)
		assert.Equal(t, ok, false)
		assert.Equal(t, v, nil)
	})

	t.Run("Should removes item", func(t *testing.T) {
		lm := NewLinkedMap(4)

		lm.PushBack(0, "-")
		lm.PushBack(2, "b")
		lm.PushBack(1, "a")
		lm.Remove(2)

		v, ok := lm.Get(2)
		assert.Equal(t, ok, false)
		assert.Equal(t, v, nil)

	})

	t.Run("Should updates v", func(t *testing.T) {
		lm := NewLinkedMap(4)

		lm.PushBack(0xa, "a")
		lm.PushBack(0xa, "A")

		v, ok := lm.Get(0xa)
		assert.Equal(t, ok, true)
		assert.Equal(t, v, "A")
	})

	t.Run("Should prunes oldest item", func(t *testing.T) {
		lm := NewLinkedMap(4)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")
		lm.PushBack(4, "d")

		v, ok := lm.Get(1)
		assert.Equal(t, ok, true)
		assert.Equal(t, v, "a")

		lm.PushBack(5, "e")

		v, ok = lm.Get(1)
		assert.Equal(t, ok, false)
		assert.Equal(t, v, nil)
	})

	t.Run("Should prunes by changing capacity", func(t *testing.T) {
		lm := NewLinkedMap(4)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")
		lm.PushBack(4, "d")

		lm.SetCapacity(6)

		v, ok := lm.Get(2)
		assert.Equal(t, ok, true)
		assert.Equal(t, v, "b")

		lm.SetCapacity(2)
		assert.True(t, lm.Full())

		v, ok = lm.Get(2)
		assert.Equal(t, ok, false)
		assert.Equal(t, v, nil)
	})

	t.Run("Should returns first", func(t *testing.T) {
		lm := NewLinkedMap(3)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")
		lm.PushBack(4, "d") // pruning happens here

		el := lm.FirstElement()
		assert.Equal(t, el.Value, &Pair{2, "b"})

		k, v := lm.First()
		assert.Equal(t, k, 2)
		assert.Equal(t, v, "b")
	})

	t.Run("Should returns last", func(t *testing.T) {
		lm := NewLinkedMap(3)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")
		lm.PushBack(4, "d") // pruning happens here

		el := lm.LastElement()
		assert.Equal(t, el.Value, &Pair{4, "d"})

		k, v := lm.Last()
		assert.Equal(t, k, 4)
		assert.Equal(t, v, "d")
	})

	t.Run("Deletd first ", func(t *testing.T) {
		lm := NewLinkedMap(3)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")

		lm.Remove(1)

		el := lm.FirstElement()
		assert.Equal(t, el.Value, &Pair{2, "b"})
	})

	t.Run("Delete last", func(t *testing.T) {
		lm := NewLinkedMap(3)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")

		lm.Remove(3)
		el := lm.LastElement()
		assert.Equal(t, el.Value, &Pair{2, "b"})
	})
}
