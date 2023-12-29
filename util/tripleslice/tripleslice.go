package tripleslice

type Triple[K comparable, V any] struct {
	FirstElement  K
	SecondElement V
	ThirdElement  int
}

type TripleSlice[K comparable, V any] struct {
	pairs    []*Triple[K, V]
	index    int
	capacity int
}

// New creates a new instance of TripleSlice with a specified capacity.
func New[K comparable, V any](capacity int) *TripleSlice[K, V] {
	return &TripleSlice[K, V]{
		pairs:    make([]*Triple[K, V], capacity),
		capacity: capacity,
	}
}

// Append adds the Triple to the end of the slice. If the capacity is full,
// it automatically removes the first element and appends the new Triple to the end of the TripleSlice.
func (ps *TripleSlice[K, V]) Append(firstElement K, secondElement V) {
	ps.pairs[ps.index] = &Triple[K, V]{firstElement, secondElement, ps.index}
	ps.index++
	ps.resetIndex()
}

// Pop removes the first element from TripleSlice.
func (ps *TripleSlice[K, V]) Pop() {
	ps.pairs = ps.pairs[1:]
	ps.index--
	ps.resetIndex()
}

// Has checks the index exists or not.
func (ps *TripleSlice[K, V]) Has(index int) bool {
	if index >= ps.capacity {
		return false
	}

	return ps.pairs[index] != nil
}

// Get returns the Triple at the specified index. If the index doesn't exist, it returns nil.
func (ps *TripleSlice[K, V]) Get(index int) *Triple[K, V] {
	if index >= ps.capacity {
		return nil
	}
	return ps.pairs[index]
}

// First returns the first Triple in the TripleSlice.
func (ps *TripleSlice[K, V]) First() *Triple[K, V] {
	return ps.Get(0)
}

// Last returns the last Triple in the TripleSlice.
func (ps *TripleSlice[K, V]) Last() *Triple[K, V] {
	if ps.index == 0 {
		return ps.Get(len(ps.pairs) - 1)
	}
	return ps.Get(ps.index - 1)
}

// All returns whole pairs of slice.
func (ps *TripleSlice[K, V]) All() []*Triple[K, V] {
	return ps.pairs
}

func (ps *TripleSlice[K, V]) resetIndex() {
	if ps.index < 0 {
		ps.index = 0
	}

	if ps.index == ps.capacity {
		ps.index = 0
	}
}
