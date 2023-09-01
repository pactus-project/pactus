package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: check error types and clean the test.
func TestLatestBlocksRequestMessages(t *testing.T) {
	config := testConfig()
	config.NodeNetwork = false
	td := setup(t, config)
	td.addBlocks(t, td.state, 10)

	t.Run("NodeNetwork flag is not set", func(t *testing.T) {
		curHeight := td.state.LastBlockHeight()

		sid := td.RandInt(100)
		pid := td.RandPeerID()

		t.Run("Reject request from unknown peers", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, curHeight-1, 1)
			assert.Error(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).From, uint32(0))
		})

		pub, _ := td.RandBLSKeyPair()
		td.addPeer(t, pub, pid, false)

		t.Run("Reject requests not within `LatestBlockInterval`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 1, 2)
			assert.Error(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).From, uint32(0))
		})

		t.Run("Request blocks more than `LatestBlockInterval`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 10, LatestBlockInterval+1)
			assert.Error(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
		})

		t.Run("Accept request within `LatestBlockInterval`", func(t *testing.T) {
			t.Run("Peer needs more block", func(t *testing.T) {
				msg := message.NewBlocksRequestMessage(sid, curHeight-config.BlockPerMessage, config.BlockPerMessage)
				assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

				bdl1 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).From, curHeight-config.BlockPerMessage)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).To(), curHeight-1)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).Count(), config.BlockPerMessage)
				assert.Zero(t, bdl1.Message.(*message.BlocksResponseMessage).LastCertificateHeight())

				bdl2 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
				assert.Equal(t, bdl2.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeNoMoreBlocks)
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).From)
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).To())
				assert.Zero(t, bdl2.Message.(*message.BlocksResponseMessage).Count())
				assert.Zero(t, bdl1.Message.(*message.BlocksResponseMessage).LastCertificateHeight())
			})

			t.Run("Peer synced", func(t *testing.T) {
				msg := message.NewBlocksRequestMessage(sid, curHeight-config.BlockPerMessage+1, config.BlockPerMessage)
				assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

				bdl1 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).From, curHeight-config.BlockPerMessage+1)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).To(), curHeight)
				assert.Equal(t, bdl1.Message.(*message.BlocksResponseMessage).Count(), config.BlockPerMessage)
				assert.Zero(t, bdl1.Message.(*message.BlocksResponseMessage).LastCertificateHeight())

				bdl2 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
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

			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeSynced)
		})
	})

	t.Run("NodeNetwork flag set", func(t *testing.T) {
		td.sync.config.NodeNetwork = true

		sid := td.RandInt(100)
		pid := td.RandPeerID()
		pub, _ := td.RandBLSKeyPair()
		td.addPeer(t, pub, pid, false)

		t.Run("Requesting one block", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 1, 2)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			msg1 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
			assert.Equal(t, msg1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)

			msg2 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
			assert.Equal(t, msg2.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeNoMoreBlocks)
		})

		t.Run("Peer is busy", func(t *testing.T) {
			td.sync.peerSet.OpenSession(td.RandPeerID())
			td.sync.peerSet.OpenSession(td.RandPeerID())
			td.sync.peerSet.OpenSession(td.RandPeerID())
			td.sync.peerSet.OpenSession(td.RandPeerID())
			td.sync.peerSet.OpenSession(td.RandPeerID())
			require.Equal(t, td.sync.peerSet.NumberOfOpenSessions(), 5)

			s := td.sync.peerSet.OpenSession(td.network.SelfID())
			msg := message.NewBlocksRequestMessage(s.SessionID(), 100, 105)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))
			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
		})
	})
}
