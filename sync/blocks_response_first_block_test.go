package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// A block larger than the bound is still sent: an empty response would tell that the range
// is fully delivered.
func TestPrepareBlocksSendsALargeBlock(t *testing.T) {
	conf := testConfig()
	conf.MaxBlocksDataBytes = 1
	td := setup(t, conf)

	from := td.state.LastBlockHeight() - 3
	blocks := td.sync.prepareBlocks(from, 3)

	assert.Len(t, blocks, 1)
}
