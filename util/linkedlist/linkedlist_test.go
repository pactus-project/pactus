package linkedlist_test

import (
	"testing"

	"github.com/pactus-project/pactus/util/linkedlist"
	"github.com/stretchr/testify/assert"
)

func TestDoublyLink_InsertAtHead(t *testing.T) {
	list := linkedlist.New[int]()
	list.InsertAtHead(1)
	list.InsertAtHead(2)
	list.InsertAtHead(3)
	list.InsertAtHead(4)

	assert.Equal(t, []int{4, 3, 2, 1}, list.Values())
	assert.Equal(t, 4, list.Length())
	assert.Equal(t, 4, list.Head.Data)
	assert.Equal(t, 1, list.Tail.Data)
}

func TestSinglyLink_InsertAtTail(t *testing.T) {
	list := linkedlist.New[int]()
	list.InsertAtTail(1)
	list.InsertAtTail(2)
	list.InsertAtTail(3)
	list.InsertAtTail(4)

	assert.Equal(t, []int{1, 2, 3, 4}, list.Values())
	assert.Equal(t, 4, list.Length())
	assert.Equal(t, 1, list.Head.Data)
	assert.Equal(t, 4, list.Tail.Data)
}

func TestDeleteAtHead(t *testing.T) {
	list := linkedlist.New[int]()
	list.InsertAtTail(1)
	list.InsertAtTail(2)
	list.InsertAtTail(3)

	list.DeleteAtHead()
	assert.Equal(t, []int{2, 3}, list.Values())
	assert.Equal(t, 2, list.Length())

	list.DeleteAtHead()
	assert.Equal(t, []int{3}, list.Values())
	assert.Equal(t, 1, list.Length())

	list.DeleteAtHead()
	assert.Equal(t, []int{}, list.Values())
	assert.Equal(t, 0, list.Length())

	list.DeleteAtHead()
	assert.Equal(t, []int{}, list.Values())
	assert.Equal(t, 0, list.Length())
}

func TestDeleteAtTail(t *testing.T) {
	list := linkedlist.New[int]()
	list.InsertAtTail(1)
	list.InsertAtTail(2)
	list.InsertAtTail(3)

	list.DeleteAtTail()
	assert.Equal(t, []int{1, 2}, list.Values())
	assert.Equal(t, 2, list.Length())

	list.DeleteAtTail()
	assert.Equal(t, []int{1}, list.Values())
	assert.Equal(t, 1, list.Length())

	list.DeleteAtTail()
	assert.Equal(t, []int{}, list.Values())
	assert.Equal(t, 0, list.Length())

	list.DeleteAtTail()
	assert.Equal(t, []int{}, list.Values())
	assert.Equal(t, 0, list.Length())
}

func TestDelete(t *testing.T) {
	list := linkedlist.New[int]()
	elm1 := list.InsertAtTail(1)
	elm2 := list.InsertAtTail(2)
	elm3 := list.InsertAtTail(3)
	elm4 := list.InsertAtTail(4)

	list.Delete(elm1)
	assert.Equal(t, []int{2, 3, 4}, list.Values())
	assert.Equal(t, 3, list.Length())

	list.Delete(elm4)
	assert.Equal(t, []int{2, 3}, list.Values())
	assert.Equal(t, 2, list.Length())

	list.Delete(elm2)
	assert.Equal(t, []int{3}, list.Values())
	assert.Equal(t, 1, list.Length())

	list.Delete(elm3)
	assert.Equal(t, []int{}, list.Values())
	assert.Equal(t, 0, list.Length())
}

func TestClear(t *testing.T) {
	list := linkedlist.New[int]()
	list.InsertAtTail(1)
	list.InsertAtTail(2)
	list.InsertAtTail(3)

	list.Clear()
	assert.Equal(t, []int{}, list.Values())
	assert.Equal(t, 0, list.Length())
}

func TestInsertAfter(t *testing.T) {
	list := linkedlist.New[int]()
	elm1 := list.InsertAtHead(1)
	elm2 := list.InsertAfter(2, elm1)
	list.InsertAfter(3, elm2)
	list.InsertAfter(4, list.Head)
	list.InsertAfter(5, list.Tail)

	assert.Equal(t, []int{1, 4, 2, 3, 5}, list.Values())
}

func TestInsertBefore(t *testing.T) {
	list := linkedlist.New[int]()
	elm1 := list.InsertAtHead(1)
	elm2 := list.InsertBefore(2, elm1)
	list.InsertBefore(3, elm2)
	list.InsertBefore(4, list.Head)
	list.InsertBefore(5, list.Tail)

	assert.Equal(t, []int{4, 3, 2, 5, 1}, list.Values())
}

func TestDefensiveOperations(t *testing.T) {
	t.Run("Delete is idempotent and ignores nil", func(t *testing.T) {
		list := linkedlist.New[int]()
		head := list.InsertAtTail(1)
		tailElm := list.InsertAtTail(2)

		list.Delete(nil)
		list.Delete(head)
		list.Delete(head) // The second delete has no effect.

		assert.Equal(t, 1, list.Length())
		assert.Equal(t, []int{2}, list.Values())

		list.Delete(tailElm)
		assert.Zero(t, list.Length())
		assert.Nil(t, list.Head)
		assert.Nil(t, list.Tail)
	})

	t.Run("an element of another list is ignored", func(t *testing.T) {
		first := linkedlist.New[int]()
		second := linkedlist.New[int]()
		first.InsertAtTail(1)
		foreign := second.InsertAtTail(2)

		first.Delete(foreign)
		assert.Equal(t, 1, first.Length())
		assert.Equal(t, 1, second.Length())
	})

	t.Run("inserting at a nil or removed position", func(t *testing.T) {
		list := linkedlist.New[int]()
		head := list.InsertAtTail(1)
		list.Delete(head)

		list.InsertBefore(2, nil)
		list.InsertAfter(3, head) // a is not in the list anymore.
		list.InsertBefore(4, head)

		assert.Equal(t, 3, list.Length())
		assert.Len(t, list.Values(), 3)
	})

	t.Run("elements removed by Clear can't corrupt the list", func(t *testing.T) {
		list := linkedlist.New[int]()
		head := list.InsertAtTail(1)
		list.Clear()
		list.InsertAtTail(2)

		list.Delete(head)
		assert.Equal(t, 1, list.Length())
		assert.Equal(t, []int{2}, list.Values())
	})
}
