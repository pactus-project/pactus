package block_test

import (
	"sync"
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A block is shared between the services of the node.
// Run with the race detector: the first calls of Hash and Bytes store their result.
func TestConcurrentHashAndBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blk, _ := ts.GenerateTestBlock(ts.RandHeight())
	data, err := blk.Bytes()
	require.NoError(t, err)

	decoded, err := block.FromBytes(data)
	require.NoError(t, err)

	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			assert.Equal(t, blk.Hash(), decoded.Hash())

			got, _ := decoded.Bytes()
			assert.Equal(t, data, got)
		})
	}
	wg.Wait()
}
