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
		slice := test.height.EncodeAsSlice()
		height := types.HeightFromSlice(slice)

		assert.Equal(t, test.expected, slice)
		assert.Equal(t, test.height, height)
	}
}
