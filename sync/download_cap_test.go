package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The sessions don't request the blocks after the estimated height of the network:
// a peer would reject them, and they would be requested again.
func TestDownloadStopsAtTheEstimatedHeight(t *testing.T) {
	conf := testConfig()
	conf.MaxSessions = 4
	conf.BlockPerSession = 27
	td := setup(t, conf)
	for range 4 {
		td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	}

	// The last block is 5 block intervals old: about 5 blocks are missing.
	blockInterval := td.state.Params().BlockInterval()
	blkTime := util.RoundNow(int(blockInterval.Seconds())).Add(-5 * blockInterval)
	blk, cert := td.GenerateTestBlock(td.state.LastBlockHeight()+1, testsuite.BlockWithTime(blkTime))
	require.NoError(t, td.state.CommitBlock(blk, cert))
	lastHeight := td.state.LastBlockHeight()

	td.sync.updateBlockchain()

	bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksRequest)
	req := bdl.Message.(*message.BlocksRequestMessage)
	assert.Equal(t, lastHeight+1, req.From)
	assert.LessOrEqual(t, uint32(req.From)+req.Count-1, uint32(lastHeight)+5)
	td.shouldNotPublishAnyMessage(t)
	assert.Equal(t, 1, td.sync.peerSet.NumberOfSessions())
}
