package util

import (
	"encoding/hex"
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

func TestCompress(t *testing.T) {
	a := []byte{1, 2, 3, 4, 5, 6, 7}
	c, err := CompressBuffer(a)
	assert.NoError(t, err)
	b, err := DecompressBuffer(c)
	assert.NoError(t, err)
	assert.Equal(t, a, b)
}

func TestDecompress(t *testing.T) {
	d, _ := hex.DecodeString("1f8b08000000000000ff5accb8929191492afefe9620e60805060280254221ac2238cb57f8d6da3ecfc47b617bd47bf80fbe503b11b7aef385a6c0ba159a2142ac110a1d8f2e447cd46a3f3d71d6fc5c9eac45377ec4efffa0b76c33bb1377fead15f5cdf7d9085bc44e58094784c2169fcd92c947ee35a43a49ff5d57b563eeaad9415b8ed6d685bd72aaf9afd3b5898b334455a26edf71fd634957941ead7f15ad5fe0e96517cef48d79216323616702020000ffffa63359ef1b010000")
	_, err := DecompressBuffer(d[1:])
	assert.Error(t, err)
	_, err = DecompressBuffer(d)
	assert.NoError(t, err)
}

func TestSubtractAndSubset(t *testing.T) {
	t.Run("Case 1", func(t *testing.T) {
		s1 := []int{1, 2, 3, 4}
		s2 := []int{1, 2, 3}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, s3, []int{4})
	})

	t.Run("Case 2", func(t *testing.T) {
		s1 := []int{1, 2, 3, 4}
		s2 := []int{2, 3, 5}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, s3, []int{1, 4})
	})

	t.Run("Case 3", func(t *testing.T) {
		s1 := []int{1, 2, 3, 4}
		s2 := []int{}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, s3, []int{1, 2, 3, 4})
	})

	t.Run("Case 4", func(t *testing.T) {
		s1 := []int{}
		s2 := []int{1, 2, 3, 4}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, s3, []int{})
	})

	t.Run("Case 5", func(t *testing.T) {
		s1 := []int{1, 2, 3, 4}
		s2 := []int{1, 2, 3, 4}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, s3, []int{})
	})

	t.Run("Case 6", func(t *testing.T) {
		s1 := []int{1, 3, 5}
		s2 := []int{1, 2, 3, 4, 5}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, s3, []int{})
	})

	t.Run("Case 7", func(t *testing.T) {
		s1 := []int{1, 2, 3, 4}
		s3 := Subtracts(s1, nil)
		assert.Equal(t, s3, s1)
	})

	t.Run("Case 8", func(t *testing.T) {
		s2 := []int{1, 2, 3, 4}
		s3 := Subtracts(nil, s2)
		assert.Equal(t, s3, []int{})
	})
}

func TestEqual(t *testing.T) {
	assert.True(t, Equal([]int{1, 2, 3}, []int{1, 2, 3}))
	assert.False(t, Equal([]int{1, 2, 3}, []int{1, 3, 2}))
	assert.False(t, Equal([]int{1, 2, 3}, []int{1, 2, 3, 4}))
	assert.True(t, Equal([]int{}, []int{}))
	assert.True(t, Equal([]int{}, nil))
}

func TestContains(t *testing.T) {
	assert.True(t, Contains([]int{1, 2, 3, 4}, 2))
	assert.False(t, Contains([]int{1, 2, 3, 4}, 5))
	assert.False(t, Contains([]int{}, 0))
}
