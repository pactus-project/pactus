package pairslice

import "log"

type Pair[K comparable, V any] struct {
	FirstElement  K
	SecondElement V
}

type PairSlice[K comparable, V any] struct {
	pairs    []*Pair[K, V]
	index    int
	capacity int
}

// New creates a new instance of PairSlice with a specified capacity.
func New[K comparable, V any](capacity int) *PairSlice[K, V] {
	return &PairSlice[K, V]{
		pairs:    make([]*Pair[K, V], capacity),
		capacity: capacity,
	}
}

// Append adds the Pair to the end of the slice. If the capacity is full,
// it automatically removes the first element and appends the new Pair to the end of the PairSlice.
func (ps *PairSlice[K, V]) Append(firstElement K, secondElement V) {
	ps.pairs[ps.index] = &Pair[K, V]{firstElement, secondElement}
	log.Printf("memory address: %p", ps.pairs)
	ps.index++
	ps.resetIndex()
}

// Pop removes the first element from PairSlice.
func (ps *PairSlice[K, V]) Pop() {
	ps.pairs = ps.pairs[1:]
	ps.index--
	ps.resetIndex()
}

// Has checks the index exists or not.
func (ps *PairSlice[K, V]) Has(index int) bool {
	if index >= ps.capacity {
		return false
	}

	return ps.pairs[index] != nil
}

// Get returns the Pair at the specified index. If the index doesn't exist, it returns nil.
func (ps *PairSlice[K, V]) Get(index int) *Pair[K, V] {
	if index >= ps.capacity {
		return nil
	}
	return ps.pairs[index]
}

// First returns the first Pair in the PairSlice.
func (ps *PairSlice[K, V]) First() *Pair[K, V] {
	return ps.Get(0)
}

// Last returns the last Pair in the PairSlice.
func (ps *PairSlice[K, V]) Last() *Pair[K, V] {
	if ps.index == 0 {
		return ps.Get(len(ps.pairs) - 1)
	}
	return ps.Get(ps.index - 1)
}

func (ps *PairSlice[K, V]) resetIndex() {
	if ps.index < 0 {
		ps.index = 0
	}

	if ps.index == ps.capacity {
		ps.index = 0
	}
}
