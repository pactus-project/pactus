package sync

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/vote"
)

func TestRequestForBlocksInvalidLastBlocHash(t *testing.T) {
	setup(t)

	invHash := crypto.GenerateTestHash()

	// Send block request, but block hash is invalid, ignore it
	tSync.broadcastBlocksReq(7, tState.LastBlockHeight(), invHash)

	tNetAPI.shouldNotReceiveAnyMessageWithThisType(t, payload.PayloadTypeBlocks)
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
	tNetAPI.shouldReceiveMessageWithThisType(t, payload.PayloadTypeVote)

	tSync.broadcastProposal(p)
	tNetAPI.shouldReceiveMessageWithThisType(t, payload.PayloadTypeProposal)

	assert.Equal(t, tConsensus.Votes[0].Hash(), v.Hash())
	assert.Equal(t, tConsensus.Proposal.Hash(), p.Hash())
}

func TestMoveToConsensus(t *testing.T) {
	setup(t)

	// Bad peer send us invalid height
	msg := message.NewSalamMessage(tState.GenHash, 100000000)
	tSync.publishMessage(msg)
	tNetAPI.shouldReceiveMessageWithThisType(t, payload.PayloadTypeAleyk)
	tNetAPI.shouldReceiveMessageWithThisType(t, payload.PayloadTypeBlocksReq)

	tSync.maybeSynced(false)
	assert.False(t, tConsensus.Moved)

	ourHeight := tState.LastBlockHeight()
	// We send all blocks we have and set LastCommit to true
	blocks := make([]*block.Block, 0)
	var commit *block.Commit
	for i := 0; i < 15; i++ {
		b, _ := block.GenerateTestBlock(nil)
		commit = block.GenerateTestCommit(b.Hash())
		blocks = append(blocks, b)

		// To make sure block will be committed
		tCache.AddCommit(b.Hash(), commit)
	}

	assert.NotNil(t, tState.LastBlockCommit)
	tSync.broadcastBlocks(tState.LastBlockHeight()+1, blocks, commit)

	// We send all blocks we have and set LastCommit to true
	tNetAPI.waitingForMessage(t, message.NewBlocksReqMessage(ourHeight+15+1, 100000000, blocks[14].Hash()))

	assert.True(t, tConsensus.Moved)
}

func TestSendInvalidBlock(t *testing.T) {
	setup(t)

	fmt.Println(tState.LastBlockHeight())
	networkHeight := tState.LastBlockHeight() + 15
	msg := message.NewSalamMessage(tState.GenHash, networkHeight)
	tSync.publishMessage(msg)
	tNetAPI.shouldReceiveMessageWithThisType(t, payload.PayloadTypeAleyk)
	tNetAPI.shouldReceiveMessageWithThisType(t, payload.PayloadTypeBlocksReq)

	tSync.maybeSynced(false)
	assert.False(t, tConsensus.Moved)

	ourHeight := tState.LastBlockHeight()
	// We send all blocks we have and set LastCommit to true
	blocks := make([]*block.Block, 0)
	var commit *block.Commit
	for i := 0; i < 15; i++ {
		b, _ := block.GenerateTestBlock(nil)
		commit = block.GenerateTestCommit(b.Hash())
		blocks = append(blocks, b)

		// To make sure block will be committed
		tCache.AddCommit(b.Hash(), commit)
	}

	tState.InvalidBlockHash = blocks[5].Hash()
	assert.NotNil(t, tState.LastBlockCommit)
	tSync.broadcastBlocks(tState.LastBlockHeight()+1, blocks, commit)

	// We send all blocks we have and set LastCommit to true
	tNetAPI.waitingForMessage(t, message.NewBlocksReqMessage(ourHeight+6, networkHeight, blocks[4].Hash()))

	assert.False(t, tConsensus.Moved)
}
