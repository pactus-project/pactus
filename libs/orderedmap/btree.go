package orderedmap

import (
	"github.com/google/btree"
)

type OrderedMap struct {
	bt     *btree.BTree
	lesser func(l, r interface{}) bool
}

type item struct {
	less  func(l, r interface{}) bool
	key   interface{}
	value interface{}
}

func (me item) Less(right btree.Item) bool {
	return me.less(me.key, right.(*item).key)
}

func NewMap(lesser func(l, r interface{}) bool) *OrderedMap {
	return &OrderedMap{
		bt:     btree.New(32),
		lesser: lesser,
	}
}

func (me *OrderedMap) Set(key interface{}, value interface{}) {
	me.bt.ReplaceOrInsert(&item{me.lesser, key, value})
}

func (me *OrderedMap) Get(key interface{}) interface{} {
	ret, _ := me.GetOk(key)
	return ret
}

func (me *OrderedMap) GetOk(key interface{}) (interface{}, bool) {
	i := me.bt.Get(&item{me.lesser, key, nil})
	if i == nil {
		return nil, false
	}
	return i.(*item).value, true
}

// Callback receives a value and returns true if another value should be
// received or false to stop iteration.
type callback func(key, value interface{}) (more bool)

func (me *OrderedMap) Iter(f callback) {
	me.bt.Ascend(func(i btree.Item) bool {
		return f(i.(*item).key, i.(*item).value)
	})
}

func (me *OrderedMap) Unset(key interface{}) {
	me.bt.Delete(&item{me.lesser, key, nil})
}

func (me *OrderedMap) Len() int {
	return me.bt.Len()
}
