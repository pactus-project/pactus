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

	for _, tt := range tests {
		s1 := Uint16ToSlice(uint16(tt.in))
		s2 := Int16ToSlice(tt.in)
		assert.Equal(t, s1, s2)
		assert.Equal(t, tt.slice, s1)

		v1 := SliceToInt16(tt.slice)
		v2 := SliceToUint16(tt.slice)
		assert.Equal(t, int16(v2), v1)
		assert.Equal(t, tt.in, v1)
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

	for _, tt := range tests {
		s1 := Uint32ToSlice(uint32(tt.in))
		s2 := Int32ToSlice(tt.in)
		assert.Equal(t, s1, s2)
		assert.Equal(t, tt.slice, s1)

		v1 := SliceToInt32(tt.slice)
		v2 := SliceToUint32(tt.slice)
		assert.Equal(t, int32(v2), v1)
		assert.Equal(t, tt.in, v1)
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

	for _, tt := range tests {
		s1 := Uint64ToSlice(uint64(tt.in))
		s2 := Int64ToSlice(tt.in)
		assert.Equal(t, s1, s2)
		assert.Equal(t, tt.slice, s1)

		v1 := SliceToInt64(tt.slice)
		v2 := SliceToUint64(tt.slice)
		assert.Equal(t, int64(v2), v1)
		assert.Equal(t, tt.in, v1)
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
	data, _ := hex.DecodeString(
		"1f8b08000000000000ff5accb8929191492afefe9620e60805060280254221ac2238cb57f8d6da3ecfc47b617bd47bf80fbe503b11b7aef3" +
			"85a6c0ba159a2142ac110a1d8f2e447cd46a3f3d71d6fc5c9eac45377ec4efffa0b76c33bb1377fead15f5cdf7d9085bc44e58094784c216" +
			"9fcd92c947ee35a43a49ff5d57b563eeaad9415b8ed6d685bd72aaf9afd3b5898b334455a26edf71fd634957941ead7f15ad5fe0e96517ce" +
			"f48d79216323616702020000ffffa63359ef1b010000")
	_, err := DecompressBuffer(data[1:])
	assert.Error(t, err)
	_, err = DecompressBuffer(data)
	assert.NoError(t, err)
}

func TestSubtractAndSubset(t *testing.T) {
	t.Run("Case 1", func(t *testing.T) {
		s1 := []int32{1, 2, 3, 4}
		s2 := []int32{1, 2, 3}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, []int32{4}, s3)
	})

	t.Run("Case 2", func(t *testing.T) {
		s1 := []int32{1, 2, 3, 4}
		s2 := []int32{2, 3, 5}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, []int32{1, 4}, s3)
	})

	t.Run("Case 3", func(t *testing.T) {
		s1 := []int32{1, 2, 3, 4}
		s2 := []int32{}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, []int32{1, 2, 3, 4}, s3)
	})

	t.Run("Case 4", func(t *testing.T) {
		s1 := []int32{}
		s2 := []int32{1, 2, 3, 4}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, []int32{}, s3)
	})

	t.Run("Case 5", func(t *testing.T) {
		s1 := []int32{1, 2, 3, 4}
		s2 := []int32{1, 2, 3, 4}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, []int32{}, s3)
	})

	t.Run("Case 6", func(t *testing.T) {
		s1 := []int32{1, 3, 5}
		s2 := []int32{1, 2, 3, 4, 5}
		s3 := Subtracts(s1, s2)
		assert.Equal(t, []int32{}, s3)
	})

	t.Run("Case 7", func(t *testing.T) {
		s1 := []int32{1, 2, 3, 4}
		s3 := Subtracts(s1, nil)
		assert.Equal(t, s3, s1)
	})

	t.Run("Case 8", func(t *testing.T) {
		s2 := []int32{1, 2, 3, 4}
		s3 := Subtracts(nil, s2)
		assert.Equal(t, []int32{}, s3)
	})
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

	for _, tt := range tests {
		merged := Merge(tt.slices...)
		assert.Equal(t, tt.merged, merged)
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

	for _, tt := range tests {
		Reverse(tt.slice)
		assert.Equal(t, tt.slice, tt.reversed)
	}
}

func TestPadToLeft(t *testing.T) {
	tests := []struct {
		in   []int
		size int
		want []int
	}{
		{[]int{1, 2, 3}, 5, []int{0, 0, 1, 2, 3}},
		{[]int{1, 2, 3}, 3, []int{1, 2, 3}},
		{[]int{1, 2, 3}, 2, []int{1, 2, 3}},
		{[]int{}, 5, []int{0, 0, 0, 0, 0}},
		{[]int{}, 0, []int{}},
	}

	for _, tt := range tests {
		got := PadToLeft(tt.in, tt.size)
		assert.Equal(t, tt.want, got, "PadToLeft failed, got %v, want %v", got, tt.want)
	}
}

