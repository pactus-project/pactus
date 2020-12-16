package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceToInt(t *testing.T) {
	i1 := -1
	s := IntToSlice(i1)
	i2 := SliceToInt(s)
	assert.Equal(t, i1, i2)
}

func TestSliceToUInt(t *testing.T) {
	i1 := uint(0)
	i1--
	s := UIntToSlice(i1)
	i2 := SliceToUInt(s)
	assert.Equal(t, i1, i2)
}

func TestSliceToInt64(t *testing.T) {
	i1 := MaxInt64
	s := Int64ToSlice(i1)
	i2 := SliceToInt64(s)
	assert.Equal(t, i1, i2)
}

func TestSliceToUInt64(t *testing.T) {
	i1 := MaxUint64
	s := UInt64ToSlice(i1)
	i2 := SliceToUInt64(s)
	assert.Equal(t, i1, i2)
}
