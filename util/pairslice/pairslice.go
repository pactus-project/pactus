package pairslice

import (
	"golang.org/x/exp/slices"
)

type Pair[K comparable, V any] struct {
	First  K
	Second V
}

type PairSlice[K comparable, V any] struct {
	pairs []*Pair[K, V]
}

// New creates a new instance of PairSlice with a specified capacity.
func New[K comparable, V any](capacity int) *PairSlice[K, V] {
	return &PairSlice[K, V]{
		pairs: make([]*Pair[K, V], 0, capacity),
	}
}

// Append adds the Pair to the end of the slice. If the capacity is full,
// it automatically removes the first element and appends the new Pair to the end of the PairSlice.
func (ps *PairSlice[K, V]) Append(first K, second V) {
	ps.pairs = append(ps.pairs, &Pair[K, V]{first, second})
}

// RemoveFirst removes the first element from PairSlice.
func (ps *PairSlice[K, V]) RemoveFirst() {
	ps.remove(0)
}

// RemoveLast removes the first element from PairSlice.
func (ps *PairSlice[K, V]) RemoveLast() {
	ps.remove(ps.Len() - 1)
}

func (ps *PairSlice[K, V]) Len() int {
	return len(ps.pairs)
}

func (ps *PairSlice[K, V]) remove(index int) {
	ps.pairs = slices.Delete(ps.pairs, index, index+1)
}

// Get returns the Pair at the specified index. If the index doesn't exist, it returns nil.
func (ps *PairSlice[K, V]) Get(index int) (K, V, bool) {
	if index >= len(ps.pairs) || index < 0 {
		var first K
		var second V
		return first, second, false
	}
	pair := ps.pairs[index]
	return pair.First, pair.Second, true
}

// First returns the first Pair in the PairSlice.
func (ps *PairSlice[K, V]) First() (K, V, bool) {
	return ps.Get(0)
}

// Last returns the last Pair in the PairSlice.
func (ps *PairSlice[K, V]) Last() (K, V, bool) {
	return ps.Get(ps.Len() - 1)
}
