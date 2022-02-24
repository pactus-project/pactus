package sync

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/util"
)

func TestSessionTimeout(t *testing.T) {
	tConfig.SessionTimeout = 200 * time.Millisecond
	setup(t)

	t.Run("An unknown peers claims to have more blocks. Session should be closed after timeout", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		pid := util.RandomPeerID()
		msg := message.NewHelloMessage(pid, "Oscar", 6666, message.FlagInitialBlockDownload, tState.GenHash)
		signer.SignMsg(msg)

		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksRequest)

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
		msg := message.NewBlocksRequestMessage(sid, 100, 105)
		assert.Error(t, testReceiveingNewMessage(tSync, msg, pid))

		bnd := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
		assert.Equal(t, bnd.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
		assert.Equal(t, bnd.Message.(*message.BlocksResponseMessage).From, 0)
	})

	pub, _ := bls.GenerateTestKeyPair()
	testAddPeer(t, pub, pid)

	t.Run("Reject request with invalid range", func(t *testing.T) {
		msg := message.NewBlocksRequestMessage(sid, 0, 20)
		assert.Error(t, testReceiveingNewMessage(tSync, msg, pid))

		bnd := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
		assert.Equal(t, bnd.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
		assert.Equal(t, bnd.Message.(*message.BlocksResponseMessage).From, 0)
	})

	t.Run("InitialBlockDownload flag is not set", func(t *testing.T) {
		heightBob := tState.LastBlockHeight()

		t.Run("Reject requests with more than `LatestBlockInterval`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 0, heightBob)
			assert.Error(t, testReceiveingNewMessage(tSync, msg, pid))

			bnd := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
			assert.Equal(t, bnd.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
			assert.Equal(t, bnd.Message.(*message.BlocksResponseMessage).From, 0)
		})

		t.Run("Accept request within `LatestBlockInterval`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, heightBob-5, heightBob)
			assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

			msg1 := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
			assert.Equal(t, msg1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)
			assert.Equal(t, msg1.Message.(*message.BlocksResponseMessage).From, heightBob-5)
			assert.Equal(t, msg1.Message.(*message.BlocksResponseMessage).To(), heightBob)

			msg2 := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
			assert.Equal(t, msg2.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeSynced)
			assert.Equal(t, msg2.Message.(*message.BlocksResponseMessage).From, heightBob)
			assert.Equal(t, msg2.Message.(*message.BlocksResponseMessage).To(), heightBob)
		})

		t.Run("Peer requests to send the blocks again, It should be rejected", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, heightBob-1, heightBob)
			assert.Error(t, testReceiveingNewMessage(tSync, msg, pid))

			bnd := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
			assert.Equal(t, bnd.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
		})

		t.Run("Peer doesn't have requested blocks", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 100, 105)
			assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))

			bnd := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
			assert.Equal(t, bnd.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeSynced)
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
		msg := message.NewBlocksRequestMessage(s.SessionID(), 100, 105)
		assert.NoError(t, testReceiveingNewMessage(tSync, msg, pid))
		bnd := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
		assert.Equal(t, bnd.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeBusy)
	})
}
