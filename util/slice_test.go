package util

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceToInt16(t *testing.T) {
	tests := []struct {
		in    int16
		slice []byte
	}{
		{MinInt16, []byte{0x00, 0x80}},
		{int16(-128), []byte{0x80, 0xff}},
		{int16(-1), []byte{0xff, 0xff}},
		{int16(0), []byte{0x00, 0x00}},
		{int16(1), []byte{0x01, 0x00}},
		{int16(256), []byte{0x00, 0x01}},
		{MaxInt16, []byte{0xff, 0x7f}},
	}

	for _, test := range tests {
		s1 := Uint16ToSlice(uint16(test.in))
		s2 := Int16ToSlice(test.in)
		assert.Equal(t, s1, s2)
		assert.Equal(t, s1, test.slice)

		v1 := SliceToInt16(test.slice)
		v2 := SliceToUint16(test.slice)
		assert.Equal(t, v1, int16(v2))
		assert.Equal(t, v1, test.in)
	}
}

func TestSliceToInt32(t *testing.T) {
	tests := []struct {
		in    int32
		slice []byte
	}{
		{MinInt32, []byte{0x00, 0x00, 0x00, 0x80}},
		{int32(-128), []byte{0x80, 0xff, 0xff, 0xff}},
		{int32(-1), []byte{0xff, 0xff, 0xff, 0xff}},
		{int32(0), []byte{0x00, 0x00, 0x00, 0x00}},
		{int32(1), []byte{0x01, 0x00, 0x00, 0x00}},
		{int32(256), []byte{0x00, 0x01, 0x00, 0x00}},
		{MaxInt32, []byte{0xff, 0xff, 0xff, 0x7f}},
	}

	for _, test := range tests {
		s1 := Uint32ToSlice(uint32(test.in))
		s2 := Int32ToSlice(test.in)
		assert.Equal(t, s1, s2)
		assert.Equal(t, s1, test.slice)

		v1 := SliceToInt32(test.slice)
		v2 := SliceToUint32(test.slice)
		assert.Equal(t, v1, int32(v2))
		assert.Equal(t, v1, test.in)
	}
}

func TestSliceToInt64(t *testing.T) {
	tests := []struct {
		in    int64
		slice []byte
	}{
		{MinInt64, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}},
		{int64(-128), []byte{0x80, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{int64(-1), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{int64(0), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{int64(1), []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{int64(256), []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{MaxInt64, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}},
	}

	for _, test := range tests {
		s1 := Uint64ToSlice(uint64(test.in))
		s2 := Int64ToSlice(test.in)
		assert.Equal(t, s1, s2)
		assert.Equal(t, s1, test.slice)

		v1 := SliceToInt64(test.slice)
		v2 := SliceToUint64(test.slice)
		assert.Equal(t, v1, int64(v2))
		assert.Equal(t, v1, test.in)
	}
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
	d, _ := hex.DecodeString(
		"1f8b08000000000000ff5accb8929191492afefe9620e60805060280254221ac2238cb57f8d6da3ecfc47b617bd47bf80fbe503b11b7aef3" +
			"85a6c0ba159a2142ac110a1d8f2e447cd46a3f3d71d6fc5c9eac45377ec4efffa0b76c33bb1377fead15f5cdf7d9085bc44e58094784c216" +
			"9fcd92c947ee35a43a49ff5d57b563eeaad9415b8ed6d685bd72aaf9afd3b5898b334455a26edf71fd634957941ead7f15ad5fe0e96517ce" +
			"f48d79216323616702020000ffffa63359ef1b010000")
	_, err := DecompressBuffer(d[1:])
	assert.Error(t, err)
	_, err = DecompressBuffer(d)
	assert.NoError(t, err)
}

func TestSubtractAndSubset(t *testing.T) {
	t.Run("Case 1", func(t *testing.T) {
		s1 := []int32{1, 2, 3, 4}
		s2 := []int32{1, 2, 3}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, s3, []int32{4})
	})

	t.Run("Case 2", func(t *testing.T) {
		s1 := []int32{1, 2, 3, 4}
		s2 := []int32{2, 3, 5}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, s3, []int32{1, 4})
	})

	t.Run("Case 3", func(t *testing.T) {
		s1 := []int32{1, 2, 3, 4}
		s2 := []int32{}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, s3, []int32{1, 2, 3, 4})
	})

	t.Run("Case 4", func(t *testing.T) {
		s1 := []int32{}
		s2 := []int32{1, 2, 3, 4}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, s3, []int32{})
	})

	t.Run("Case 5", func(t *testing.T) {
		s1 := []int32{1, 2, 3, 4}
		s2 := []int32{1, 2, 3, 4}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, s3, []int32{})
	})

	t.Run("Case 6", func(t *testing.T) {
		s1 := []int32{1, 3, 5}
		s2 := []int32{1, 2, 3, 4, 5}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, s3, []int32{})
	})

	t.Run("Case 7", func(t *testing.T) {
		s1 := []int32{1, 2, 3, 4}
		s3 := Subtracts(s1, nil)
		assert.Equal(t, s3, s1)
	})

	t.Run("Case 8", func(t *testing.T) {
		s2 := []int32{1, 2, 3, 4}
		s3 := Subtracts(nil, s2)
		assert.Equal(t, s3, []int32{})
	})
}

func TestEqual(t *testing.T) {
	assert.True(t, Equal([]int32{1, 2, 3}, []int32{1, 2, 3}))
	assert.False(t, Equal([]int32{1, 2, 3}, []int32{1, 3, 2}))
	assert.False(t, Equal([]int32{1, 2, 3}, []int32{1, 2, 3, 4}))
	assert.True(t, Equal([]int32{}, []int32{}))
	assert.True(t, Equal([]int32{}, nil))
}

func TestContains(t *testing.T) {
	assert.True(t, Contains([]int32{1, 2, 3, 4}, 2))
	assert.False(t, Contains([]int{1, 2, 3, 4}, 5))
	assert.False(t, Contains([]int64{}, 0))
	assert.True(t, Contains([]string{"foo", "bar"}, "foo"))
	assert.False(t, Contains([]string{"foo", "bar"}, "zoo"))
}

func TestSafeCmp(t *testing.T) {
	assert.True(t, SafeCmp([]byte{1, 2, 3}, []byte{1, 2, 3}))
	assert.False(t, SafeCmp([]byte{1, 2, 3, 3}, []byte{1, 2, 3}))
}

func TestMerge(t *testing.T) {
	tests := []struct {
		slices [][]byte
		merged []byte
	}{
		{[][]byte{nil}, []byte{}},
		{[][]byte{{0, 1, 2}}, []byte{0, 1, 2}},
		{[][]byte{{}}, []byte{}},
		{[][]byte{{}, {}}, []byte{}},
		{[][]byte{{0}, {0}}, []byte{0, 0}},
		{[][]byte{{0}, {1}, {2}}, []byte{0, 1, 2}},
	}

	for _, test := range tests {
		merged := Merge(test.slices...)
		assert.Equal(t, merged, test.merged)
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		slice    []byte
		reversed []byte
	}{
		{[]byte{}, []byte{}},
		{[]byte{0}, []byte{0}},
		{[]byte{1, 2, 3}, []byte{3, 2, 1}},
		{[]byte{1, 2}, []byte{2, 1}},
	}

	for _, test := range tests {
		Reverse(test.slice)
		assert.Equal(t, test.slice, test.reversed)
	}
}
