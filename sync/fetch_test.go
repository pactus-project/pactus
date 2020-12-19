package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/vote"
)

func TestRequestForBlocksInvalidLastBlocHash(t *testing.T) {
	setup(t)

	invHash := crypto.GenerateTestHash()

	// Send block request, but block hash is invalid, ignore it
	tSync.broadcastBlocksReq(7, tState.LastBlockHeight(), invHash)

	tNetAPI.shouldNotReceiveAnyMessageWithThisType(t, message.PayloadTypeBlocks)
}

func TestRequestForBlocks(t *testing.T) {
	setup(t)

	h := tState.Store.Blocks[7].Header().LastBlockHash()
	tSync.broadcastBlocksReq(7, 11, h)

	blocks := make([]*block.Block, 0)
	for i := 7; i <= 11; i++ {
		blocks = append(blocks, tState.Store.Blocks[i])
	}

	expectedMsg := message.NewBlocksMessage(7, blocks, nil)
	tNetAPI.waitingForMessage(t, expectedMsg)
}

func TestRequestForBlocksWithLastCommist(t *testing.T) {
	setup(t)

	h := tState.Store.Blocks[7].Header().LastBlockHash()
	tSync.broadcastBlocksReq(7, tState.LastBlockHeight(), h)

	blocks := make([]*block.Block, 0)
	for i := 7; i <= tState.LastBlockHeight(); i++ {
		blocks = append(blocks, tState.Store.Blocks[i])
	}

	assert.NotNil(t, tState.LastBlockCommit)
	expectedMsg := message.NewBlocksMessage(7, blocks, tState.LastBlockCommit)
	tNetAPI.waitingForMessage(t, expectedMsg)
}

func TestUpdateConsensus(t *testing.T) {
	setup(t)

	v, _ := vote.GenerateTestPrecommitVote(1, 1)
	p, _ := vote.GenerateTestProposal(1, 1)

	tSync.broadcastVote(v)
	tNetAPI.shouldReceiveMessageWithThisType(t, message.PayloadTypeVote)

	tSync.broadcastProposal(p)
	tNetAPI.shouldReceiveMessageWithThisType(t, message.PayloadTypeProposal)

	assert.Equal(t, tConsensus.Votes[0].Hash(), v.Hash())
	assert.Equal(t, tConsensus.Proposal.Hash(), p.Hash())
}
