package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
)

func TestAddBlockToCache(t *testing.T) {
	setup(t)

	b1, trxs1 := block.GenerateTestBlock(nil, nil)
	b2, trxs2 := block.GenerateTestBlock(nil, nil)

	// Alice send a block to another peer, bob should cache it
	tAliceSync.stateSync.BroadcastLatestBlocksResponse(0, 1001, []*block.Block{b1}, trxs1, nil)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)
	assert.Equal(t, tBobSync.cache.GetBlock(1001).Hash(), b1.Hash())

	// Alice send a block to bob, bob should cache it
	tAliceSync.stateSync.BroadcastLatestBlocksResponse(0, 1002, []*block.Block{b2}, trxs2, nil)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)
	assert.Equal(t, tBobSync.cache.GetBlock(1002).Hash(), b2.Hash())
}

func TestAddTxToCache(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestBondTx()

	// Alice send transaction to bob, bob should cache it
	tAliceSync.stateSync.BroadcastTransactions([]*tx.Tx{trx1})
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

	tAliceBroadcastCh <- message.NewQueryTransactionsMessage([]crypto.Hash{trx1.ID()})
	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeQueryTransactions)

	tAliceSync.cache.AddTransaction(trx1)
	tAliceBroadcastCh <- message.NewQueryTransactionsMessage([]crypto.Hash{trx1.ID(), trx2.ID()})
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryTransactions)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeTransactions)

	assert.NotNil(t, tAliceSync.cache.GetTransaction(trx2.ID()))
}

func TestRequestForBlocksVeryFar(t *testing.T) {
	setup(t)

	tAliceSync.stateSync.BroadcastLatestBlocksRequest(tBobPeerID, 2)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	tBobNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)
}

func TestSendLastCommit(t *testing.T) {
	setup(t)

	tAliceSync.stateSync.BroadcastLatestBlocksRequest(tBobPeerID, 90)

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	msg := tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)
	pld := msg.Payload.(*payload.LatestBlocksResponsePayload)

	assert.Equal(t, pld.LastCommit, tBobState.LastBlockCommit)
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

	tAliceSync.stateSync.BroadcastLatestBlocksResponse(0, aliceHeight+1, blocks, trxs, commit)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)

	assert.True(t, tBobConsensus.Started)
}

func TestPrepareLastBlock(t *testing.T) {
	setup(t)

	h := tAliceState.LastBlockHeight()
	b, _ := tAliceSync.stateSync.prepareBlocksAndTransactions(h, 10)
	assert.Equal(t, len(b), 1)
}

func TestDownloadBlock(t *testing.T) {
	setup(t)

	// Clear bob store
	tBobState.Store.Blocks = make(map[int]*block.Block)

	tBobSync.sendBlocksRequestIfWeAreBehind()
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse)

	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse)

	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse)

	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse)
}
