package sync

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
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

func TestRequestForBlocksNotVeryFar(t *testing.T) {
	setup(t)

	addMoreBlocksForBobAndSendBlockAnnounceMessage(t, 15)

	tAliceSync.stateSync.BroadcastLatestBlocksRequest(tBobPeerID)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // blocks 21-30
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // blocks 31-35
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // last commit + sync response
}

func TestRequestForBlocksVeryFar(t *testing.T) {
	setup(t)
	addMoreBlocksForBobAndSendBlockAnnounceMessage(t, 40)

	tAliceSync.stateSync.BroadcastLatestBlocksRequest(tBobPeerID)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	tBobNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)
}

func TestSendLastCommit(t *testing.T) {
	setup(t)

	addMoreBlocksForBobAndSendBlockAnnounceMessage(t, 5)

	tAliceSync.stateSync.BroadcastLatestBlocksRequest(tBobPeerID)
	tAliceConsensus.Started = false

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)
	msg := tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse)
	pld := msg.Payload.(*payload.LatestBlocksResponsePayload)

	assert.Equal(t, pld.LastCommit, tBobState.LastBlockCommit)
	assert.True(t, tAliceConsensus.Started)
}

func TestPrepareLastBlock(t *testing.T) {
	setup(t)

	h := tAliceState.LastBlockHeight()
	b, _ := tAliceSync.stateSync.prepareBlocksAndTransactions(h, 10)
	assert.Equal(t, len(b), 1)
}

func TestMoveToConcensusByBlockAnnounceMessage(t *testing.T) {
	setup(t)

	tAliceConsensus.Started = false
	joinBobToTheSetAndUpdateHRS(t)
	addMoreBlocksForBobAndSendBlockAnnounceMessage(t, 1)

	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeBlockAnnounce)

	assert.True(t, tAliceConsensus.Started)
}
func TestInvalidRangeForDownload(t *testing.T) {
	setup(t)

	t.Run("Bob is not target", func(t *testing.T) {
		pld := &payload.DownloadRequestPayload{
			SessionID: 1,
			Initiator: tAnotherPeerID,
			Target:    util.RandomPeerID(),
			From:      1000,
			To:        1001,
		}
		tBobSync.stateSync.ProcessDownloadRequestPayload(pld)
		tBobNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse)
	})

	t.Run("Ask Bob to send big range of blocks", func(t *testing.T) {
		pld := &payload.DownloadRequestPayload{
			SessionID: 1,
			Initiator: tAnotherPeerID,
			Target:    tBobPeerID,
			From:      1000,
			To:        2000,
		}
		tBobSync.stateSync.ProcessDownloadRequestPayload(pld)
		tBobNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse)
	})

	t.Run("Ask bob for the blocks that he doesn't have", func(t *testing.T) {
		pld := &payload.DownloadRequestPayload{
			SessionID: 1,
			Initiator: tAnotherPeerID,
			Target:    tBobPeerID,
			From:      1000,
			To:        1010,
		}
		tBobSync.stateSync.ProcessDownloadRequestPayload(pld)
		msg := tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // No more block
		assert.Equal(t, msg.Payload.(*payload.DownloadResponsePayload).ResponseCode, payload.ResponseCodeNoMoreBlocks)
	})

}

func TestDownloadBlocks(t *testing.T) {
	setup(t)

	disableHeartbeat(t)

	// Clear alice store
	tAliceSync.cache.Clear()
	tAliceState.Store.Blocks = make(map[int]*block.Block)
	tAliceConsensus.Started = false

	joinBobToTheSetAndUpdateHRS(t)
	addMoreBlocksForBobAndSendBlockAnnounceMessage(t, 80)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeBlockAnnounce)

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 1-10
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 11-20
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 21-31 (one extra block)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // NoMoreBlock

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 40-49
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 50-59
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 60-70 (one extra block)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // NoMoreBlock

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 61-70
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 71-80
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 81-91 (one extra block)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // NoMoreBlock

	// Latest block requests
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // 91-100
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // 101-101
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // Synced

	assert.True(t, tBobConsensus.Started)
	assert.False(t, tAliceSync.peerSet.HasAnyValidSession())
	assert.False(t, tBobSync.peerSet.HasAnyValidSession())
}

func TestDownloadBlockBobIsNotInSet(t *testing.T) {
	setup(t)

	// Clear alice store
	tAliceSync.cache.Clear()
	tAliceState.Store.Blocks = make(map[int]*block.Block)
	tAliceConsensus.Started = false

	addMoreBlocksForBobAndUpdateHRS(t, 80)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeHeartBeat)

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 1-10
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 11-20
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 21-31 (one extra block)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // NoMoreBlock

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 40-49
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 50-59
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 60-70 (one extra block)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // NoMoreBlock

	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadRequest)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 61-70
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 71-80
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // 81-91 (one extra block)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeDownloadResponse) // NoMoreBlock

	// Latest block requests
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // 91-100
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // 101-101
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksResponse) // Synced

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

func TestOneBlockBehind(t *testing.T) {
	setup(t)

	addMoreBlocksForBobAndUpdateHRS(t, 1)

	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeHeartBeat)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)

}
