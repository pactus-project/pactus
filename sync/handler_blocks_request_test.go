package sync

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionTimeout(t *testing.T) {
	config := testConfig()
	config.SessionTimeout = 200 * time.Millisecond
	td := setup(t, config)

	t.Run("An unknown peers claims to have more blocks. Session should be closed after timeout", func(t *testing.T) {
		signer := td.RandomSigner()
		pid := td.RandomPeerID()
		claimedHeight := uint32(6666)
		msg := message.NewHelloMessage(pid, "Oscar", claimedHeight, message.FlagNodeNetwork,
			td.state.LastBlockHash(), td.state.Genesis().Hash())
		signer.SignMsg(msg)

		assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

		td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeBlocksRequest)

		assert.True(t, td.sync.peerSet.HasAnyOpenSession())
		time.Sleep(2 * config.SessionTimeout)
		assert.False(t, td.sync.peerSet.HasAnyOpenSession())

		peer := td.sync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Height, claimedHeight)
		assert.Equal(t, td.sync.peerSet.MaxClaimedHeight(), claimedHeight)
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
	config := testConfig()
	config.NodeNetwork = false
	td := setup(t, config)
	td.addBlocks(t, td.state, 10)

	t.Run("NodeNetwork flag is not set", func(t *testing.T) {
		curHeight := td.state.LastBlockHeight()

		sid := td.RandInt(100)
		pid := td.RandomPeerID()

		t.Run("Reject request from unknown peers", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, curHeight-1, 1)
			assert.Error(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).From, uint32(0))
		})

		pub, _ := td.RandomBLSKeyPair()
		td.addPeer(t, pub, pid, false)

		t.Run("Reject requests not within `LatestBlockInterval`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 1, 2)
			assert.Error(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).From, uint32(0))
		})

		t.Run("Request blocks more than `LatestBlockInterval`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 0, LatestBlockInterval+1)
			assert.Error(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
		})

		t.Run("Accept request within `LatestBlockInterval`", func(t *testing.T) {
			t.Run("Peer needs more block", func(t *testing.T) {
				msg := message.NewBlocksRequestMessage(sid, curHeight-config.BlockPerMessage, config.BlockPerMessage)
				assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

				bdl1 := td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeBlocksResponse)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).From, curHeight-config.BlockPerMessage)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).To(), curHeight-1)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).Count(), config.BlockPerMessage)
				assert.Zero(t, bdl1.Message.(*message.BlocksResponseMessage).LastCertificateHeight())

				bdl2 := td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeBlocksResponse)
				assert.Equal(t, bdl2.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeNoMoreBlocks)
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).From)
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).To())
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).Count())
				assert.Zero(t, bdl1.Message.(*message.BlocksResponseMessage).LastCertificateHeight())
			})

			t.Run("Peer synced", func(t *testing.T) {
				msg := message.NewBlocksRequestMessage(sid, curHeight-config.BlockPerMessage+1, config.BlockPerMessage)
				assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

				bdl1 := td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeBlocksResponse)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).From, curHeight-config.BlockPerMessage+1)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).To(), curHeight)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).Count(), config.BlockPerMessage)
				assert.Zero(t, bdl1.Message.(*message.BlocksResponseMessage).LastCertificateHeight())

				bdl2 := td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeBlocksResponse)
				assert.Equal(t, bdl2.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeSynced)
				assert.Equal(t, bdl2.Message.(*message.BlocksResponseMessage).From, curHeight)
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).To())
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).Count())
				assert.Equal(t, bdl2.Message.(*message.BlocksResponseMessage).LastCertificateHeight(), curHeight)
			})
		})

		t.Run("Peer requests blocks that we don't have", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, curHeight+100, 1)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeSynced)
		})
	})

	t.Run("NodeNetwork flag set", func(t *testing.T) {
		td.sync.config.NodeNetwork = true

		sid := td.RandInt(100)
		pid := td.RandomPeerID()
		pub, _ := td.RandomBLSKeyPair()
		td.addPeer(t, pub, pid, false)

		t.Run("Requesting one block", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 1, 2)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			msg1 := td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeBlocksResponse)
			assert.Equal(t, msg1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)

			msg2 := td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeBlocksResponse)
			assert.Equal(t, msg2.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeNoMoreBlocks)
		})

		t.Run("Peer is busy", func(t *testing.T) {
			td.sync.peerSet.OpenSession(td.RandomPeerID())
			td.sync.peerSet.OpenSession(td.RandomPeerID())
			td.sync.peerSet.OpenSession(td.RandomPeerID())
			td.sync.peerSet.OpenSession(td.RandomPeerID())
			td.sync.peerSet.OpenSession(td.RandomPeerID())
			require.Equal(t, td.sync.peerSet.NumberOfOpenSessions(), 5)

			s := td.sync.peerSet.OpenSession(td.network.SelfID())
			msg := message.NewBlocksRequestMessage(s.SessionID(), 100, 105)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))
			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.MessageTypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeBusy)
		})
	})
}
