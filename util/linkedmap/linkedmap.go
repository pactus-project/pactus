package linkedmap

import (
	ll "github.com/pactus-project/pactus/util/linkedlist"
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

type LinkedMap[K comparable, V any] struct {
	list     *ll.LinkedList[Pair[K, V]]
	hashmap  map[K]*ll.Element[Pair[K, V]]
	capacity int
}

// New creates a new LinkedMap with the specified capacity.
func New[K comparable, V any](capacity int) *LinkedMap[K, V] {
	return &LinkedMap[K, V]{
		list:     ll.New[Pair[K, V]](),
		hashmap:  make(map[K]*ll.Element[Pair[K, V]]),
		capacity: capacity,
	}
}

// SetCapacity sets the capacity of the LinkedMap and prunes the excess elements if needed.
func (lm *LinkedMap[K, V]) SetCapacity(capacity int) {
	lm.capacity = capacity

	lm.prune()
}

// Has checks if the specified key exists in the LinkedMap.
func (lm *LinkedMap[K, V]) Has(key K) bool {
	_, found := lm.hashmap[key]

	return found
}

// PushBack adds a new key-value pair to the end of the LinkedMap.
func (lm *LinkedMap[K, V]) PushBack(key K, value V) {
	ln, found := lm.hashmap[key]
	if found {
		// Update the value if the key already exists
		ln.Data.Value = value

		return
	}

	p := Pair[K, V]{Key: key, Value: value}
	ln = lm.list.InsertAtTail(p)
	lm.hashmap[key] = ln

	lm.prune()
}

// PushFront adds a new key-value pair to the beginning of the LinkedMap.
func (lm *LinkedMap[K, V]) PushFront(key K, value V) {
	ln, found := lm.hashmap[key]
	if found {
		// Update the value if the key already exists
		ln.Data.Value = value

		return
	}

	p := Pair[K, V]{Key: key, Value: value}
	ln = lm.list.InsertAtHead(p)
	lm.hashmap[key] = ln

	lm.prune()
}

// GetNode returns the LinkNode corresponding to the specified key.
func (lm *LinkedMap[K, V]) GetNode(key K) *ll.Element[Pair[K, V]] {
	ln, found := lm.hashmap[key]
	if found {
		return ln
	}

	return nil
}

// TailNode returns the LinkNode at the end (tail) of the LinkedMap.
func (lm *LinkedMap[K, V]) TailNode() *ll.Element[Pair[K, V]] {
	ln := lm.list.Tail
	if ln == nil {
		return nil
	}

	return ln
}

func (lm *LinkedMap[K, V]) RemoveTail() {
	lm.doRemove(lm.list.Tail, lm.list.Tail.Data.Key)
}

// HeadNode returns the LinkNode at the beginning (head) of the LinkedMap.
func (lm *LinkedMap[K, V]) HeadNode() *ll.Element[Pair[K, V]] {
	ln := lm.list.Head
	if ln == nil {
		return nil
	}

	return ln
}

func (lm *LinkedMap[K, V]) RemoveHead() {
	lm.doRemove(lm.list.Head, lm.list.Head.Data.Key)
}

// Remove removes the key-value pair with the specified key from the LinkedMap.
// It returns true if the key was found and removed, otherwise false.
func (lm *LinkedMap[K, V]) Remove(key K) bool {
	ln, found := lm.hashmap[key]
	if found {
		lm.list.Delete(ln)
		delete(lm.hashmap, ln.Data.Key)
	}

	return found
}

// doRemove removes the key-value pair with the specified key from the LinkedMap and LinkedList.
func (lm *LinkedMap[K, V]) doRemove(element *ll.Element[Pair[K, V]], key K) {
	lm.list.Delete(element)
	delete(lm.hashmap, key)
}

// Empty checks if the LinkedMap is empty (contains no key-value pairs).
func (lm *LinkedMap[K, V]) Empty() bool {
	return lm.Size() == 0
}

// Capacity returns the capacity of the LinkedMap.
func (lm *LinkedMap[K, V]) Capacity() int {
	return lm.capacity
}

// Size returns the number of key-value pairs in the LinkedMap.
func (lm *LinkedMap[K, V]) Size() int {
	return lm.list.Length()
}

// Full checks if the LinkedMap is full (reached its capacity).
func (lm *LinkedMap[K, V]) Full() bool {
	return lm.list.Length() == lm.capacity
}

// Clear removes all key-value pairs from the LinkedMap, making it empty.
func (lm *LinkedMap[K, V]) Clear() {
	lm.list.Clear()
	lm.hashmap = make(map[K]*ll.Element[Pair[K, V]])
}

// prune removes excess elements from the LinkedMap if its size exceeds the capacity.
func (lm *LinkedMap[K, V]) prune() {
	if lm.capacity == 0 {
		return
	}

	for lm.list.Length() > lm.capacity {
		head := lm.list.Head
		key := head.Data.Key
		lm.list.Delete(head)
		delete(lm.hashmap, key)
	}
}
