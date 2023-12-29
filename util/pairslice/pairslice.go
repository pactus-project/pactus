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

func New[K comparable, V any](capacity int) *PairSlice[K, V] {
	return &PairSlice[K, V]{
		pairs:    make([]*Pair[K, V], capacity),
		capacity: capacity,
	}
}

func (ps *PairSlice[K, V]) Append(firstElement K, secondElement V) {
	ps.pairs[ps.index] = &Pair[K, V]{firstElement, secondElement}
	log.Printf("memory address: %p", ps.pairs)
	ps.index++
	ps.resetIndex()
}

func (ps *PairSlice[K, V]) Pop() {
	ps.pairs = ps.pairs[1:]
	ps.index--
	ps.resetIndex()
}

func (ps *PairSlice[K, V]) Has(index int) bool {
	if index >= ps.capacity {
		return false
	}

	return ps.pairs[index] != nil
}

func (ps *PairSlice[K, V]) Get(index int) *Pair[K, V] {
	if index >= ps.capacity {
		return nil
	}
	return ps.pairs[index]
}

func (ps *PairSlice[K, V]) First() *Pair[K, V] {
	return ps.Get(0)
}

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
