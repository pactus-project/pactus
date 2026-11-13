package grpc

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/stretchr/testify/assert"
)

func TestConfirmationsOf(t *testing.T) {
	tests := []struct {
		last, height types.Height
		want         int
	}{
		{last: 5, height: 5, want: 1}, // the last block is final
		{last: 5, height: 4, want: 2},
		{last: 5, height: 1, want: 5},
		{last: 5, height: 0, want: 0}, // not committed
		{last: 5, height: 6, want: 0}, // not committed yet
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, confirmationsOf(tt.last, tt.height), "last %d, height %d", tt.last, tt.height)
	}
}
