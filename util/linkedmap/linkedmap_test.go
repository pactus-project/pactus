package linkedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedMap(t *testing.T) {
	t.Run("Test FirstNode", func(t *testing.T) {
		lm := NewLinkedMap[int, string](4)
		assert.Nil(t, lm.HeadNode())

		lm.PushFront(3, "c")
		lm.PushFront(2, "b")
		lm.PushFront(1, "a")

		assert.Equal(t, lm.HeadNode().Data.Key, 1)
		assert.Equal(t, lm.HeadNode().Data.Value, "a")
	})

	t.Run("Test LastNode", func(t *testing.T) {
		lm := NewLinkedMap[int, string](4)
		assert.Nil(t, lm.TailNode())

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")

		assert.Equal(t, lm.TailNode().Data.Key, 3)
		assert.Equal(t, lm.TailNode().Data.Value, "c")
	})

	t.Run("Test Get", func(t *testing.T) {
		lm := NewLinkedMap[int, string](4)

		lm.PushBack(2, "b")
		lm.PushBack(1, "a")

		n := lm.GetNode(2)
		assert.Equal(t, n.Data.Key, 2)
		assert.Equal(t, n.Data.Value, "b")

		n = lm.GetNode(5)
		assert.Nil(t, n)
	})

	t.Run("Test Remove", func(t *testing.T) {
		lm := NewLinkedMap[int, string](4)

		lm.PushBack(0, "-")
		lm.PushBack(2, "b")
		lm.PushBack(1, "a")
		assert.True(t, lm.Remove(2))
		assert.False(t, lm.Remove(2))
	})

	t.Run("Should updates v", func(t *testing.T) {
		lm := NewLinkedMap[int, string](4)
		lm.PushBack(1, "a")

		lm.PushBack(1, "b")
		n := lm.GetNode(1)
		assert.Equal(t, n.Data.Key, 1)
		assert.Equal(t, n.Data.Value, "b")

		lm.PushFront(1, "c")
		n = lm.GetNode(1)
		assert.Equal(t, n.Data.Key, 1)
		assert.Equal(t, n.Data.Value, "c")
	})

	t.Run("Should prunes oldest item", func(t *testing.T) {
		lm := NewLinkedMap[int, string](4)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")
		lm.PushBack(4, "d")

		n := lm.GetNode(1)
		assert.Equal(t, n.Data.Key, 1)
		assert.Equal(t, n.Data.Value, "a")

		lm.PushBack(5, "e")

		n = lm.GetNode(1)
		assert.Nil(t, n)
	})

	t.Run("Should prunes by changing capacity", func(t *testing.T) {
		lm := NewLinkedMap[int, string](4)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")
		lm.PushBack(4, "d")

		lm.SetCapacity(6)

		n := lm.GetNode(2)
		assert.Equal(t, n.Data.Key, 2)
		assert.Equal(t, n.Data.Value, "b")

		lm.SetCapacity(2)
		assert.True(t, lm.Full())

		n = lm.GetNode(2)
		assert.Nil(t, n)
	})

	t.Run("Test PushBack and prune", func(t *testing.T) {
		lm := NewLinkedMap[int, string](3)

		lm.PushBack(1, "a") // This item should be pruned
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")
		lm.PushBack(4, "d")

		n := lm.HeadNode()
		assert.Equal(t, n.Data.Key, 2)
		assert.Equal(t, n.Data.Value, "b")
	})

	t.Run("Test PushFront and prune", func(t *testing.T) {
		lm := NewLinkedMap[int, string](3)

		lm.PushFront(1, "a")
		lm.PushFront(2, "b")
		lm.PushFront(3, "c")
		lm.PushFront(4, "d") // This item should be pruned

		n := lm.TailNode()
		assert.Equal(t, n.Data.Key, 1)
		assert.Equal(t, n.Data.Value, "a")
	})

	t.Run("Delete first ", func(t *testing.T) {
		lm := NewLinkedMap[int, string](3)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")

		lm.Remove(1)

		assert.Equal(t, lm.HeadNode().Data.Key, 2)
		assert.Equal(t, lm.HeadNode().Data.Value, "b")
	})

	t.Run("Delete last", func(t *testing.T) {
		lm := NewLinkedMap[int, string](3)

		lm.PushBack(1, "a")
		lm.PushBack(2, "b")
		lm.PushBack(3, "c")

		lm.Remove(3)

		assert.Equal(t, lm.TailNode().Data.Key, 2)
		assert.Equal(t, lm.TailNode().Data.Value, "b")
	})

	t.Run("Test Has function", func(t *testing.T) {
		lm := NewLinkedMap[int, string](2)

		lm.PushBack(1, "a")

		assert.True(t, lm.Has(1))
		assert.False(t, lm.Has(2))
	})

	t.Run("Test Clear", func(t *testing.T) {
		lm := NewLinkedMap[int, string](2)

		lm.PushBack(1, "a")
		lm.Clear()
		assert.True(t, lm.Empty())
	})
}

func TestCapacity(t *testing.T) {
	capacity := 100
	lm := NewLinkedMap[int, string](capacity)
	assert.Equal(t, lm.Capacity(), capacity)
}
