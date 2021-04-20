package sync

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func TestSessionTimeout(t *testing.T) {
	tBobConfig.SessionTimeout = 200 * time.Millisecond
	setup(t)

	t.Run("An unknown peers claims has more blocks. Bob requests for more blocks. Bob doesn't get any response. Session should be closed", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		pld := payload.NewSalamPayload("devil", pub, tBobState.GenHash, 6666, 0x1) // InitialBlockDownload:  true
		tBobNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeDownloadRequest)

		assert.True(t, tBobSync.peerSet.HasAnyValidSession())
		time.Sleep(2 * tBobConfig.SessionTimeout)
		assert.False(t, tBobSync.peerSet.HasAnyValidSession())
	})
}

func TestLatestBlocksRequestMessages(t *testing.T) {
	setup(t)
	disableHeartbeat(t)

	addMoreBlocksForBob(t, 10)

	t.Run("Bob should not pay attention to requests for other peers", func(t *testing.T) {
		pld := payload.NewLatestBlocksRequestPayload(6, util.RandomPeerID(), 10, 20)
		tBobNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldNotPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeLatestBlocksResponse)
	})

	t.Run("Bob should not pay attention to requests with invalid ranges", func(t *testing.T) {
		pld := payload.NewLatestBlocksRequestPayload(6, tBobPeerID, 0, 20)
		tBobNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldNotPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeLatestBlocksResponse)
	})

	t.Run("Bob doesn't have request blocks", func(t *testing.T) {
		pld := payload.NewLatestBlocksRequestPayload(6, tBobPeerID, 100, 105)
		tBobNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tBobNet, payload.PayloadTypeLatestBlocksResponse, payload.ResponseCodeSynced)
	})
}
