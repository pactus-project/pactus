package orderedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func slice(om *OrderedMap) (ret []interface{}) {
	om.Iter(func(k, v interface{}) bool {
		ret = append(ret, om.Get(k))
		return true
	})
	return
}

func TestSimple(t *testing.T) {
	om := NewMap(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
	om.Set(3, 1)
	om.Set(2, 2)
	om.Set(1, 3)
	assert.EqualValues(t, []interface{}{3, 2, 1}, slice(om))
	om.Set(3, 2)
	om.Unset(2)
	assert.EqualValues(t, []interface{}{3, 2}, slice(om))
	om.Set(-1, 4)
	assert.EqualValues(t, []interface{}{4, 3, 2}, slice(om))
}

func TestIterEmpty(t *testing.T) {
	om := NewMap(nil)
	om.Iter(func(key, value interface{}) (more bool) {
		assert.Fail(t, "Iterating empty map.")
		return false
	})
}

func TestGetMinMax(t *testing.T) {
	om := NewMap(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})

	_, ok := om.MinKey()
	assert.False(t, ok)

	om.Set(3, 'a')
	om.Set(5, 'b')
	om.Set(1, 'c')
	om.Set(4, 'd')

	min, ok := om.MinKey()
	assert.True(t, ok)
	assert.Equal(t, min, 1)

	max, ok := om.MinKey()
	assert.True(t, ok)
	assert.Equal(t, max, 1)
}
