package tripleslice

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

		assert.Equal(t, &Triple[int, string]{2, "c", 2}, ps.pairs[2])
	})

	t.Run("Test Pop", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		ps.Pop()

		assert.Equal(t, &Triple[int, string]{1, "b", 1}, ps.pairs[0])
		assert.Equal(t, &Triple[int, string]{2, "c", 2}, ps.pairs[1])
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

		assert.Equal(t, &Triple[int, string]{0, "a", 0}, ps.Get(0))
		assert.Equal(t, &Triple[int, string]{2, "c", 2}, ps.Get(2))
		assert.Nil(t, ps.Get(10))
	})

	t.Run("Test First", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		assert.Equal(t, &Triple[int, string]{0, "a", 0}, ps.First())
	})

	t.Run("Test Last", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		assert.Equal(t, &Triple[int, string]{2, "c", 2}, ps.Last())
	})

	t.Run("Test Last with Pop", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")
		ps.Append(3, "d")

		ps.Pop()

		assert.Equal(t, &Triple[int, string]{3, "d", 3}, ps.Last())
	})

	t.Run("Test All", func(t *testing.T) {
		ps := New[int, string](4)

		ps.Append(0, "a")
		ps.Append(1, "b")
		ps.Append(2, "c")

		all := ps.All()

		assert.Equal(t, all[0].FirstElement, 0)
		assert.Equal(t, all[1].FirstElement, 1)
		assert.Equal(t, all[2].FirstElement, 2)
		assert.Equal(t, all[0].SecondElement, "a")
		assert.Equal(t, all[1].SecondElement, "b")
		assert.Equal(t, all[2].SecondElement, "c")
	})
}
