package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/vote"
)

func TestRequestForBlocksInvalidLastBlocHash(t *testing.T) {
	setup(t)

	invHash := crypto.GenerateTestHash()

	// Alice asks bob to send blocks but last block hash is invalid
	tAliceSync.broadcastBlocksReq(4, 6, invHash)
	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeBlocksReq)

	tBobNetAPI.shouldNotPublishMessageWithThisType(t, payload.PayloadTypeBlocks)
}

func TestSendLastCommit(t *testing.T) {
	setup(t)

	tAliceSync.broadcastBlocksReq(4, tBobState.LastBlockHeight(), tBobState.Store.Blocks[3].Hash())

	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeBlocksReq)
	msg := tBobNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeBlocks)
	pld := msg.Payload.(*payload.BlocksPayload)

	assert.Equal(t, pld.LastCommit, tBobState.LastBlockCommit)
}

func TestUpdateConsensus(t *testing.T) {
	setup(t)

	v, _ := vote.GenerateTestPrecommitVote(1, 1)
	p, _ := vote.GenerateTestProposal(1, 1)

	tAliceSync.broadcastVote(v)
	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeVote)

	tAliceSync.broadcastProposal(p)
	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeProposal)

	assert.Equal(t, tBobConsensus.Votes[0].Hash(), v.Hash())
	assert.Equal(t, tBobConsensus.Proposal.Hash(), p.Hash())
}

func TestMoveToConsensus(t *testing.T) {
	setup(t)

	aliceHeight := tAliceState.LastBlockHeight()
	aliceLastHash := tAliceState.LastBlockHash()
	// Another peers send all blocks he has and set the LastCommit
	blocks := make([]*block.Block, 0)
	var commit *block.Commit
	lastHash := aliceLastHash
	for i := 0; i < 5; i++ {
		b, _ := block.GenerateTestBlock(nil, &lastHash)
		commit = block.GenerateTestCommit(b.Hash())
		lastHash = b.Hash()
		blocks = append(blocks, b)
	}

	tBobConsensus.Started = false

	tAliceSync.broadcastBlocks(aliceHeight+1, blocks, commit)
	tAliceNetAPI.shouldPublishMessageWithThisType(t, payload.PayloadTypeBlocks)

	assert.True(t, tBobConsensus.Started)
}
