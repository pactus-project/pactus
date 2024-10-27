package linkedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedMap(t *testing.T) {
	t.Run("Test FirstNode", func(t *testing.T) {
		lmp := New[int, string](4)
		assert.Nil(t, lmp.HeadNode())

		lmp.PushFront(3, "c")
		lmp.PushFront(2, "b")
		lmp.PushFront(1, "a")

		assert.Equal(t, lmp.HeadNode().Data.Key, 1)
		assert.Equal(t, lmp.HeadNode().Data.Value, "a")
	})

	t.Run("Test LastNode", func(t *testing.T) {
		lmp := New[int, string](4)
		assert.Nil(t, lmp.TailNode())

		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")
		lmp.PushBack(3, "c")

		assert.Equal(t, lmp.TailNode().Data.Key, 3)
		assert.Equal(t, lmp.TailNode().Data.Value, "c")
	})

	t.Run("Test Get", func(t *testing.T) {
		lmp := New[int, string](4)

		lmp.PushBack(2, "b")
		lmp.PushBack(1, "a")

		node := lmp.GetNode(2)
		assert.Equal(t, node.Data.Key, 2)
		assert.Equal(t, node.Data.Value, "b")

		node = lmp.GetNode(5)
		assert.Nil(t, node)
	})

	t.Run("Test Remove", func(t *testing.T) {
		lmp := New[int, string](4)

		lmp.PushBack(0, "-")
		lmp.PushBack(2, "b")
		lmp.PushBack(1, "a")
		assert.True(t, lmp.Remove(2))
		assert.False(t, lmp.Remove(2))
	})

	t.Run("Test RemoveTail", func(t *testing.T) {
		lmp := New[int, string](4)
		lmp.PushBack(0, "-")
		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")

		lmp.RemoveTail()
		assert.Equal(t, lmp.TailNode().Data.Value, "a")
		assert.NotEqual(t, lmp.TailNode().Data.Value, "b")
	})

	t.Run("Test RemoveHead", func(t *testing.T) {
		lmp := New[int, string](4)
		lmp.PushBack(0, "-")
		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")

		lmp.RemoveHead()
		assert.Equal(t, lmp.HeadNode().Data.Value, "a")
		assert.NotEqual(t, lmp.HeadNode().Data.Value, "-")
	})

	t.Run("Should updates v", func(t *testing.T) {
		lmp := New[int, string](4)
		lmp.PushBack(1, "a")

		lmp.PushBack(1, "b")
		node := lmp.GetNode(1)
		assert.Equal(t, node.Data.Key, 1)
		assert.Equal(t, node.Data.Value, "b")

		lmp.PushFront(1, "c")
		node = lmp.GetNode(1)
		assert.Equal(t, node.Data.Key, 1)
		assert.Equal(t, node.Data.Value, "c")
	})

	t.Run("Should prunes oldest item", func(t *testing.T) {
		lmp := New[int, string](4)

		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")
		lmp.PushBack(3, "c")
		lmp.PushBack(4, "d")

		node := lmp.GetNode(1)
		assert.Equal(t, node.Data.Key, 1)
		assert.Equal(t, node.Data.Value, "a")

		lmp.PushBack(5, "e")

		node = lmp.GetNode(1)
		assert.Nil(t, node)
	})

	t.Run("Should prunes by changing capacity", func(t *testing.T) {
		lmp := New[int, string](4)

		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")
		lmp.PushBack(3, "c")
		lmp.PushBack(4, "d")

		lmp.SetCapacity(6)

		node := lmp.GetNode(2)
		assert.Equal(t, node.Data.Key, 2)
		assert.Equal(t, node.Data.Value, "b")

		lmp.SetCapacity(2)
		assert.True(t, lmp.Full())

		node = lmp.GetNode(2)
		assert.Nil(t, node)
	})

	t.Run("Test PushBack and prune", func(t *testing.T) {
		lmp := New[int, string](3)

		lmp.PushBack(1, "a") // This item should be pruned
		lmp.PushBack(2, "b")
		lmp.PushBack(3, "c")
		lmp.PushBack(4, "d")

		node := lmp.HeadNode()
		assert.Equal(t, node.Data.Key, 2)
		assert.Equal(t, node.Data.Value, "b")
	})

	t.Run("Test PushFront and prune", func(t *testing.T) {
		lmp := New[int, string](3)

		lmp.PushFront(1, "a")
		lmp.PushFront(2, "b")
		lmp.PushFront(3, "c")
		lmp.PushFront(4, "d") // This item should be pruned

		node := lmp.TailNode()
		assert.Equal(t, node.Data.Key, 1)
		assert.Equal(t, node.Data.Value, "a")
	})

	t.Run("Delete first ", func(t *testing.T) {
		lmp := New[int, string](3)

		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")
		lmp.PushBack(3, "c")

		lmp.Remove(1)

		assert.Equal(t, lmp.HeadNode().Data.Key, 2)
		assert.Equal(t, lmp.HeadNode().Data.Value, "b")
	})

	t.Run("Delete last", func(t *testing.T) {
		lmp := New[int, string](3)

		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")
		lmp.PushBack(3, "c")

		lmp.Remove(3)

		assert.Equal(t, lmp.TailNode().Data.Key, 2)
		assert.Equal(t, lmp.TailNode().Data.Value, "b")
	})

	t.Run("Test Has function", func(t *testing.T) {
		lmp := New[int, string](2)

		lmp.PushBack(1, "a")

		assert.True(t, lmp.Has(1))
		assert.False(t, lmp.Has(2))
	})

	t.Run("Test Clear", func(t *testing.T) {
		lmp := New[int, string](2)

		lmp.PushBack(1, "a")
		lmp.Clear()
		assert.True(t, lmp.Empty())
	})
}

