package linkedlist

type Element[T any] struct {
	Data T
	Next *Element[T]
	Prev *Element[T]
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

// InsertAtHead inserts a new node at the head of the list.
func (l *LinkedList[T]) InsertAtHead(data T) *Element[T] {
	newNode := NewElement(data)

	if l.Head == nil {
		// Empty list case
		l.Head = newNode
		l.Tail = newNode
	} else {
		newNode.Next = l.Head
		l.Head.Prev = newNode
		l.Head = newNode
	}

	l.length++

	return newNode
}

// InsertAtTail appends a new node at the tail of the list.
func (l *LinkedList[T]) InsertAtTail(data T) *Element[T] {
	newNode := NewElement(data)

	if l.Head == nil {
		// Empty list case
		l.Head = newNode
		l.Tail = newNode
	} else {
		newNode.Prev = l.Tail
		l.Tail.Next = newNode
		l.Tail = newNode
	}

	l.length++

	return newNode
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *LinkedList[T]) insertValue(data T, at *Element[T]) *Element[T] {
	return l.insert(NewElement[T](data), at)
}

func (l *LinkedList[T]) insert(e, at *Element[T]) *Element[T] {
	e.Prev = at
	e.Next = at.Next
	e.Prev.Next = e
	e.Next.Prev = e
	l.length++
	return e
}

func (l *LinkedList[T]) InsertBefore(data T, at *Element[T]) *Element[T] {
	if at == l.Head {
		e := NewElement[T](data)
		l.Head = e
		l.Head.Next = at
		at.Prev = e
		l.length++
		return e
	} else {
		return l.insertValue(data, at.Prev)
	}
}

func (l *LinkedList[T]) InsertAfter(data T, at *Element[T]) *Element[T] {
	if at == l.Tail {
		e := NewElement[T](data)
		l.Tail = e
		l.Tail.Prev = at
		at.Next = e
		l.length++
		return e
	} else {
		return l.insertValue(data, at)
	}
}

// DeleteAtHead deletes the node at the head of the list.
func (l *LinkedList[T]) DeleteAtHead() {
	if l.Head == nil {
		// Empty list case
		return
	}

	l.Head = l.Head.Next
	if l.Head != nil {
		l.Head.Prev = nil
	} else {
		l.Tail = nil
	}

	l.length--
}

// DeleteAtTail deletes the node at the tail of the list.
func (l *LinkedList[T]) DeleteAtTail() {
	if l.Tail == nil {
		// Empty list case
		return
	}

	l.Tail = l.Tail.Prev
	if l.Tail != nil {
		l.Tail.Next = nil
	} else {
		l.Head = nil
	}

	l.length--
}

// Delete removes a specific node from the list.
func (l *LinkedList[T]) Delete(ln *Element[T]) {
	if ln.Prev != nil {
		ln.Prev.Next = ln.Next
	} else {
		l.Head = ln.Next
	}

	if ln.Next != nil {
		ln.Next.Prev = ln.Prev
	} else {
		l.Tail = ln.Prev
	}

	l.length--
}

// Length returns the number of nodes in the list.
func (l *LinkedList[T]) Length() int {
	return l.length
}

// Values returns a slice of values in the list.
func (l *LinkedList[T]) Values() []T {
	values := []T{}
	cur := l.Head
	for cur != nil {
		values = append(values, cur.Data)
		cur = cur.Next
	}
	return values
}

// Clear removes all nodes from the list, making it empty.
func (l *LinkedList[T]) Clear() {
	l.Head = nil
	l.Tail = nil
	l.length = 0
}
