package linkedmap

import (
	"github.com/pactus-project/pactus/util/linkedlist"
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

type LinkedMap[K comparable, V any] struct {
	list     *linkedlist.LinkedList[Pair[K, V]]
	hashmap  map[K]*linkedlist.Element[Pair[K, V]]
	capacity int
}

// New creates a new LinkedMap with the specified capacity.
func New[K comparable, V any](capacity int) *LinkedMap[K, V] {
	return &LinkedMap[K, V]{
		list:     linkedlist.New[Pair[K, V]](),
		hashmap:  make(map[K]*linkedlist.Element[Pair[K, V]]),
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
	node, found := lm.hashmap[key]
	if found {
		// Update the value if the key already exists
		node.Data.Value = value

		return
	}

	p := Pair[K, V]{Key: key, Value: value}
	node = lm.list.InsertAtTail(p)
	lm.hashmap[key] = node

	lm.prune()
}

// PushFront adds a new key-value pair to the beginning of the LinkedMap.
func (lm *LinkedMap[K, V]) PushFront(key K, value V) {
	node, found := lm.hashmap[key]
	if found {
		// Update the value if the key already exists
		node.Data.Value = value

		return
	}

	p := Pair[K, V]{Key: key, Value: value}
	node = lm.list.InsertAtHead(p)
	lm.hashmap[key] = node

	lm.prune()
}

// GetNode returns the Element corresponding to the specified key.
func (lm *LinkedMap[K, V]) GetNode(key K) *linkedlist.Element[Pair[K, V]] {
	node, found := lm.hashmap[key]
	if found {
		return node
	}

	return nil
}

// TailNode returns the Element at the end (tail) of the LinkedMap.
func (lm *LinkedMap[K, V]) TailNode() *linkedlist.Element[Pair[K, V]] {
	node := lm.list.Tail
	if node == nil {
		return nil
	}

	return node
}

func (lm *LinkedMap[K, V]) RemoveTail() {
	lm.remove(lm.list.Tail)
}

// HeadNode returns the Element at the beginning (head) of the LinkedMap.
func (lm *LinkedMap[K, V]) HeadNode() *linkedlist.Element[Pair[K, V]] {
	node := lm.list.Head
	if node == nil {
		return nil
	}

	return node
}

func (lm *LinkedMap[K, V]) RemoveHead() {
	lm.remove(lm.list.Head)
}

// Remove removes the key-value pair with the specified key from the LinkedMap.
// It returns true if the key was found and removed, otherwise false.
func (lm *LinkedMap[K, V]) Remove(key K) bool {
	element, found := lm.hashmap[key]
	if found {
		lm.remove(element)
	}

	return found
}

// remove removes the specified element pair from the LinkedMap.
func (lm *LinkedMap[K, V]) remove(element *linkedlist.Element[Pair[K, V]]) {
	lm.list.Delete(element)
	delete(lm.hashmap, element.Data.Key)
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
	lm.hashmap = make(map[K]*linkedlist.Element[Pair[K, V]])
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
