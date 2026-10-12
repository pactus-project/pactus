package sync

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

// When the last block is in the current block slot, or even ahead of the local clock,
// the node is synced, and it must not try to download blocks.
func TestUpdateBlockchainWhenSynced(t *testing.T) {
	tests := []struct {
		name   string
		offset time.Duration // The time of the last block compared to the current block slot
	}{
		{"last block is in the current block slot", 0},
		{"last block is ahead of the local clock", 30 * time.Second},
		{"last block is far ahead of the local clock", 10 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := setup(t, nil)
			td.addPeer(t, status.StatusKnown, service.New(service.FullNode))

			blockInterval := td.state.Params().BlockInterval()
			blkTime := util.RoundNow(int(blockInterval.Seconds())).Add(tt.offset)
			blk, cert := td.GenerateTestBlock(td.state.LastBlockHeight()+1, testsuite.BlockWithTime(blkTime))
			require.NoError(t, td.state.CommitBlock(blk, cert))

			td.sync.updateBlockchain()

			// There is nothing to download.
			td.shouldNotPublishAnyMessage(t)
		})
	}
}
