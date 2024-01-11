package pairslice

import (
	"golang.org/x/exp/slices"
)

// Pair represents a key-value pair.
type Pair[K comparable, V any] struct {
	First  K
	Second V
}

// PairSlice represents a slice of key-value pairs.
type PairSlice[K comparable, V any] struct {
	pairs []*Pair[K, V]
}

// New creates a new instance of PairSlice with a specified capacity.
func New[K comparable, V any](capacity int) *PairSlice[K, V] {
	return &PairSlice[K, V]{
		pairs: make([]*Pair[K, V], 0, capacity),
	}
}

// Append adds the first and second to the end of the slice.
func (ps *PairSlice[K, V]) Append(first K, second V) {
	ps.pairs = append(ps.pairs, &Pair[K, V]{first, second})
}

// RemoveFirst removes the first element from PairSlice.
func (ps *PairSlice[K, V]) RemoveFirst() {
	ps.remove(0)
}

// RemoveLast removes the last element from PairSlice.
func (ps *PairSlice[K, V]) RemoveLast() {
	ps.remove(ps.Len() - 1)
}

// Len returns the number of elements in the PairSlice.
func (ps *PairSlice[K, V]) Len() int {
	return len(ps.pairs)
}

// remove removes the element at the specified index from PairSlice.
func (ps *PairSlice[K, V]) remove(index int) {
	ps.pairs = slices.Delete(ps.pairs, index, index+1)
}

// Get returns the properties at the specified index. If the index is out of bounds, it returns false.
func (ps *PairSlice[K, V]) Get(index int) (K, V, bool) {
	if index < 0 || index >= len(ps.pairs) {
		var first K
		var second V

		return first, second, false
	}
	pair := ps.pairs[index]

	return pair.First, pair.Second, true
}

// First returns the first properties in the PairSlice. If the PairSlice is empty, it returns false.
func (ps *PairSlice[K, V]) First() (K, V, bool) {
	return ps.Get(0)
}

// Last returns the last properties in the PairSlice. If the PairSlice is empty, it returns false.
func (ps *PairSlice[K, V]) Last() (K, V, bool) {
	return ps.Get(ps.Len() - 1)
}
