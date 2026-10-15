package types_test

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/stretchr/testify/assert"
)

func TestHeightSafeSub(t *testing.T) {
	sub := types.Height(10).SafeSub(5)
	assert.Equal(t, uint32(5), sub)

	sub = types.Height(5).SafeSub(10)
	assert.Equal(t, uint32(0), sub)

	sub = types.Height(10).SafeSub(10)
	assert.Equal(t, uint32(0), sub)
}

func TestHeightEncodeAsSlice(t *testing.T) {
	tests := []struct {
		height   types.Height
		expected []byte
	}{
		{height: 0, expected: []byte{0, 0, 0, 0}},
		{height: 1, expected: []byte{1, 0, 0, 0}},
		{height: 255, expected: []byte{0xff, 0, 0, 0}},
		{height: 256, expected: []byte{0, 1, 0, 0}},
		{height: 65535, expected: []byte{0xff, 0xff, 0, 0}},
		{height: 65536, expected: []byte{0, 0, 1, 0}},
	}

	for _, test := range tests {
		slice := test.height.BytesLE()
		height := types.HeightFromBytesLE(slice)

		assert.Equal(t, test.expected, slice)
		assert.Equal(t, test.height, height)
	}
}

func TestHeightSafeIncrease(t *testing.T) {
	assert.Equal(t, types.Height(15), types.Height(10).SafeIncrease(5))
	assert.Equal(t, types.Height(10), types.Height(10).SafeIncrease(0))
	assert.Equal(t, types.Height(1), types.Height(0).SafeIncrease(1))
}

func TestHeightSafeDecrease(t *testing.T) {
	assert.Equal(t, types.Height(5), types.Height(10).SafeDecrease(5))
	assert.Equal(t, types.Height(10), types.Height(10).SafeDecrease(0))

	// The result never goes below zero.
	assert.Equal(t, types.Height(0), types.Height(10).SafeDecrease(10))
	assert.Equal(t, types.Height(0), types.Height(10).SafeDecrease(11))
	assert.Equal(t, types.Height(0), types.Height(0).SafeDecrease(1))
	assert.Equal(t, types.Height(0), types.Height(5).SafeDecrease(1<<32-1))
}
