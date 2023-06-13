package linkedmap

// TODO: should be thread safe or not?

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

type LinkedMap[K comparable, V any] struct {
	list     *DoublyLinkedList[Pair[K, V]]
	hashmap  map[K]*LinkNode[Pair[K, V]]
	capacity int
}

func NewLinkedMap[K comparable, V any](capacity int) *LinkedMap[K, V] {
	return &LinkedMap[K, V]{
		list:     NewDoublyLinkedList[Pair[K, V]](),
		hashmap:  make(map[K]*LinkNode[Pair[K, V]]),
		capacity: capacity,
	}
}

func (lm *LinkedMap[K, V]) SetCapacity(capacity int) {
	lm.capacity = capacity

	lm.prune()
}

func (lm *LinkedMap[K, V]) Has(key K) bool {
	_, found := lm.hashmap[key]
	return found
}

func (lm *LinkedMap[K, V]) PushBack(key K, value V) {
	ln, found := lm.hashmap[key]
	if found {
		// update the second
		ln.Data.Value = value
		return
	}

	p := Pair[K, V]{Key: key, Value: value}
	ln = lm.list.InsertAtTail(p)
	lm.hashmap[key] = ln

	lm.prune()
}

func (lm *LinkedMap[K, V]) PushFront(key K, value V) {
	ln, found := lm.hashmap[key]
	if found {
		// update the second
		ln.Data.Value = value
		return
	}

	p := Pair[K, V]{Key: key, Value: value}
	ln = lm.list.InsertAtHead(p)
	lm.hashmap[key] = ln

	lm.prune()
}

func (lm *LinkedMap[K, V]) GetNode(key K) *LinkNode[Pair[K, V]] {
	ln, found := lm.hashmap[key]
	if found {
		return ln
	}
	return nil
}

func (lm *LinkedMap[K, V]) LastNode() *LinkNode[Pair[K, V]] {
	ln := lm.list.Tail
	if ln == nil {
		return nil
	}
	return ln
}

func (lm *LinkedMap[K, V]) FirstNode() *LinkNode[Pair[K, V]] {
	ln := lm.list.Head
	if ln == nil {
		return nil
	}
	return ln
}

func (lm *LinkedMap[K, V]) Remove(key K) bool {
	nl, found := lm.hashmap[key]
	if found {
		lm.list.Delete(nl)
		delete(lm.hashmap, nl.Data.Key)
	}
	return found
}

func (lm *LinkedMap[K, V]) Empty() bool {
	return lm.Size() == 0
}

func (lm *LinkedMap[K, V]) Capacity() int {
	return lm.capacity
}

func (lm *LinkedMap[K, V]) Size() int {
	return lm.list.Length()
}

func (lm *LinkedMap[K, V]) Full() bool {
	return lm.list.Length() == lm.capacity
}

func (lm *LinkedMap[K, V]) Clear() {
	lm.list.Clear()
	lm.hashmap = make(map[K]*LinkNode[Pair[K, V]])
}

func (lm *LinkedMap[K, V]) prune() {
	for lm.list.Length() > lm.capacity {
		front := lm.list.Head
		key := front.Data.Key
		lm.list.Delete(front)
		delete(lm.hashmap, key)
	}
}
