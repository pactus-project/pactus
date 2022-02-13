package sync

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func TestSessionTimeout(t *testing.T) {
	tBobConfig.SessionTimeout = 200 * time.Millisecond
	setup(t)

	t.Run("An unknown peers claims has more blocks. Alice requests for more blocks. Alice doesn't get any response. Session should be closed", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		pid := util.RandomPeerID()
		pld := payload.NewAleykPayload(pid, "Oscar", 6666, 0x1, tAlicePeerID, payload.ResponseCodeOK, "ok") // InitialBlockDownload: true
		signer.SignMsg(pld)
		simulatingReceiveingNewMessage(t, tAliceSync, pld, pid)

		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)

		assert.True(t, tAliceSync.peerSet.HasAnyOpenSession())
		time.Sleep(2 * tAliceConfig.SessionTimeout)
		assert.False(t, tAliceSync.peerSet.HasAnyOpenSession())
	})
}

func TestLatestBlocksRequestMessages(t *testing.T) {
	tBobConfig.InitialBlockDownload = false

	setup(t)
	disableHeartbeat(t)

	addMoreBlocksForBob(t, 12)
	sid := util.RandInt(100)

	t.Run("Bob should reject requests from unknown peer", func(t *testing.T) {
		pld := payload.NewBlocksRequestPayload(sid, 100, 105)
		initiator := util.RandomPeerID()
		simulatingReceiveingNewMessage(t, tBobSync, pld, initiator)

		msg := shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse)
		assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeRejected)
		assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).From, 0)
	})

	t.Run("Bob should reject requests with invalid ranges", func(t *testing.T) {
		pld := payload.NewBlocksRequestPayload(sid, 0, 20)
		simulatingReceiveingNewMessage(t, tBobSync, pld, tAlicePeerID)

		msg := shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse)
		assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeRejected)
		assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).From, 0)
	})

	t.Run("Bob didn't set the InitialBlockDownload flag", func(t *testing.T) {
		pid := util.RandomPeerID()
		heightBob := tBobState.LastBlockHeight()

		t.Run("Bob handshakes with the new peer", func(t *testing.T) {
			signer := bls.GenerateTestSigner()
			pld := payload.NewSalamPayload(pid, "new-peer", 0, 0, tBobState.GenHash)
			signer.SignMsg(pld)
			simulatingReceiveingNewMessage(t, tBobSync, pld, pid)

			shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeAleyk)
		})

		t.Run("Bob rejects request with more than `LatestBlockInterval`", func(t *testing.T) {
			pld := payload.NewBlocksRequestPayload(sid, 0, heightBob)
			simulatingReceiveingNewMessage(t, tBobSync, pld, pid)

			msg := shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse)
			assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeRejected)
			assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).From, 0)
		})

		t.Run("Bob accepts request within `LatestBlockInterval`", func(t *testing.T) {
			pld := payload.NewBlocksRequestPayload(sid, heightBob-5, heightBob)
			simulatingReceiveingNewMessage(t, tBobSync, pld, pid)

			msg1 := shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse)
			assert.Equal(t, msg1.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeMoreBlocks)
			assert.Equal(t, msg1.Payload.(*payload.BlocksResponsePayload).From, heightBob-5)
			assert.Equal(t, msg1.Payload.(*payload.BlocksResponsePayload).To(), heightBob)

			msg2 := shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse)
			assert.Equal(t, msg2.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeSynced)
			assert.Equal(t, msg2.Payload.(*payload.BlocksResponsePayload).From, heightBob)
			assert.Equal(t, msg2.Payload.(*payload.BlocksResponsePayload).To(), heightBob)
		})

		t.Run("Peer requests from Bob to send the blocks again, Bob should reject it.", func(t *testing.T) {
			pld := payload.NewBlocksRequestPayload(sid, heightBob-5, heightBob)
			simulatingReceiveingNewMessage(t, tBobSync, pld, pid)

			msg := shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse)
			assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeRejected)
		})

		t.Run("Bob doesn't have request blocks", func(t *testing.T) {
			pld := payload.NewBlocksRequestPayload(sid, 100, 105)
			simulatingReceiveingNewMessage(t, tBobSync, pld, pid)

			msg := shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse)
			assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeSynced)
		})
	})

	t.Run("Bob is busy", func(t *testing.T) {
		tBobSync.peerSet.OpenSession(util.RandomPeerID())
		tBobSync.peerSet.OpenSession(util.RandomPeerID())
		tBobSync.peerSet.OpenSession(util.RandomPeerID())
		tBobSync.peerSet.OpenSession(util.RandomPeerID())
		tBobSync.peerSet.OpenSession(util.RandomPeerID())
		require.Equal(t, tBobSync.peerSet.NumberOfOpenSessions(), 5)

		s := tAliceSync.peerSet.OpenSession(tBobPeerID)
		pld := payload.NewBlocksRequestPayload(s.SessionID(), 100, 105)
		simulatingReceiveingNewMessage(t, tBobSync, pld, tAlicePeerID)
		assert.Equal(t, tAliceSync.peerSet.NumberOfOpenSessions(), 2)

		msg := shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlocksResponse)
		assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeBusy)
		assert.Equal(t, tAliceSync.peerSet.NumberOfOpenSessions(), 1, "Alice should close the session")
	})
}
