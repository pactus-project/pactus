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
	tConfig.SessionTimeout = 200 * time.Millisecond
	setup(t)

	t.Run("An unknown peers claims to have more blocks. Session should be closed after timeout", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		pid := util.RandomPeerID()
		pld := payload.NewHelloPayload(pid, "Oscar", 6666, payload.FlagInitialBlockDownload, tState.GenHash)
		signer.SignMsg(pld)

		assert.NoError(t, testReceiveingNewMessage(tSync, pld, pid))

		shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlocksRequest)

		assert.True(t, tSync.peerSet.HasAnyOpenSession())
		time.Sleep(2 * tConfig.SessionTimeout)
		assert.False(t, tSync.peerSet.HasAnyOpenSession())
	})
}

func TestLatestBlocksRequestMessages(t *testing.T) {
	tConfig.InitialBlockDownload = false
	setup(t)

	sid := util.RandInt(100)
	pid := util.RandomPeerID()

	t.Run("Reject request from unknown peers", func(t *testing.T) {
		pld := payload.NewBlocksRequestPayload(sid, 100, 105)
		assert.Error(t, testReceiveingNewMessage(tSync, pld, pid))

		msg := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlocksResponse)
		assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeRejected)
		assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).From, 0)
	})

	pub, _ := bls.GenerateTestKeyPair()
	testAddPeer(t, pub, pid)

	t.Run("Reject request with invalid range", func(t *testing.T) {
		pld := payload.NewBlocksRequestPayload(sid, 0, 20)
		assert.Error(t, testReceiveingNewMessage(tSync, pld, pid))

		msg := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlocksResponse)
		assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeRejected)
		assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).From, 0)
	})

	t.Run("InitialBlockDownload flag is not set", func(t *testing.T) {
		heightBob := tState.LastBlockHeight()

		t.Run("Reject requests with more than `LatestBlockInterval`", func(t *testing.T) {
			pld := payload.NewBlocksRequestPayload(sid, 0, heightBob)
			assert.Error(t, testReceiveingNewMessage(tSync, pld, pid))

			msg := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlocksResponse)
			assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeRejected)
			assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).From, 0)
		})

		t.Run("Accept request within `LatestBlockInterval`", func(t *testing.T) {
			pld := payload.NewBlocksRequestPayload(sid, heightBob-5, heightBob)
			assert.NoError(t, testReceiveingNewMessage(tSync, pld, pid))

			msg1 := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlocksResponse)
			assert.Equal(t, msg1.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeMoreBlocks)
			assert.Equal(t, msg1.Payload.(*payload.BlocksResponsePayload).From, heightBob-5)
			assert.Equal(t, msg1.Payload.(*payload.BlocksResponsePayload).To(), heightBob)

			msg2 := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlocksResponse)
			assert.Equal(t, msg2.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeSynced)
			assert.Equal(t, msg2.Payload.(*payload.BlocksResponsePayload).From, heightBob)
			assert.Equal(t, msg2.Payload.(*payload.BlocksResponsePayload).To(), heightBob)
		})

		t.Run("Peer requests to send the blocks again, It should be rejected", func(t *testing.T) {
			pld := payload.NewBlocksRequestPayload(sid, heightBob-1, heightBob)
			assert.Error(t, testReceiveingNewMessage(tSync, pld, pid))

			msg := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlocksResponse)
			assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeRejected)
		})

		t.Run("Peer doesn't have requested blocks", func(t *testing.T) {
			pld := payload.NewBlocksRequestPayload(sid, 100, 105)
			assert.NoError(t, testReceiveingNewMessage(tSync, pld, pid))

			msg := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlocksResponse)
			assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeSynced)
		})
	})

	t.Run("Peer is busy", func(t *testing.T) {
		tSync.peerSet.OpenSession(util.RandomPeerID())
		tSync.peerSet.OpenSession(util.RandomPeerID())
		tSync.peerSet.OpenSession(util.RandomPeerID())
		tSync.peerSet.OpenSession(util.RandomPeerID())
		tSync.peerSet.OpenSession(util.RandomPeerID())
		require.Equal(t, tSync.peerSet.NumberOfOpenSessions(), 5)

		s := tSync.peerSet.OpenSession(tNetwork.SelfID())
		pld := payload.NewBlocksRequestPayload(s.SessionID(), 100, 105)
		assert.NoError(t, testReceiveingNewMessage(tSync, pld, pid))
		msg := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlocksResponse)
		assert.Equal(t, msg.Payload.(*payload.BlocksResponsePayload).ResponseCode, payload.ResponseCodeBusy)
	})
}
