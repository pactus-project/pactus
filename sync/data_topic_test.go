package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/tx"
)

func TestAddBlockToCache(t *testing.T) {
	setup(t)

	b1, trxs := block.GenerateTestBlock(nil, nil)

	// Alice send block to bob, bob should cache it
	tAliceSync.dataTopic.BroadcastLatestBlocks(1001, []*block.Block{b1}, trxs, nil)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocks)
	assert.Equal(t, tBobSync.cache.GetBlock(1001).Hash(), b1.Hash())
}

func TestAddTxToCache(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestBondTx()

	// Alice send transaction to bob, bob should cache it
	tAliceSync.dataTopic.BroadcastTransactions([]*tx.Tx{trx1})
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeTransactions)
	assert.NotNil(t, tBobSync.cache.GetTransaction(trx1.ID()))
}

func TestSendTxs(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestBondTx()
	trx2, _ := tx.GenerateTestSendTx()

	// Alice has trx1 in his cache
	tAliceSync.cache.AddTransaction(trx1)
	tBobSync.cache.AddTransaction(trx2)

	tAliceBroadcastCh <- message.NewTransactionsRequestMessage([]crypto.Hash{trx1.ID()})
	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeTransactionsRequest)

	tAliceSync.cache.AddTransaction(trx1)
	tAliceBroadcastCh <- message.NewTransactionsRequestMessage([]crypto.Hash{trx1.ID(), trx2.ID()})
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeTransactionsRequest)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeTransactions)

	assert.NotNil(t, tAliceSync.cache.GetTransaction(trx2.ID()))
}

func TestRequestForBlocksInvalidLastBlocHash(t *testing.T) {
	setup(t)

	invHash := crypto.GenerateTestHash()

	// Alice asks bob to send blocks but last block hash is invalid
	tAliceSync.dataTopic.BroadcastLatestBlocksRequest(8, invHash)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)

	tBobNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocks)
}

func TestRequestForBlocksVeryFar(t *testing.T) {
	setup(t)

	tAliceSync.dataTopic.BroadcastLatestBlocksRequest(2, tBobState.Store.Blocks[1].Hash())
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)

	tBobNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocks)
}

func TestSendLastCommit(t *testing.T) {
	setup(t)

	tAliceSync.dataTopic.BroadcastLatestBlocksRequest(8, tBobState.Store.Blocks[7].Hash())

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	msg := tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocks)
	pld := msg.Payload.(*payload.LatestBlocksPayload)

	assert.Equal(t, pld.Commit, tBobState.LastBlockCommit)
}

func TestMoveToConsensus(t *testing.T) {
	setup(t)

	aliceHeight := tAliceState.LastBlockHeight()
	aliceLastHash := tAliceState.LastBlockHash()
	// Another peers send all blocks he has and set the LastCommit
	blocks := make([]*block.Block, 0)
	trxs := make([]*tx.Tx, 0)
	var commit *block.Commit
	lastHash := aliceLastHash
	for i := 0; i < 5; i++ {
		b, t := block.GenerateTestBlock(nil, &lastHash)
		commit = block.GenerateTestCommit(b.Hash())
		lastHash = b.Hash()
		blocks = append(blocks, b)
		trxs = append(trxs, t...)
	}

	tBobConsensus.Started = false

	tAliceSync.dataTopic.BroadcastLatestBlocks(aliceHeight+1, blocks, trxs, commit)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocks)

	assert.True(t, tBobConsensus.Started)
}
