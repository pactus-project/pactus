package types_test

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/stretchr/testify/assert"
)

func TestRoundSafeSub(t *testing.T) {
	sub := types.Round(10).SafeSub(5)
	assert.Equal(t, uint32(5), sub)

	sub = types.Round(5).SafeSub(10)
	assert.Equal(t, uint32(0), sub)

	sub = types.Round(10).SafeSub(10)
	assert.Equal(t, uint32(0), sub)
}

func TestRoundEncodeAsSlice(t *testing.T) {
	tests := []struct {
		round    types.Round
		expected []byte
	}{
		{round: 0, expected: []byte{0, 0}},
		{round: 1, expected: []byte{1, 0}},
		{round: 255, expected: []byte{0xff, 0}},
		{round: 256, expected: []byte{0, 1}},
		{round: -1, expected: []byte{0xff, 0xff}},
	}

	for _, test := range tests {
		slice := test.round.EncodeAsSlice()
		round := types.RoundFromSlice(slice)

		assert.Equal(t, test.expected, slice)
		assert.Equal(t, test.round, round)
	}
}
