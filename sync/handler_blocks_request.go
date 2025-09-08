package sync

import (
	"fmt"

	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/util"
)

type blocksRequestHandler struct {
	*synchronizer
}

func newBlocksRequestHandler(sync *synchronizer) messageHandler {
	return &blocksRequestHandler{
		sync,
	}
}

func (handler *blocksRequestHandler) ParseMessage(m message.Message, pid peer.ID) {
	msg := m.(*message.BlocksRequestMessage)
	handler.logger.Trace("parsing BlocksRequest message", "msg", msg)

	peer := handler.peerSet.GetPeer(pid)
	if peer == nil {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			fmt.Sprintf("unknown peer (%s)", pid.String()), msg.SessionID, 0, nil, nil)

		handler.respond(response, pid)

		return
	}

	if !peer.Status.IsKnown() {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			fmt.Sprintf("not handshaked (%s)", peer.Status.String()), msg.SessionID, 0, nil, nil)

		handler.respond(response, pid)

		return
	}

	ourHeight := handler.state.LastBlockHeight()
	if msg.From > ourHeight {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			fmt.Sprintf("requested blocks from %v exceed current height %v",
				msg.From, ourHeight), msg.SessionID, 0, nil, nil)

		handler.respond(response, pid)

		return
	}

	if msg.Count > handler.config.BlockPerSession {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			fmt.Sprintf("requested block range %v-%v exceeds the allowed %v blocks per session",
				msg.From, msg.To(), handler.config.BlockPerSession), msg.SessionID, 0, nil, nil)

		handler.respond(response, pid)

		return
	}

	// Help this peer to sync up
	height := msg.From
	count := msg.Count
	for {
		blockCount := util.Min(handler.config.BlockPerMessage, count)
		blocksData := handler.prepareBlocks(height, blockCount)
		if len(blocksData) == 0 {
			break
		}

		response := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks,
			message.ResponseCodeMoreBlocks.String(), msg.SessionID, height, blocksData, nil)
		handler.respond(response, pid)

		height += uint32(len(blocksData))
		count -= uint32(len(blocksData))
		if count <= 0 {
			break
		}
	}

	if msg.To() >= ourHeight {
		lastCert := handler.state.LastCertificate()
		response := message.NewBlocksResponseMessage(message.ResponseCodeSynced,
			message.ResponseCodeSynced.String(), msg.SessionID, lastCert.Height(), nil, lastCert)

		handler.respond(response, pid)

		return
	}

	response := message.NewBlocksResponseMessage(message.ResponseCodeNoMoreBlocks,
		message.ResponseCodeNoMoreBlocks.String(), msg.SessionID, 0, nil, nil)

	handler.respond(response, pid)
}

func (*blocksRequestHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(m)
}

func (handler *blocksRequestHandler) respond(msg *message.BlocksResponseMessage, pid peer.ID) {
	if msg.ResponseCode == message.ResponseCodeRejected {
		handler.logger.Debug("rejecting block request message", "msg", msg,
			"pid", pid, "reason", msg.Reason)

		handler.sendTo(msg, pid)
	} else {
		handler.logger.Info("responding block request message", "msg", msg, "pid", pid)

		handler.sendTo(msg, pid)
	}
}
