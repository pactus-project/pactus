package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
)

func TestOneBlockShorter(t *testing.T) {
	setup(t)
	disableHeartbeat(t)

	t.Run("Bob commits two blocks. Alice should request for the lastest block.", func(t *testing.T) {
		joinBobToCommittee(t)
		addMoreBlocksForBobAndAnnounceLastBlock(t, 2)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlockAnnounce)

		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // 22
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // No more block

		assert.Equal(t, tAliceState.LastBlockHash(), tBobState.LastBlockHash())
		assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
	})
}

func TestDownloadBlocks(t *testing.T) {
	LatestBlockInterval = 30
	//
	tBobConfig.InitialBlockDownload = true

	setup(t)
	disableHeartbeat(t)

	// Reset state for Alice
	tAliceSync.cache.Clear()
	tAliceState.Store.Blocks = make(map[int]*block.Block)
	tBobSync.peerSet.RemovePeer(tAlicePeerID)
	p := tBobSync.peerSet.MustGetPeer(tAlicePeerID)
	p.UpdateStatus(peerset.StatusCodeOK)

	joinBobToCommittee(t)
	addMoreBlocksForBobAndAnnounceLastBlock(t, 79) // total blocks: 21+79 = 100
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlockAnnounce)

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // 1-10
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // 11-20
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // 21-30
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // NoMoreBlock

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // 31-40
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // 41-50
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // 51-60
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // NoMoreBlock

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // 61-70
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // 71-80
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // 81-90
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // NoMoreBlock

	// Latest block requests
	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // 91-100
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse) // Synced

	assert.Equal(t, tAliceState.LastBlockHash(), tBobState.LastBlockHash())
	assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
	assert.False(t, tAliceSync.peerSet.HasAnyOpenSession())
	assert.False(t, tBobSync.peerSet.HasAnyOpenSession())
}