func TestCapacity(t *testing.T) {
	t.Run("Check Capacity", func(t *testing.T) {
		capacity := 100
		lmp := New[int, string](capacity)
		assert.Equal(t, lmp.Capacity(), capacity)
	})

	t.Run("Test FirstNode with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)
		assert.Nil(t, lmp.HeadNode())

		lmp.PushFront(3, "c")
		lmp.PushFront(2, "b")
		lmp.PushFront(1, "a")

		assert.Equal(t, lmp.HeadNode().Data.Key, 1)
		assert.Equal(t, lmp.HeadNode().Data.Value, "a")
	})

	t.Run("Test LastNode with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)
		assert.Nil(t, lmp.TailNode())

		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")
		lmp.PushBack(3, "c")

		assert.Equal(t, lmp.TailNode().Data.Key, 3)
		assert.Equal(t, lmp.TailNode().Data.Value, "c")
	})

	t.Run("Test Get with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)

		lmp.PushBack(2, "b")
		lmp.PushBack(1, "a")

		node := lmp.GetNode(2)
		assert.Equal(t, node.Data.Key, 2)
		assert.Equal(t, node.Data.Value, "b")

		node = lmp.GetNode(5)
		assert.Nil(t, node)
	})

	t.Run("Test Remove with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)

		lmp.PushBack(0, "-")
		lmp.PushBack(2, "b")
		lmp.PushBack(1, "a")
		assert.True(t, lmp.Remove(2))
		assert.False(t, lmp.Remove(2))
	})

	t.Run("Test RemoveTail with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)
		lmp.PushBack(0, "-")
		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")

		lmp.RemoveTail()
		assert.Equal(t, lmp.TailNode().Data.Value, "a")
		assert.NotEqual(t, lmp.TailNode().Data.Value, "b")
	})

	t.Run("Test RemoveHead with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)
		lmp.PushBack(0, "-")
		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")

		lmp.RemoveHead()
		assert.Equal(t, lmp.HeadNode().Data.Value, "a")
		assert.NotEqual(t, lmp.HeadNode().Data.Value, "-")
	})

	t.Run("Should updates v with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)
		lmp.PushBack(1, "a")

		lmp.PushBack(1, "b")
		node := lmp.GetNode(1)
		assert.Equal(t, node.Data.Key, 1)
		assert.Equal(t, node.Data.Value, "b")

		lmp.PushFront(1, "c")
		node = lmp.GetNode(1)
		assert.Equal(t, node.Data.Key, 1)
		assert.Equal(t, node.Data.Value, "c")
	})

	t.Run("Should not prunes oldest item with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)

		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")
		lmp.PushBack(3, "c")
		lmp.PushBack(4, "d")

		node := lmp.GetNode(1)
		assert.Equal(t, node.Data.Key, 1)
		assert.Equal(t, node.Data.Value, "a")

		lmp.PushBack(5, "e")

		node = lmp.GetNode(1)
		assert.NotNil(t, node)
	})

	t.Run("Should prunes by changing capacity with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)

		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")
		lmp.PushBack(3, "c")
		lmp.PushBack(4, "d")

		lmp.SetCapacity(6)

		node := lmp.GetNode(2)
		assert.Equal(t, node.Data.Key, 2)
		assert.Equal(t, node.Data.Value, "b")

		lmp.SetCapacity(2)
		assert.True(t, lmp.Full())

		node = lmp.GetNode(2)
		assert.Nil(t, node)
	})

	t.Run("Test PushBack and should not prune with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)

		lmp.PushBack(1, "a") // This item should be pruned
		lmp.PushBack(2, "b")
		lmp.PushBack(3, "c")
		lmp.PushBack(4, "d")

		node := lmp.TailNode()
		assert.Equal(t, node.Data.Key, 4)
		assert.Equal(t, node.Data.Value, "d")
	})

	t.Run("Test PushFront and prune with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)

		lmp.PushFront(1, "a")
		lmp.PushFront(2, "b")
		lmp.PushFront(3, "c")
		lmp.PushFront(4, "d") // This item should be pruned

		n := lmp.TailNode()
		assert.Equal(t, n.Data.Key, 1)
		assert.Equal(t, n.Data.Value, "a")
	})

	t.Run("Delete first with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)

		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")
		lmp.PushBack(3, "c")

		lmp.Remove(1)

		assert.Equal(t, lmp.HeadNode().Data.Key, 2)
		assert.Equal(t, lmp.HeadNode().Data.Value, "b")
	})

	t.Run("Delete last with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)

		lmp.PushBack(1, "a")
		lmp.PushBack(2, "b")
		lmp.PushBack(3, "c")

		lmp.Remove(3)

		assert.Equal(t, lmp.TailNode().Data.Key, 2)
		assert.Equal(t, lmp.TailNode().Data.Value, "b")
	})

	t.Run("Test Has function with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)

		lmp.PushBack(1, "a")

		assert.True(t, lmp.Has(1))
		assert.False(t, lmp.Has(2))
	})

	t.Run("Test Clear with Zero Capacity", func(t *testing.T) {
		lmp := New[int, string](0)

		lmp.PushBack(1, "a")
		lmp.Clear()
		assert.True(t, lmp.Empty())
	})
}
