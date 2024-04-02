package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/peerset/service"
	"github.com/stretchr/testify/assert"
)

func TestBlocksRequestMessages(t *testing.T) {
	config := testConfig()
	config.NodeNetwork = false

	td := setup(t, config)
	sid := td.RandInt(100)

	td.state.CommitTestBlocks(31)

	t.Run("NodeNetwork flag is not set", func(t *testing.T) {
		curHeight := td.state.LastBlockHeight()

		t.Run("Reject request from unknown peers", func(t *testing.T) {
			pid := td.RandPeerID()
			msg := message.NewBlocksRequestMessage(sid, curHeight-1, 1)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			res := bdl.Message.(*message.BlocksResponseMessage)
			assert.Equal(t, message.ResponseCodeRejected, res.ResponseCode)
			assert.Contains(t, res.Reason, "unknown peer")
			assert.Zero(t, res.From)
			assert.Equal(t, res.SessionID, sid)
		})

		t.Run("Reject request from peers without handshaking", func(t *testing.T) {
			pid := td.addPeer(t, peerset.StatusCodeConnected, service.New(service.None))
			msg := message.NewBlocksRequestMessage(sid, curHeight-1, 1)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			res := bdl.Message.(*message.BlocksResponseMessage)
			assert.Equal(t, message.ResponseCodeRejected, res.ResponseCode)
			assert.Contains(t, res.Reason, "not handshaked")
		})

		pid := td.addPeer(t, peerset.StatusCodeKnown, service.New(service.None))

		t.Run("Peer requested blocks that we don't have", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, curHeight+1, 1)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			res := bdl.Message.(*message.BlocksResponseMessage)
			assert.Equal(t, message.ResponseCodeRejected, res.ResponseCode)
			assert.Contains(t, res.Reason, "don't have requested block")
		})

		t.Run("Reject requests not within `LatestBlockInterval`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, curHeight-config.LatestBlockInterval-1, 1)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			res := bdl.Message.(*message.BlocksResponseMessage)
			assert.Equal(t, message.ResponseCodeRejected, res.ResponseCode)
			assert.Contains(t, res.Reason, "the request height is not acceptable")
		})

		t.Run("Request blocks more than `LatestBlockInterval`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 10, config.LatestBlockInterval+1)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			res := bdl.Message.(*message.BlocksResponseMessage)
			assert.Equal(t, message.ResponseCodeRejected, res.ResponseCode)
			assert.Contains(t, res.Reason, "too many blocks requested")
		})

		t.Run("Accept request within `LatestBlockInterval`", func(t *testing.T) {
			t.Run("Peer needs more block", func(t *testing.T) {
				msg := message.NewBlocksRequestMessage(sid, curHeight-config.BlockPerMessage, config.BlockPerMessage)
				assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

				bdl1 := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
				res1 := bdl1.Message.(*message.BlocksResponseMessage)
				assert.Equal(t, res1.ResponseCode, message.ResponseCodeMoreBlocks)
				assert.Equal(t, res1.From, curHeight-config.BlockPerMessage)
				assert.Equal(t, res1.To(), curHeight-1)
				assert.Equal(t, res1.Count(), config.BlockPerMessage)

				bdl2 := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
				res2 := bdl2.Message.(*message.BlocksResponseMessage)
				assert.Equal(t, res2.ResponseCode, message.ResponseCodeNoMoreBlocks)
				assert.Zero(t, res2.From)
				assert.Zero(t, res2.To())
				assert.Zero(t, res2.Count())
			})

			t.Run("Peer synced", func(t *testing.T) {
				msg := message.NewBlocksRequestMessage(sid, curHeight-config.BlockPerMessage+1, config.BlockPerMessage)
				assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

				bdl1 := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
				res1 := bdl1.Message.(*message.BlocksResponseMessage)
				assert.Equal(t, res1.ResponseCode, message.ResponseCodeMoreBlocks)
				assert.Equal(t, res1.From, curHeight-config.BlockPerMessage+1)
				assert.Equal(t, res1.To(), curHeight)
				assert.Equal(t, res1.Count(), config.BlockPerMessage)

				bdl2 := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
				res2 := bdl2.Message.(*message.BlocksResponseMessage)
				assert.Equal(t, res2.ResponseCode, message.ResponseCodeSynced)
				assert.Equal(t, res2.From, curHeight)
				assert.Zero(t, res2.To())
				assert.Zero(t, res2.Count())
			})
		})
	})

	t.Run("NodeNetwork flag set", func(t *testing.T) {
		td.sync.config.NodeNetwork = true
		pid := td.addPeer(t, peerset.StatusCodeKnown, service.New(service.None))

		t.Run("Requesting one block", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 1, 2)
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			msg1 := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			assert.Equal(t, msg1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)

			msg2 := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			assert.Equal(t, msg2.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeNoMoreBlocks)
		})
	})
}
