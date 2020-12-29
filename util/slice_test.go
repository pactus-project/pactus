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
	c, err := CompressSlice(a)
	assert.NoError(t, err)
	b, err := DecompressSlice(c)
	assert.NoError(t, err)
	assert.Equal(t, a, b)
}

func TestDecompress(t *testing.T) {
	d, _ := hex.DecodeString("1f8b08000000000000ff5accb8929191492afefe9620e60805060280254221ac2238cb57f8d6da3ecfc47b617bd47bf80fbe503b11b7aef385a6c0ba159a2142ac110a1d8f2e447cd46a3f3d71d6fc5c9eac45377ec4efffa0b76c33bb1377fead15f5cdf7d9085bc44e58094784c2169fcd92c947ee35a43a49ff5d57b563eeaad9415b8ed6d685bd72aaf9afd3b5898b334455a26edf71fd634957941ead7f15ad5fe0e96517cef48d79216323616702020000ffffa63359ef1b010000")
	_, err := DecompressSlice(d[1:])
	assert.Error(t, err)
	_, err = DecompressSlice(d)
	assert.NoError(t, err)
}
