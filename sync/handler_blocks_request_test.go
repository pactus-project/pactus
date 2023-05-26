package sync

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionTimeout(t *testing.T) {
	tConfig.SessionTimeout = 200 * time.Millisecond
	setup(t)

	t.Run("An unknown peers claims to have more blocks. Session should be closed after timeout", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		pid := network.TestRandomPeerID()
		claimedHeight := uint32(6666)
		msg := message.NewHelloMessage(pid, "Oscar", claimedHeight, message.FlagNodeNetwork, tState.Genesis().Hash())
		signer.SignMsg(msg)

		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksRequest)

		assert.True(t, tSync.peerSet.HasAnyOpenSession())
		time.Sleep(2 * tConfig.SessionTimeout)
		assert.False(t, tSync.peerSet.HasAnyOpenSession())

		peer := tSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Height, claimedHeight)
		assert.Equal(t, tSync.peerSet.MaxClaimedHeight(), claimedHeight)
		// TODO: This is not really good that a bad peer can manipulate the MaxCalim height
		// Here is a possible solution:
		// 1- A peer claims that he has more blocks
		// 2- We ask him a block_request message
		// 3- He doesn't respond
		// 4- We close the session, marking the peer as a bad peer.
		// 5- If MaxClaimedHeight is same as peer height, we set it to zero
	})
}

func TestLatestBlocksRequestMessages(t *testing.T) {
	tConfig.NodeNetwork = false
	setup(t)
	testAddBlocks(t, tState, 10)

	t.Run("NodeNetwork flag is not set", func(t *testing.T) {
		curHeight := tState.LastBlockHeight()

		sid := int(util.RandInt32(100))
		pid := network.TestRandomPeerID()

		t.Run("Reject request from unknown peers", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, curHeight-1, 1)
			assert.Error(t, testReceivingNewMessage(tSync, msg, pid))

			bdl := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).From, uint32(0))
		})

		pub, _ := bls.GenerateTestKeyPair()
		testAddPeer(t, pub, pid, false)

		t.Run("Reject requests not within `LatestBlockInterval`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 1, 2)
			assert.Error(t, testReceivingNewMessage(tSync, msg, pid))

			bdl := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).From, uint32(0))
		})

		t.Run("Request blocks more than `LatestBlockInterval`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 0, LatestBlockInterval+1)
			assert.Error(t, testReceivingNewMessage(tSync, msg, pid))

			bdl := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
		})

		t.Run("Accept request within `LatestBlockInterval`", func(t *testing.T) {
			t.Run("Peer needs more block", func(t *testing.T) {
				msg := message.NewBlocksRequestMessage(sid, curHeight-tConfig.BlockPerMessage, tConfig.BlockPerMessage)
				assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

				bdl1 := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).From, curHeight-tConfig.BlockPerMessage)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).To(), curHeight-1)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).Count(), tConfig.BlockPerMessage)
				assert.Zero(t, bdl1.Message.(*message.BlocksResponseMessage).LastCertificateHeight())

				bdl2 := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
				assert.Equal(t, bdl2.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeNoMoreBlocks)
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).From)
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).To())
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).Count())
				assert.Zero(t, bdl1.Message.(*message.BlocksResponseMessage).LastCertificateHeight())

				peer := tSync.peerSet.GetPeer(pid)
				assert.Equal(t, peer.Height, curHeight-1) // Peer needs one more block
			})

			t.Run("Peer synced", func(t *testing.T) {
				msg := message.NewBlocksRequestMessage(sid, curHeight-tConfig.BlockPerMessage+1, tConfig.BlockPerMessage)
				assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

				bdl1 := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).From, curHeight-tConfig.BlockPerMessage+1)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).To(), curHeight)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).Count(), tConfig.BlockPerMessage)
				assert.Zero(t, bdl1.Message.(*message.BlocksResponseMessage).LastCertificateHeight())

				bdl2 := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
				assert.Equal(t, bdl2.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeSynced)
				assert.Equal(t, bdl2.Message.(*message.BlocksResponseMessage).From, curHeight)
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).To())
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).Count())
				assert.Equal(t, bdl2.Message.(*message.BlocksResponseMessage).LastCertificateHeight(), curHeight)

				peer := tSync.peerSet.GetPeer(pid)
				assert.Equal(t, peer.Height, curHeight) // Peer is synced
			})
		})

		// t.Run("Peer requests to send the blocks again, It should be rejected", func(t *testing.T) {
		// 	msg := message.NewBlocksRequestMessage(sid, curHeight-1, 1)
		// 	assert.Error(t, testReceivingNewMessage(tSync, msg, pid))

		// 	bdl := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
		// 	assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
		// 	assert.Zero(t, bdl.Message.(*message.BlocksResponseMessage).LastCertificateHeight())

		// })

		t.Run("Peer requests blocks that we don't have", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, curHeight+100, 1)
			assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

			bdl := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeSynced)

			peer := tSync.peerSet.GetPeer(pid)
			assert.Equal(t, peer.Height, curHeight+99)
		})
	})

	t.Run("NodeNetwork flag set", func(t *testing.T) {
		tSync.config.NodeNetwork = true

		sid := int(util.RandInt32(100))
		pid := network.TestRandomPeerID()
		pub, _ := bls.GenerateTestKeyPair()
		testAddPeer(t, pub, pid, false)

		t.Run("Requesting one block", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 1, 2)
			assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

			msg1 := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
			assert.Equal(t, msg1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)

			msg2 := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
			assert.Equal(t, msg2.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeNoMoreBlocks)

			peer := tSync.peerSet.GetPeer(pid)
			assert.Equal(t, peer.Height, uint32(2))
		})

		t.Run("Peer is busy", func(t *testing.T) {
			tSync.peerSet.OpenSession(network.TestRandomPeerID())
			tSync.peerSet.OpenSession(network.TestRandomPeerID())
			tSync.peerSet.OpenSession(network.TestRandomPeerID())
			tSync.peerSet.OpenSession(network.TestRandomPeerID())
			tSync.peerSet.OpenSession(network.TestRandomPeerID())
			require.Equal(t, tSync.peerSet.NumberOfOpenSessions(), 5)

			s := tSync.peerSet.OpenSession(tNetwork.SelfID())
			msg := message.NewBlocksRequestMessage(s.SessionID(), 100, 105)
			assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))
			bdl := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeBusy)
		})
	})
}
