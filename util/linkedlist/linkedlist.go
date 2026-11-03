package linkedlist

type Element[T any] struct {
	Data T
	Next *Element[T]
	Prev *Element[T]

	// list is the list that the element belongs to. It is nil for a removed element.
	list *LinkedList[T]
}

func NewElement[T any](data T) *Element[T] {
	return &Element[T]{
		Data: data,
		Next: nil,
		Prev: nil,
	}
}

// LinkedList represents a doubly linked list.
type LinkedList[T any] struct {
	Head   *Element[T]
	Tail   *Element[T]
	length int
}

func New[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		Head:   nil,
		Tail:   nil,
		length: 0,
	}
}

// InsertAtHead inserts a new element at the head of the list.
func (ll *LinkedList[T]) InsertAtHead(data T) *Element[T] {
	elm := NewElement(data)
	elm.list = ll
	if ll.Head == nil {
		// Empty list case
		ll.Head = elm
		ll.Tail = elm
	} else {
		elm.Next = ll.Head
		ll.Head.Prev = elm
		ll.Head = elm
	}

	ll.length++

	return elm
}

// InsertAtTail appends a new element at the tail of the list.
func (ll *LinkedList[T]) InsertAtTail(data T) *Element[T] {
	elm := NewElement(data)
	elm.list = ll
	if ll.Head == nil {
		// Empty list case
		ll.Head = elm
		ll.Tail = elm
	} else {
		elm.Prev = ll.Tail
		ll.Tail.Next = elm
		ll.Tail = elm
	}

	ll.length++

	return elm
}

// InsertBefore inserts a new element before pos. If pos is not an element of the list,
// the new element is inserted at the head.
func (ll *LinkedList[T]) InsertBefore(data T, pos *Element[T]) *Element[T] {
	if pos == nil || pos.list != ll {
		return ll.InsertAtHead(data)
	}

	elm := NewElement[T](data)
	elm.list = ll
	if pos == ll.Head {
		ll.Head = elm
		elm.Next = pos
		elm.Next.Prev = elm
	} else {
		elm.Prev = pos.Prev
		elm.Next = pos
		elm.Next.Prev = elm
		elm.Prev.Next = elm
	}
	ll.length++

	return elm
}

// InsertAfter inserts a new element after pos. If pos is not an element of the list,
// the new element is inserted at the tail.
func (ll *LinkedList[T]) InsertAfter(data T, pos *Element[T]) *Element[T] {
	if pos == nil || pos.list != ll {
		return ll.InsertAtTail(data)
	}

	elm := NewElement[T](data)
	elm.list = ll
	if pos == ll.Tail {
		ll.Tail = elm
		elm.Prev = pos
		elm.Prev.Next = elm
	} else {
		elm.Prev = pos
		elm.Next = pos.Next
		elm.Prev.Next = elm
		elm.Next.Prev = elm
	}
	ll.length++

	return elm
}

// DeleteAtHead deletes the element at the head of the list.
func (ll *LinkedList[T]) DeleteAtHead() {
	if ll.Head == nil {
		// Empty list case
		return
	}

	removed := ll.Head
	ll.Head = ll.Head.Next
	if ll.Head != nil {
		ll.Head.Prev = nil
	} else {
		ll.Tail = nil
	}
	removed.list = nil

	ll.length--
}

// DeleteAtTail deletes the element at the tail of the list.
func (ll *LinkedList[T]) DeleteAtTail() {
	if ll.Tail == nil {
		// Empty list case
		return
	}

	removed := ll.Tail
	ll.Tail = ll.Tail.Prev
	if ll.Tail != nil {
		ll.Tail.Next = nil
	} else {
		ll.Head = nil
	}
	removed.list = nil

	ll.length--
}

// Delete removes a specific element from the list.
// It does nothing if the element is nil, was already removed, or belongs to another list.
func (ll *LinkedList[T]) Delete(elm *Element[T]) {
	if elm == nil || elm.list != ll {
		return
	}
	elm.list = nil

	if elm.Prev != nil {
		elm.Prev.Next = elm.Next
	} else {
		ll.Head = elm.Next
	}

	if elm.Next != nil {
		elm.Next.Prev = elm.Prev
	} else {
		ll.Tail = elm.Prev
	}

	ll.length--
}

// Length returns the number of elements in the list.
func (ll *LinkedList[T]) Length() int {
	return ll.length
}

// Values returns a slice of values in the list.
func (ll *LinkedList[T]) Values() []T {
	values := []T{}
	cur := ll.Head
	for cur != nil {
		values = append(values, cur.Data)
		cur = cur.Next
	}

	return values
}

// Clear removes all elements from the list, making it empty.
func (ll *LinkedList[T]) Clear() {
	for elm := ll.Head; elm != nil; elm = elm.Next {
		elm.list = nil
	}

	ll.Head = nil
	ll.Tail = nil
	ll.length = 0
}
