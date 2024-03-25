package sync

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
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

func (handler *blocksRequestHandler) ParseMessage(m message.Message, pid peer.ID) error {
	msg := m.(*message.BlocksRequestMessage)
	handler.logger.Trace("parsing BlocksRequest message", "msg", msg)

	p := handler.peerSet.GetPeer(pid)
	if p == nil {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			fmt.Sprintf("unknown peer (%s)", pid.String()), msg.SessionID, 0, nil, nil)

		handler.respond(response, pid)

		return nil
	}

	if !p.IsKnownOrTrusty() {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			fmt.Sprintf("not handshaked (%s)", p.Status.String()), msg.SessionID, 0, nil, nil)

		handler.respond(response, pid)

		return nil
	}

	ourHeight := handler.state.LastBlockHeight()
	if !handler.config.NodeNetwork {
		if ourHeight > handler.config.LatestBlockInterval && msg.From < ourHeight-handler.config.LatestBlockInterval {
			response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
				fmt.Sprintf("the request height is not acceptable: %v", msg.From), msg.SessionID, 0, nil, nil)

			handler.respond(response, pid)

			return nil
		}
	}

	if msg.From > ourHeight {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			fmt.Sprintf("don't have requested blocks: %v", msg.From), msg.SessionID, 0, nil, nil)

		handler.respond(response, pid)

		return nil
	}

	if msg.Count > handler.config.LatestBlockInterval {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			fmt.Sprintf("too many blocks requested: %v-%v", msg.From, msg.Count), msg.SessionID, 0, nil, nil)

		handler.respond(response, pid)

		return nil
	}

	// Help this peer to sync up
	height := msg.From
	count := msg.Count
	for {
		blockToRead := util.Min(handler.config.BlockPerMessage, count)
		blocksData := handler.prepareBlocks(height, blockToRead)
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

		return nil
	}

	response := message.NewBlocksResponseMessage(message.ResponseCodeNoMoreBlocks,
		message.ResponseCodeNoMoreBlocks.String(), msg.SessionID, 0, nil, nil)

	handler.respond(response, pid)

	return nil
}

func (handler *blocksRequestHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(m)
}

func (handler *blocksRequestHandler) respond(msg *message.BlocksResponseMessage, to peer.ID) {
	if msg.ResponseCode == message.ResponseCodeRejected {
		handler.logger.Debug("rejecting block request message", "msg", msg,
			"to", to, "reason", msg.Reason)

		handler.sendTo(msg, to)
	} else {
		handler.logger.Info("responding block request message", "msg", msg, "to", to)

		handler.sendTo(msg, to)
	}
}
