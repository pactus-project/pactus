package pairslice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ps := New[int, string](10)

	assert.NotNil(t, ps)
	assert.Equal(t, 10, cap(ps.pairs))
	assert.Equal(t, 0, len(ps.pairs))
}

func TestPairSlice(t *testing.T) {

	t.Run("Test Append", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(1, "a")
		ps.Append(2, "b")
		ps.Append(3, "c")

		assert.Equal(t, 3, ps.Len())
		assert.Equal(t, ps.pairs[2].First, 3)
		assert.Equal(t, ps.pairs[2].Second, "c")
	})

	t.Run("Test RemoveFirst", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(1, "a")
		ps.Append(2, "b")
		ps.Append(3, "c")

		ps.RemoveFirst()

		assert.Equal(t, ps.pairs[0].First, 2)
		assert.Equal(t, ps.pairs[0].Second, "b")
	})

	t.Run("Test RemoveLast", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(1, "a")
		ps.Append(2, "b")
		ps.Append(3, "c")
		ps.Append(4, "d")

		ps.RemoveLast()

		assert.Equal(t, ps.pairs[2].First, 3)
		assert.Equal(t, ps.pairs[2].Second, "c")
	})

	t.Run("Test Len", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(1, "a")
		ps.Append(2, "b")

		assert.Equal(t, 2, ps.Len())
	})

	t.Run("Test Remove", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(1, "a")
		ps.Append(2, "b")
		ps.Append(3, "c")
		ps.Append(4, "d")

		ps.remove(1)

		assert.Equal(t, ps.pairs[1].First, 3)
		assert.Equal(t, ps.pairs[1].Second, "c")
	})

	t.Run("Test Get", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(1, "a")
		ps.Append(2, "b")
		ps.Append(3, "c")
		ps.Append(4, "d")

		first, second, _ := ps.Get(2)
		assert.Equal(t, ps.pairs[2].First, first)
		assert.Equal(t, ps.pairs[2].Second, second)
	})

	t.Run("Test Get negative index or bigger than len", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(1, "a")
		ps.Append(4, "d")

		_, _, result1 := ps.Get(-1)
		_, _, result2 := ps.Get(10)
		assert.False(t, result1)
		assert.False(t, result2)
	})

	t.Run("Test First", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(1, "a")
		ps.Append(2, "b")
		ps.Append(3, "c")
		ps.Append(4, "d")

		first, second, _ := ps.First()
		assert.Equal(t, ps.pairs[0].First, first)
		assert.Equal(t, ps.pairs[0].Second, second)
	})

	t.Run("Test Last", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(1, "a")
		ps.Append(2, "b")
		ps.Append(3, "c")
		ps.Append(4, "d")

		first, second, _ := ps.Last()
		assert.Equal(t, ps.pairs[3].First, first)
		assert.Equal(t, ps.pairs[3].Second, second)
	})
}
