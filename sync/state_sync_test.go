package sync

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

func TestAddBlockToCache(t *testing.T) {
	setup(t)

	b1, trxs1 := block.GenerateTestBlock(nil, nil)
	b2, trxs2 := block.GenerateTestBlock(nil, nil)

	// Alice send a block to another peer, bob should cache it
	tAliceSync.stateSync.BroadcastLatestBlocksResponse(payload.ResponseCodeMoreBlocks, tAnotherPeerID, 123, 1001, []*block.Block{b1}, trxs1, nil)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)
	assert.Equal(t, tBobSync.cache.GetBlock(1001).Hash(), b1.Hash())

	// Alice send a block to bob, bob should cache it
	tAliceSync.stateSync.BroadcastLatestBlocksResponse(payload.ResponseCodeMoreBlocks, tBobPeerID, 123, 1002, []*block.Block{b2}, trxs2, nil)
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
	assert.True(t, tBobSync.txPool.HasTx(trx1.ID()))
}
func TestRequestForBlocksVeryFar(t *testing.T) {
	setup(t)

	tAliceSync.stateSync.BroadcastLatestBlocksRequest(tBobPeerID, 2)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	tBobNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)
}

func TestSendLastCommit(t *testing.T) {
	setup(t)

	tAliceSync.stateSync.BroadcastLatestBlocksRequest(tBobPeerID, 95)

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)
	msg := tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)
	pld := msg.Payload.(*payload.LatestBlocksResponsePayload)

	assert.Equal(t, pld.LastCommit, tBobState.LastBlockCommit)
}

func TestPrepareLastBlock(t *testing.T) {
	setup(t)

	h := tAliceState.LastBlockHeight()
	b, _ := tAliceSync.stateSync.prepareBlocksAndTransactions(h, 10)
	assert.Equal(t, len(b), 1)
}

func TestProcessHeartbeatForSyncing(t *testing.T) {
	setup(t)

	lastHash := tAliceState.LastBlockHash()
	height := tAliceState.LastBlockHeight()
	for i := 0; i < 5; i++ {
		b, trxs := block.GenerateTestBlock(nil, &lastHash)
		c := block.GenerateTestCommit(b.Hash())
		tAliceSync.cache.AddTransactions(trxs)
		lastHash = b.Hash()
		assert.NoError(t, tAliceState.ApplyBlock(height+i+1, *b, *c))
	}

	val := validator.NewValidator(tAliceSync.signer.PublicKey(), 4, tAliceState.LastBlockHeight())
	assert.NoError(t, tAliceState.ValSet.UpdateTheSet(0, []*validator.Validator{val}))
	tAliceConsensus.HRS_ = hrs.NewHRS(tAliceState.LastBlockHeight()+1, 0, 3)
	tBobConsensus.Started = false

	tAliceSync.broadcastHeartBeat()
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeHeartBeat)

	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)        // blocks 101-105
	msg := tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // Synced response code
	pld := msg.Payload.(*payload.LatestBlocksResponsePayload)
	assert.Equal(t, pld.ResponseCode, payload.ResponseCodeSynced)
	assert.False(t, tAliceSync.peerSet.HasAnyValidSession())
	assert.False(t, tBobSync.peerSet.HasAnyValidSession())

	assert.True(t, tBobConsensus.Started)
}

func TestDownloadBlock(t *testing.T) {
	setup(t)

	// Clear bob store
	tBobSync.cache.Clear()
	tBobState.Store.Blocks = make(map[int]*block.Block)
	tBobConsensus.Started = false

	tBobSync.sendBlocksRequestIfWeAreBehind()
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 1-10
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 11-20
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 21-31 (one extra block)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // NoMoreBlock

	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 40-49
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 50-59
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 60-70 (one extra block)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // NoMoreBlock

	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 61-70
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 71-80
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 81-91 (one extra block)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // NoMoreBlock

	// Latest block requests
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // 91-100
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // Synced

	assert.True(t, tBobConsensus.Started)
	assert.False(t, tAliceSync.peerSet.HasAnyValidSession())
	assert.False(t, tBobSync.peerSet.HasAnyValidSession())
}

func TestSessionTieout(t *testing.T) {
	tAliceConfig.SessionTimeout = 200 * time.Millisecond
	setup(t)

	p := tAliceSync.peerSet.MustGetPeer(tAnotherPeerID)
	p.UpdateInitialBlockDownload(true)
	p.UpdateHeight(1000)
	tAliceSync.peerSet.UpdateMaxClaimedHeight(1000)
	tAliceSync.sendBlocksRequestIfWeAreBehind()
	assert.True(t, tAliceSync.peerSet.HasAnyValidSession())
	time.Sleep(tAliceConfig.SessionTimeout)
	assert.False(t, tAliceSync.peerSet.HasAnyValidSession())
}
