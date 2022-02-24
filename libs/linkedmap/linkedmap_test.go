package linkedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedMap(t *testing.T) {
	t.Run("Test FirstElement", func(t *testing.T) {
		lm := NewLinkedMap(4)
		k, v := lm.First()
		assert.Nil(t, lm.FirstElement())
		assert.Nil(t, k)
		assert.Nil(t, v)

		lm.PushFront(3, "c")
		lm.PushFront(2, "b")
		lm.PushFront(1, "a")

		k, v = lm.First()
		assert.Equal(t, lm.FirstElement().Value, &Pair{1, "a"})
		assert.Equal(t, k, 1)
		assert.Equal(t, v, "a")

		// Updating the pair
		lm.PushFront(1, "d")
		k, v = lm.First()
		assert.Equal(t, k, 1)
		assert.Equal(t, v, "d")
		assert.Equal(t, lm.Size(), 3)
	})

	t.Run("Test LastElement", func(t *testing.T) {
		lm := NewLinkedMap(4)
		k, v := lm.Last()
		assert.Nil(t, lm.LastElement())
		assert.Nil(t, k)
		assert.Nil(t, v)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")

		k, v = lm.Last()
		assert.Equal(t, lm.LastElement().Value, &Pair{3, "c"})
		assert.Equal(t, k, 3)
		assert.Equal(t, v, "c")

		// Updating the pair
		lm.PushBack(3, "d")
		k, v = lm.Last()
		assert.Equal(t, k, 3)
		assert.Equal(t, v, "d")
		assert.Equal(t, lm.Size(), 3)
	})

	t.Run("Test Get", func(t *testing.T) {
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

	t.Run("test Remove", func(t *testing.T) {
		lm := NewLinkedMap(4)

		lm.PushBack(0, "-")
		lm.PushBack(2, "b")
		lm.PushBack(1, "a")
		assert.True(t, lm.Remove(2))
		assert.False(t, lm.Remove(2))
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

	t.Run("Test PushBack and prune", func(t *testing.T) {
		lm := NewLinkedMap(3)

		lm.PushBack(1, "a") // This item should be pruned
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")
		lm.PushBack(4, "d")

		k, v := lm.First()
		assert.Equal(t, lm.FirstElement().Value, &Pair{2, "b"})
		assert.Equal(t, k, 2)
		assert.Equal(t, v, "b")
	})

	t.Run("Test PushFront and prune", func(t *testing.T) {
		lm := NewLinkedMap(3)

		lm.PushFront(1, "a")
		lm.PushFront(2, "b")
		lm.PushFront(3, "c")
		lm.PushFront(4, "d") // This item should be pruned

		k, v := lm.Last()
		assert.Equal(t, lm.LastElement().Value, &Pair{1, "a"})
		assert.Equal(t, k, 1)
		assert.Equal(t, v, "a")
	})

	t.Run("Deletd first ", func(t *testing.T) {
		lm := NewLinkedMap(3)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")

		lm.Remove(1)

		assert.Equal(t, lm.FirstElement().Value, &Pair{2, "b"})
	})

	t.Run("Delete last", func(t *testing.T) {
		lm := NewLinkedMap(3)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")

		lm.Remove(3)
		assert.Equal(t, lm.LastElement().Value, &Pair{2, "b"})
	})

	t.Run("Test Has function", func(t *testing.T) {
		lm1 := NewLinkedMap(2)

		lm1.PushBack(1, "a")

		assert.True(t, lm1.Has(1))
		assert.False(t, lm1.Has(2))
	})

	t.Run("Test Clear", func(t *testing.T) {
		lm1 := NewLinkedMap(2)

		lm1.PushBack(1, "a")
		lm1.Clear()
		assert.True(t, lm1.Empty())
	})
}

func TestSortingLinkedMap(t *testing.T) {
	lm1 := NewLinkedMap(6)

	lm1.PushBack(3, "c")
	lm1.PushBack(5, "e")
	lm1.PushBack(1, "a")
	lm1.PushBack(2, "b")
	lm1.PushBack(4, "d")

	lessThan := func(left interface{}, right interface{}) bool {
		return left.(string) < right.(string)
	}
	lm1.SortList(lessThan)
	assert.Equal(t, lm1.FirstElement().Value, &Pair{1, "a"})
	assert.Equal(t, lm1.LastElement().Value, &Pair{5, "e"})
	assert.Equal(t, lm1.Size(), 5)
}
