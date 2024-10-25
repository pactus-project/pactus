package pairslice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	slice := New[int, string](10)

	assert.NotNil(t, slice)
	assert.Equal(t, 10, cap(slice.pairs))
	assert.Equal(t, 0, len(slice.pairs))
}

func TestPairSlice(t *testing.T) {
	t.Run("Test Append", func(t *testing.T) {
		slice := New[int, string](4)

		slice.Append(1, "a")
		slice.Append(2, "b")
		slice.Append(3, "c")

		assert.Equal(t, 3, slice.Len())
		assert.Equal(t, slice.pairs[2].First, 3)
		assert.Equal(t, slice.pairs[2].Second, "c")
	})

	t.Run("Test RemoveFirst", func(t *testing.T) {
		slice := New[int, string](4)

		slice.Append(1, "a")
		slice.Append(2, "b")
		slice.Append(3, "c")

		slice.RemoveFirst()

		assert.Equal(t, slice.pairs[0].First, 2)
		assert.Equal(t, slice.pairs[0].Second, "b")
	})

	t.Run("Test RemoveLast", func(t *testing.T) {
		slice := New[int, string](4)

		slice.Append(1, "a")
		slice.Append(2, "b")
		slice.Append(3, "c")
		slice.Append(4, "d")

		slice.RemoveLast()

		assert.Equal(t, slice.pairs[2].First, 3)
		assert.Equal(t, slice.pairs[2].Second, "c")
	})

	t.Run("Test Len", func(t *testing.T) {
		slice := New[int, string](4)

		slice.Append(1, "a")
		slice.Append(2, "b")

		assert.Equal(t, 2, slice.Len())
	})

	t.Run("Test Remove", func(t *testing.T) {
		slice := New[int, string](4)

		slice.Append(1, "a")
		slice.Append(2, "b")
		slice.Append(3, "c")
		slice.Append(4, "d")

		slice.remove(1)

		assert.Equal(t, slice.pairs[1].First, 3)
		assert.Equal(t, slice.pairs[1].Second, "c")
	})

	t.Run("Test Get", func(t *testing.T) {
		slice := New[int, string](4)

		slice.Append(1, "a")
		slice.Append(2, "b")
		slice.Append(3, "c")
		slice.Append(4, "d")

		first, second, _ := slice.Get(2)
		assert.Equal(t, slice.pairs[2].First, first)
		assert.Equal(t, slice.pairs[2].Second, second)
	})

	t.Run("Test Get negative index or bigger than len", func(t *testing.T) {
		slice := New[int, string](4)

		slice.Append(1, "a")
		slice.Append(4, "d")

		_, _, result1 := slice.Get(-1)
		_, _, result2 := slice.Get(10)
		assert.False(t, result1)
		assert.False(t, result2)
	})

	t.Run("Test First", func(t *testing.T) {
		slice := New[int, string](4)

		slice.Append(1, "a")
		slice.Append(2, "b")
		slice.Append(3, "c")
		slice.Append(4, "d")

		first, second, _ := slice.First()
		assert.Equal(t, slice.pairs[0].First, first)
		assert.Equal(t, slice.pairs[0].Second, second)
	})

	t.Run("Test Last", func(t *testing.T) {
		slice := New[int, string](4)

		slice.Append(1, "a")
		slice.Append(2, "b")
		slice.Append(3, "c")
		slice.Append(4, "d")

		first, second, _ := slice.Last()
		assert.Equal(t, slice.pairs[3].First, first)
		assert.Equal(t, slice.pairs[3].Second, second)
	})
}
