package pairslice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ps := New[int, string](10)

	assert.NotNil(t, ps)
	assert.Equal(t, 10, len(ps.pairs))
	assert.Equal(t, 0, ps.index)
	assert.Equal(t, 10, ps.capacity)
}

func TestPairSlice(t *testing.T) {
	t.Run("Test Append", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		assert.Equal(t, &Pair[int, string]{2, "c"}, ps.pairs[2])
	})

	t.Run("Test Pop", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		ps.Pop()

		assert.Equal(t, &Pair[int, string]{1, "b"}, ps.pairs[0])
		assert.Equal(t, &Pair[int, string]{2, "c"}, ps.pairs[1])
		assert.Nil(t, ps.pairs[2])
	})

	t.Run("Test Has", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		assert.True(t, ps.Has(0))
		assert.False(t, ps.Has(3))
		assert.False(t, ps.Has(10))
	})

	t.Run("Test Get", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		assert.Equal(t, &Pair[int, string]{0, "a"}, ps.Get(0))
		assert.Equal(t, &Pair[int, string]{2, "c"}, ps.Get(2))
		assert.Nil(t, ps.Get(10))
	})

	t.Run("Test First", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		assert.Equal(t, &Pair[int, string]{0, "a"}, ps.First())
	})

	t.Run("Test Last", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		assert.Equal(t, &Pair[int, string]{2, "c"}, ps.Last())
	})

	t.Run("Test Last with Pop", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")
		ps.Append(3, "d")

		ps.Pop()

		assert.Equal(t, &Pair[int, string]{3, "d"}, ps.Last())
	})
}
