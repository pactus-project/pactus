package util

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandIntsAreNotNegative(t *testing.T) {
	for range 200 {
		assert.GreaterOrEqual(t, RandInt16(0), int16(0))
		assert.GreaterOrEqual(t, RandInt16(-1), int16(0))
		assert.GreaterOrEqual(t, RandInt32(0), int32(0))
		assert.GreaterOrEqual(t, RandInt32(math.MinInt32), int32(0))
		assert.GreaterOrEqual(t, RandInt64(0), int64(0))
		assert.Less(t, RandInt32(10), int32(10))
	}
}

func TestI2OSPOutOfTheDomain(t *testing.T) {
	assert.Nil(t, I2OSP(big.NewInt(256), 1), "256 doesn't fit in one byte")
	assert.Nil(t, I2OSP(big.NewInt(-1), 4))
	assert.Nil(t, I2OSP(big.NewInt(1), -1))
	assert.Equal(t, []byte{0xff}, I2OSP(big.NewInt(255), 1))
	assert.Equal(t, []byte{0, 0, 1, 0}, I2OSP(big.NewInt(256), 4))
}

func TestSliceHelpersDoNotModifyTheCallerSlice(t *testing.T) {
	t.Run("RemoveFirstOccurrenceOf", func(t *testing.T) {
		in := []int{1, 2, 3, 2}
		res, found := RemoveFirstOccurrenceOf(in, 2)

		assert.True(t, found)
		assert.Equal(t, []int{1, 3, 2}, res)
		assert.Equal(t, []int{1, 2, 3, 2}, in)
	})

	t.Run("Prepend", func(t *testing.T) {
		elements := make([]int, 2, 10)
		elements[0], elements[1] = 1, 2
		in := []int{3, 4}
		res := Prepend(in, elements...)
		_ = append(res[:1], 9) // Writing in the result...

		assert.Equal(t, []int{1, 2}, elements[:2], "...doesn't write in the elements")
		assert.Equal(t, []int{3, 4}, in)
	})

	t.Run("Subtracts", func(t *testing.T) {
		in := []int32{1, 2, 3}
		res := Subtracts(in, nil)
		res[0] = 100

		assert.Equal(t, []int32{1, 2, 3}, in)
	})

	t.Run("Trim", func(t *testing.T) {
		assert.NotPanics(t, func() { assert.Empty(t, Trim([]int{1, 2}, -1)) })
		assert.Equal(t, []int{1}, Trim([]int{1, 2}, 1))
		assert.Equal(t, []int{1, 2}, Trim([]int{1, 2}, 5))
	})
}
