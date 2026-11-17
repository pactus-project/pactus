package block_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/stretchr/testify/require"
)

// A block that decodes is encoded back to the same bytes: it has one encoding.
func FuzzBlockRoundTrip(f *testing.F) {
	// The seeds are in testdata/fuzz/FuzzBlockRoundTrip.
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		blk, err := block.FromBytes(data)
		if err != nil {
			return
		}

		encoded, err := blk.Bytes()
		require.NoError(t, err)
		require.Equal(t, data, encoded)
	})
}
