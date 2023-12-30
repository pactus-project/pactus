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

		assert.Equal(t, &Pair[int, string]{2, "c", 2}, ps.pairs[2])
	})

	t.Run("Test RemoveFirst", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		ps.RemoveFirst()

		assert.Equal(t, &Pair[int, string]{1, "b", 1}, ps.pairs[0])
		assert.Equal(t, &Pair[int, string]{2, "c", 2}, ps.pairs[1])
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

		assert.Equal(t, &Pair[int, string]{0, "a", 0}, ps.Get(0))
		assert.Equal(t, &Pair[int, string]{2, "c", 2}, ps.Get(2))
		assert.Nil(t, ps.Get(10))
	})

	t.Run("Test First", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		assert.Equal(t, &Pair[int, string]{0, "a", 0}, ps.First())
	})

	t.Run("Test Last", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		assert.Equal(t, &Pair[int, string]{2, "c", 2}, ps.Last())
	})

	t.Run("Test Last with RemoveFirst", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")
		ps.Append(3, "d")

		ps.RemoveFirst()

		assert.Equal(t, &Pair[int, string]{3, "d", 3}, ps.Last())
	})

	t.Run("Test All", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		all := ps.All()

		assert.Equal(t, all[0].First, 0)
		assert.Equal(t, all[1].First, 1)
		assert.Equal(t, all[2].First, 2)
		assert.Equal(t, all[0].Second, "a")
		assert.Equal(t, all[1].Second, "b")
		assert.Equal(t, all[2].Second, "c")
	})
}
