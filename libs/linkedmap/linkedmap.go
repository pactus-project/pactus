package linkedmap

import (
	"container/list"
)

type pair struct {
	key, value interface{}
}

type LinkedMap struct {
	list     *list.List
	hashmap  map[interface{}]*list.Element
	capacity int
}

func NewLinkedMap(capacity int) *LinkedMap {
	return &LinkedMap{
		list:     list.New(),
		hashmap:  make(map[interface{}]*list.Element),
		capacity: capacity,
	}
}

func (lm *LinkedMap) SetCapacity(capacity int) {
	lm.capacity = capacity

	lm.prune()
}

func (lm *LinkedMap) Has(key interface{}) bool {
	_, found := lm.hashmap[key]
	return found
}

func (lm *LinkedMap) PushBack(key interface{}, value interface{}) {
	_, found := lm.hashmap[key]
	if found {
		return
	}

	lm.prune()

	el := lm.list.PushBack(pair{key, value})
	lm.hashmap[key] = el
}

func (lm *LinkedMap) Get(key interface{}) (interface{}, bool) {
	el, found := lm.hashmap[key]
	if found {
		return el.Value.(pair).value, true

	}
	return nil, false
}

func (lm *LinkedMap) Remove(key interface{}) {
	el, found := lm.hashmap[key]
	if found {
		lm.list.Remove(el)
		delete(lm.hashmap, el.Value.(pair).key)
	}
}

func (lm *LinkedMap) Empty() bool {
	return lm.Size() == 0
}

func (lm *LinkedMap) Size() int {
	return lm.list.Len()
}

func (lm *LinkedMap) Clear() {
	lm.list = list.New()
	lm.hashmap = make(map[interface{}]*list.Element)
}

func (lm *LinkedMap) prune() {
	for lm.list.Len() >= lm.capacity {
		front := lm.list.Front()
		key := front.Value.(pair).key
		lm.list.Remove(front)
		delete(lm.hashmap, key)
	}
}
