package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func TestDownloadBlocksRequestMessages(t *testing.T) {
	setup(t)
	disableHeartbeat(t)

	pid := util.RandomPeerID()

	t.Run("Alice received request from unknown peer. Request should be rejected", func(t *testing.T) {
		pld := payload.NewDownloadRequestPayload(6, tAlicePeerID, 100, 105)
		tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tAliceNet, payload.PayloadTypeDownloadResponse, payload.ResponseCodeRejected)

		t.Run("Alice handshakes with the new peer", func(t *testing.T) {
			_, pub, _ := bls.GenerateTestKeyPair()
			pld := payload.NewSalamPayload("new-peer", pub, tAliceState.GenHash, 0, 0)
			tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

			shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeAleyk)
		})
	})

	t.Run("Alice should not pay attention to requests for other peers", func(t *testing.T) {
		pld := payload.NewDownloadRequestPayload(6, util.RandomPeerID(), 0, 10)
		tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeDownloadResponse)
	})

	t.Run("Alice should not pay attention to requests with invalid ranges", func(t *testing.T) {
		pld := payload.NewDownloadRequestPayload(6, tAlicePeerID, 1000, 2000)
		tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tAliceNet, payload.PayloadTypeDownloadResponse, payload.ResponseCodeRejected)
	})

	t.Run("Alice send blocks to the new peer", func(t *testing.T) {
		pld := payload.NewDownloadRequestPayload(6, tAlicePeerID, 1, 10)
		tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tAliceNet, payload.PayloadTypeDownloadResponse, payload.ResponseCodeMoreBlocks)
		shouldPublishPayloadWithThisTypeAndResponseCode(t, tAliceNet, payload.PayloadTypeDownloadResponse, payload.ResponseCodeNoMoreBlocks)

		t.Run("Peer requests from Alice to send the blocks again, Alice should reject it.", func(t *testing.T) {
			tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

			shouldPublishPayloadWithThisTypeAndResponseCode(t, tAliceNet, payload.PayloadTypeDownloadResponse, payload.ResponseCodeRejected)
		})
	})

	t.Run("Alice doesn't have request blocks", func(t *testing.T) {
		pld := payload.NewDownloadRequestPayload(6, tAlicePeerID, 100, 105)
		tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tAliceNet, payload.PayloadTypeDownloadResponse, payload.ResponseCodeNoMoreBlocks)
	})

	t.Run("Alice is busy", func(t *testing.T) {
		tAliceSync.peerSet.OpenSession(util.RandomPeerID())
		tAliceSync.peerSet.OpenSession(util.RandomPeerID())
		tAliceSync.peerSet.OpenSession(util.RandomPeerID())
		tAliceSync.peerSet.OpenSession(util.RandomPeerID())
		tAliceSync.peerSet.OpenSession(util.RandomPeerID())
		assert.Equal(t, tAliceSync.peerSet.NumberOfOpenSessions(), 5)

		s := tBobSync.peerSet.OpenSession(tAlicePeerID)
		pld := payload.NewDownloadRequestPayload(s.SessionID(), tAlicePeerID, 100, 105)
		tAliceNet.ReceivingMessageFromOtherPeer(tBobPeerID, pld)
		assert.Equal(t, tBobSync.peerSet.NumberOfOpenSessions(), 1)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tAliceNet, payload.PayloadTypeDownloadResponse, payload.ResponseCodeBusy)
		assert.Equal(t, tBobSync.peerSet.NumberOfOpenSessions(), 0, "Bob should close the session")
	})
}
