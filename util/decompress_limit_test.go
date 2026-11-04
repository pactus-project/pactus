package util

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecompressBufferHugeLimit(t *testing.T) {
	c, err := CompressBuffer([]byte("pactus"))
	require.NoError(t, err)

	// A limit of MaxInt must not overflow into an empty result.
	out, err := DecompressBuffer(c, WithMaxDecompressedSize(math.MaxInt))
	require.NoError(t, err)
	assert.Equal(t, []byte("pactus"), out)
}
