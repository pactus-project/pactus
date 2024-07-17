package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/stretchr/testify/assert"
)

func TestBlocksRequestMessages(t *testing.T) {
	t.Run("NetworkLimited service is enabled", func(t *testing.T) {
		config := testConfig()
		config.Services = service.Services(service.PrunedNode)

		td := setup(t, config)
		sid := td.RandInt(100)

		td.state.CommitTestBlocks(31)
		curHeight := td.state.LastBlockHeight()

		t.Run("Reject request from unknown peers", func(t *testing.T) {
			pid := td.RandPeerID()
			msg := message.NewBlocksRequestMessage(sid, curHeight-1, 1)
			td.receivingNewMessage(td.sync, msg, pid)

			bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			res := bdl.Message.(*message.BlocksResponseMessage)
			assert.Equal(t, message.ResponseCodeRejected, res.ResponseCode)
			assert.Contains(t, res.Reason, "unknown peer")
			assert.Zero(t, res.From)
			assert.Equal(t, res.SessionID, sid)
		})

		t.Run("Reject request from peers without handshaking", func(t *testing.T) {
			pid := td.addPeer(t, status.StatusConnected, service.New(service.None))
			msg := message.NewBlocksRequestMessage(sid, curHeight-1, 1)
			td.receivingNewMessage(td.sync, msg, pid)

			bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			res := bdl.Message.(*message.BlocksResponseMessage)
			assert.Equal(t, message.ResponseCodeRejected, res.ResponseCode)
			assert.Contains(t, res.Reason, "not handshaked")
		})

		pid := td.addPeer(t, status.StatusKnown, service.New(service.None))

		t.Run("Peer requested blocks that we don't have", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, curHeight+1, 1)
			td.receivingNewMessage(td.sync, msg, pid)

			bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			res := bdl.Message.(*message.BlocksResponseMessage)
			assert.Equal(t, message.ResponseCodeRejected, res.ResponseCode)
			assert.Contains(t, res.Reason, "requested blocks from 32 exceed current height 31")
		})

		t.Run("Request blocks more than `BlockPerSession`", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 10, config.BlockPerSession+1)
			td.receivingNewMessage(td.sync, msg, pid)

			bdl := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			res := bdl.Message.(*message.BlocksResponseMessage)
			assert.Equal(t, message.ResponseCodeRejected, res.ResponseCode)
			assert.Contains(t, res.Reason, "requested block range 10-33 exceeds the allowed 23 blocks per session")
		})

		t.Run("Accept request within `BlockPerSession`", func(t *testing.T) {
			t.Run("Peer needs more block", func(t *testing.T) {
				msg := message.NewBlocksRequestMessage(sid, curHeight-config.BlockPerMessage, config.BlockPerMessage)
				td.receivingNewMessage(td.sync, msg, pid)

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
				td.receivingNewMessage(td.sync, msg, pid)

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

	t.Run("Network service is enabled", func(t *testing.T) {
		config := testConfig()
		config.Services = service.New(service.FullNode)

		td := setup(t, config)
		sid := td.RandInt(100)

		td.state.CommitTestBlocks(31)
		pid := td.addPeer(t, status.StatusKnown, service.New(service.None))

		t.Run("Requesting one block", func(t *testing.T) {
			msg := message.NewBlocksRequestMessage(sid, 1, 2)
			td.receivingNewMessage(td.sync, msg, pid)

			msg1 := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			assert.Equal(t, msg1.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeMoreBlocks)

			msg2 := td.shouldPublishMessageWithThisType(t, message.TypeBlocksResponse)
			assert.Equal(t, msg2.Message.(*message.BlocksResponseMessage).ResponseCode, message.ResponseCodeNoMoreBlocks)
		})
	})
}
