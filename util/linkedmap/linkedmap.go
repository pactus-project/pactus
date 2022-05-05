package linkedmap

import (
	"container/list"
)

// TODO: should be thread safe or not?

type Pair struct {
	First, Second interface{}
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

func (lm *LinkedMap) PushBack(first interface{}, second interface{}) {
	el, found := lm.hashmap[first]
	if found {
		// update the second
		el.Value.(*Pair).Second = second
		return
	}

	el = lm.list.PushBack(&Pair{first, second})
	lm.hashmap[first] = el

	lm.prune()
}

func (lm *LinkedMap) PushFront(first interface{}, second interface{}) {
	el, found := lm.hashmap[first]
	if found {
		// update the second
		el.Value.(*Pair).Second = second
		return
	}

	el = lm.list.PushFront(&Pair{first, second})
	lm.hashmap[first] = el

	lm.prune()
}

func (lm *LinkedMap) Get(first interface{}) (interface{}, bool) {
	el, found := lm.hashmap[first]
	if found {
		return el.Value.(*Pair).Second, true

	}
	return nil, false
}

func (lm *LinkedMap) Last() (interface{}, interface{}) {
	el := lm.list.Back()
	if el == nil {
		return nil, nil
	}
	p := el.Value.(*Pair)
	return p.First, p.Second
}

func (lm *LinkedMap) First() (interface{}, interface{}) {
	el := lm.list.Front()
	if el == nil {
		return nil, nil
	}
	p := el.Value.(*Pair)
	return p.First, p.Second
}

func (lm *LinkedMap) LastElement() *list.Element {
	return lm.list.Back()
}

func (lm *LinkedMap) FirstElement() *list.Element {
	return lm.list.Front()
}

func (lm *LinkedMap) Remove(first interface{}) bool {
	el, found := lm.hashmap[first]
	if found {
		lm.list.Remove(el)
		delete(lm.hashmap, el.Value.(*Pair).First)
	}
	return found
}

func (lm *LinkedMap) Empty() bool {
	return lm.Size() == 0
}

func (lm *LinkedMap) Capacity() int {
	return lm.capacity
}

func (lm *LinkedMap) Size() int {
	return lm.list.Len()
}

func (lm *LinkedMap) Full() bool {
	return lm.list.Len() == lm.capacity
}

func (lm *LinkedMap) Clear() {
	lm.list = list.New()
	lm.hashmap = make(map[interface{}]*list.Element)
}

func (lm *LinkedMap) prune() {
	for lm.list.Len() > lm.capacity {
		front := lm.list.Front()
		key := front.Value.(*Pair).First
		lm.list.Remove(front)
		delete(lm.hashmap, key)
	}
}

func (lm *LinkedMap) SortList(cmp func(left interface{}, right interface{}) bool) {
	index := lm.list.Front()
	if index == nil {
		return
	}

	for index != nil {
		current := index.Next()
		for current != nil {
			if cmp(current.Value.(*Pair).Second, index.Value.(*Pair).Second) {
				lm.list.MoveBefore(current, index)
				index = current
				current = index
			}
			current = current.Next()
		}

		index = index.Next()
	}
}
