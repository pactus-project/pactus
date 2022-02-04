package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
)

func TestDownloadBlocks(t *testing.T) {
	LatestBlockInterval = 30
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

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeDownloadRequest)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadResponse) // 1-10
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadResponse) // 11-20
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadResponse) // 21-31 (one extra block)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadResponse) // NoMoreBlock

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeDownloadRequest)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadResponse) // 40-49
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadResponse) // 50-59
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadResponse) // 60-70 (one extra block)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadResponse) // NoMoreBlock

	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeDownloadRequest)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadResponse) // 61-70
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadResponse) // 71-80
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadResponse) // 81-91 (one extra block)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadResponse) // NoMoreBlock

	// Latest block requests
	shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeLatestBlocksRequest)
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeLatestBlocksResponse) // 91-100
	shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeLatestBlocksResponse) // Synced

	assert.Equal(t, tAliceState.LastBlockHash(), tBobState.LastBlockHash())
	assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
	assert.False(t, tAliceSync.peerSet.HasAnyOpenSession())
	assert.False(t, tBobSync.peerSet.HasAnyOpenSession())
}
