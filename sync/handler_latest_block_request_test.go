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

	t.Run("An unknown peers claims has more blocks. Alice requests for more blocks. Alice doesn't get any response. Session should be closed", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		pld := payload.NewAleykPayload(tAlicePeerID, payload.ResponseCodeOK, "ok", "devil", pub, 6666, 0x1) // InitialBlockDownload:  true
		tAliceNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeDownloadRequest)

		assert.True(t, tAliceSync.peerSet.HasAnyOpenSession())
		time.Sleep(2 * tAliceConfig.SessionTimeout)
		assert.False(t, tAliceSync.peerSet.HasAnyOpenSession())
	})
}

func TestLatestBlocksRequestMessages(t *testing.T) {
	setup(t)
	disableHeartbeat(t)

	addMoreBlocksForBob(t, 12)
	pid := util.RandomPeerID()

	t.Run("Bob received request from unknown peer. Request should be rejected", func(t *testing.T) {
		pld := payload.NewLatestBlocksRequestPayload(6, tBobPeerID, 100, 105)
		tBobNet.ReceivingMessageFromOtherPeer(pid, pld)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tBobNet, payload.PayloadTypeLatestBlocksResponse, payload.ResponseCodeRejected)

		t.Run("Bob handshakes with the new peer", func(t *testing.T) {
			_, pub, _ := crypto.GenerateTestKeyPair()
			pld := payload.NewSalamPayload("new-peer", pub, tBobState.GenHash, 0, 0)
			tBobNet.ReceivingMessageFromOtherPeer(pid, pld)

			shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeAleyk)
			// First session opens here
		})
	})

	t.Run("Bob should not pay attention to requests for other peers", func(t *testing.T) {
		pld := payload.NewLatestBlocksRequestPayload(6, util.RandomPeerID(), 0, 20)
		tBobNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldNotPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeLatestBlocksResponse)
	})

	t.Run("Bob should reject requests with invalid ranges", func(t *testing.T) {
		pld := payload.NewLatestBlocksRequestPayload(6, tBobPeerID, 0, 20)
		tBobNet.ReceivingMessageFromOtherPeer(pid, pld)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tBobNet, payload.PayloadTypeLatestBlocksResponse, payload.ResponseCodeRejected)
	})

	t.Run("Bob should send blocks to peer", func(t *testing.T) {
		bobHeight := tBobState.LastBlockHeight()
		pld := payload.NewLatestBlocksRequestPayload(6, tBobPeerID, bobHeight-5, bobHeight)
		tBobNet.ReceivingMessageFromOtherPeer(pid, pld)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tBobNet, payload.PayloadTypeLatestBlocksResponse, payload.ResponseCodeMoreBlocks)
		shouldPublishPayloadWithThisTypeAndResponseCode(t, tBobNet, payload.PayloadTypeLatestBlocksResponse, payload.ResponseCodeSynced)

		t.Run("Peer requests from Bob to send the blocks again, Bob should reject it.", func(t *testing.T) {
			tBobNet.ReceivingMessageFromOtherPeer(pid, pld)

			shouldPublishPayloadWithThisTypeAndResponseCode(t, tBobNet, payload.PayloadTypeLatestBlocksResponse, payload.ResponseCodeRejected)
		})
	})

	t.Run("Bob doesn't have request blocks", func(t *testing.T) {
		pld := payload.NewLatestBlocksRequestPayload(6, tBobPeerID, 100, 105)
		tBobNet.ReceivingMessageFromOtherPeer(pid, pld)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tBobNet, payload.PayloadTypeLatestBlocksResponse, payload.ResponseCodeSynced)
	})

	t.Run("Bob is busy", func(t *testing.T) {
		tBobSync.peerSet.OpenSession(util.RandomPeerID())
		tBobSync.peerSet.OpenSession(util.RandomPeerID())
		tBobSync.peerSet.OpenSession(util.RandomPeerID())
		tBobSync.peerSet.OpenSession(util.RandomPeerID())
		tBobSync.peerSet.OpenSession(util.RandomPeerID())
		assert.Equal(t, tBobSync.peerSet.NumberOfOpenSessions(), 5)

		s := tAliceSync.peerSet.OpenSession(tBobPeerID)
		pld := payload.NewLatestBlocksRequestPayload(s.SessionID(), tBobPeerID, 100, 105)
		tBobNet.ReceivingMessageFromOtherPeer(tAlicePeerID, pld)
		assert.Equal(t, tAliceSync.peerSet.NumberOfOpenSessions(), 2)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tBobNet, payload.PayloadTypeLatestBlocksResponse, payload.ResponseCodeBusy)
		assert.Equal(t, tAliceSync.peerSet.NumberOfOpenSessions(), 1, "Alice should close the session")
	})
}
