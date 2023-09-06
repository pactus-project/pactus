package sync

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/services"
	"github.com/stretchr/testify/assert"
)

func TestLatestBlocksRequestMessages(t *testing.T) {
	config := testConfig()
	config.NodeNetwork = false

	td := setup(t, config)
	sid := td.RandInt(100)
	pid := td.RandPeerID()

	td.state.CommitTestBlocks(31)

	t.Run("NodeNetwork flag is not set", func(t *testing.T) {
		curHeight := td.state.LastBlockHeight()

		t.Run("Reject request from unknown peers", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, curHeight-1, 1)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).From, uint32(0))
		})

		pub, _ := td.RandBLSKeyPair()
		td.addPeer(t, pub, pid, services.New(services.None))

		t.Run("Reject requests not within `LatestBlockInterval`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 1, 2)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeRejected)
			assert.Equal(t, bdl.Message.(*message.BlocksResponseMessage).From, uint32(0))
		})

		t.Run("Request blocks more than `LatestBlockInterval`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 10, LatestBlockInterval+1)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

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

		t.Run("Requesting one block", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 1, 2)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			msg1 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
			assert.Equal(t, msg1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)

			msg2 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksResponse)
			assert.Equal(t, msg2.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeNoMoreBlocks)
		})
	})

	t.Run("Respond error", func(t *testing.T) {
		td.network.SendError = fmt.Errorf("send error")

		msg := message.NewBlocksRequestMessage(sid, 1, 2)
		err := td.receivingNewMessage(td.sync, msg, pid)
		assert.ErrorIs(t, err, td.network.SendError)
	})
}