func TestPadToRight(t *testing.T) {
	tests := []struct {
		in   []int
		size int
		want []int
	}{
		{[]int{1, 2, 3}, 5, []int{1, 2, 3, 0, 0}},
		{[]int{1, 2, 3}, 3, []int{1, 2, 3}},
		{[]int{1, 2, 3}, 2, []int{1, 2, 3}},
		{[]int{}, 4, []int{0, 0, 0, 0}},
		{[]int{}, 0, []int{}},
	}

	for _, tt := range tests {
		got := PadToRight(tt.in, tt.size)
		assert.Equal(t, tt.want, got, "PadToRight failed, got %v, want %v", got, tt.want)
	}
}

func TestIsSubset(t *testing.T) {
	tests := []struct {
		arr1, arr2 []int
		want       bool
	}{
		{[]int{11, 1, 13, 21, 3, 7}, []int{11, 3, 7}, true},
		{[]int{11, 1, 13, 21, 3, 7}, []int{3, 11, 7}, false},
		{[]int{1, 2, 3, 4, 5}, []int{1, 2, 3}, true},
		{[]int{1, 2, 3}, []int{1, 2, 3, 4}, false},
		{[]int{1, 2, 3, 4, 5}, []int{6, 7, 8}, false},
		{[]int{}, []int{1, 2, 3, 4, 5}, false},
		{[]int{1, 2, 3, 4, 5}, []int{}, true},
		{[]int{}, []int{}, true},
	}

	for _, tt := range tests {
		got := IsSubset(tt.arr1, tt.arr2)
		assert.Equal(t, tt.want, got,
			"isSubset(%v, %v) = %v; want %v", tt.arr1, tt.arr2, got, tt.want)
	}
}

func TestStringToBytes(t *testing.T) {
	tests := []struct {
		input  string
		output []byte
	}{
		{"Hello", []byte("Hello")},
		{"Go", []byte("Go")},
		{"", []byte("")},
	}

	for _, tt := range tests {
		got := StringToBytes(tt.input)
		assert.Equal(t, tt.output, got, "StringToBytes('%s') =  %v, want %v", tt.input, got, tt.output)
	}
}

func TestRemoveFirstOccurrenceOf(t *testing.T) {
	tests := []struct {
		name    string
		s       []int
		e       int
		want    []int
		removed bool
	}{
		{
			name:    "empty slice",
			s:       []int{},
			e:       1,
			want:    []int{},
			removed: false,
		},
		{
			name:    "element not in slice",
			s:       []int{1, 2, 3},
			e:       4,
			want:    []int{1, 2, 3},
			removed: false,
		},
		{
			name:    "element in slice",
			s:       []int{1, 2, 3},
			e:       2,
			want:    []int{1, 3},
			removed: true,
		},
		{
			name:    "two elements in slice",
			s:       []int{1, 2, 2, 3},
			e:       2,
			want:    []int{1, 2, 3},
			removed: true,
		},
	}

	for _, tt := range tests {
		got, removed := RemoveFirstOccurrenceOf(tt.s, tt.e)

		assert.Equal(t, tt.want, got, "%s failed: got %v, want %v", tt.name, got, tt.want)
		assert.Equal(t, tt.removed, removed, "%s failed: got %v, want %v", tt.name, removed, tt.removed)
	}
}

func TestTrimSlice(t *testing.T) {
	tests := []struct {
		input     []int
		newLength int
		want      []int
	}{
		{[]int{1, 2, 3, 4, 5}, 3, []int{1, 2, 3}},
		{[]int{1}, 3, []int{1}},
		{[]int{}, 3, []int{}},
		{[]int{1, 2, 3, 4, 5}, 0, []int{}},
	}

	for _, tt := range tests {
		got := Trim(tt.input, tt.newLength)
		assert.Equal(t, tt.want, got, "Trim() = %v, want %v", got, tt.want)
	}
}

func TestShuffle(t *testing.T) {
	// Create a slice with 100 integers
	ints := make([]int, 100)
	for i := range ints {
		ints[i] = i + 1
	}
	originalInts := make([]int, len(ints))
	copy(originalInts, ints)

	Shuffle(ints)

	assert.NotEqual(t, originalInts, ints, "ints slice was not shuffled")
	assert.ElementsMatch(t, originalInts, ints, "ints slice does not contain the same elements")
}
