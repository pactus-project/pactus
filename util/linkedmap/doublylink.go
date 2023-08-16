package linkedmap

type LinkNode[T any] struct {
	Data T
	Next *LinkNode[T]
	Prev *LinkNode[T]
}

func NewLinkNode[T any](data T) *LinkNode[T] {
	return &LinkNode[T]{
		Data: data,
		Next: nil,
		Prev: nil,
	}
}

// DoublyLinkedList represents a doubly linked list.
type DoublyLinkedList[T any] struct {
	Head   *LinkNode[T]
	Tail   *LinkNode[T]
	length int
}

func NewDoublyLinkedList[T any]() *DoublyLinkedList[T] {
	return &DoublyLinkedList[T]{
		Head:   nil,
		Tail:   nil,
		length: 0,
	}
}

// InsertAtHead inserts a new node at the head of the list.
func (l *DoublyLinkedList[T]) InsertAtHead(data T) *LinkNode[T] {
	newNode := NewLinkNode(data)

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
func (l *DoublyLinkedList[T]) InsertAtTail(data T) *LinkNode[T] {
	newNode := NewLinkNode(data)

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

// DeleteAtHead deletes the node at the head of the list.
func (l *DoublyLinkedList[T]) DeleteAtHead() {
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
func (l *DoublyLinkedList[T]) DeleteAtTail() {
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
func (l *DoublyLinkedList[T]) Delete(ln *LinkNode[T]) {
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
func (l *DoublyLinkedList[T]) Length() int {
	return l.length
}

// Values returns a slice of values in the list.
func (l *DoublyLinkedList[T]) Values() []T {
	values := []T{}
	cur := l.Head
	for cur != nil {
		values = append(values, cur.Data)
		cur = cur.Next
	}
	return values
}

// Clear removes all nodes from the list, making it empty.
func (l *DoublyLinkedList[T]) Clear() {
	l.Head = nil
	l.Tail = nil
	l.length = 0
}
