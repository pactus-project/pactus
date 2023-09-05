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

func (handler *blocksRequestHandler) ParseMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.BlocksRequestMessage)
	handler.logger.Trace("parsing BlocksRequest message", "message", msg)

	peer := handler.peerSet.GetPeer(initiator)
	if !peer.IsKnownOrTrusty() {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			fmt.Sprintf("peer status is %v", peer.Status), msg.SessionID, 0, nil, nil)

		return handler.respond(response, initiator, msg.SessionID)
	}

	if !handler.config.NodeNetwork {
		ourHeight := handler.state.LastBlockHeight()
		if msg.From < ourHeight-LatestBlockInterval {
			response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
				fmt.Sprintf("the request height is not acceptable: %v", msg.From), msg.SessionID, 0, nil, nil)

			return handler.respond(response, initiator, msg.SessionID)
		}
	}
	height := msg.From
	count := msg.Count

	if count > LatestBlockInterval {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			fmt.Sprintf("too many blocks requested: %v-%v", msg.From, msg.Count), msg.SessionID, 0, nil, nil)

		return handler.respond(response, initiator, msg.SessionID)
	}

	// Help this peer to sync up
	for {
		blockToRead := util.Min(handler.config.BlockPerMessage, count)
		blocksData := handler.prepareBlocks(height, blockToRead)
		if len(blocksData) == 0 {
			break
		}

		response := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks,
			message.ResponseCodeMoreBlocks.String(), msg.SessionID, height, blocksData, nil)
		err := handler.respond(response, initiator, msg.SessionID)
		if err != nil {
			return err
		}

		height += uint32(len(blocksData))
		count -= uint32(len(blocksData))
		if count <= 0 {
			break
		}
	}
	// To avoid sending blocks again, we update height for this peer
	// Height is always greater than zeo.
	peerHeight := height - 1

	if msg.To() >= handler.state.LastBlockHeight() {
		lastCert := handler.state.LastCertificate()
		response := message.NewBlocksResponseMessage(message.ResponseCodeSynced,
			message.ResponseCodeSynced.String(), msg.SessionID, peerHeight, nil, lastCert)

		return handler.respond(response, initiator, msg.SessionID)
	}

	response := message.NewBlocksResponseMessage(message.ResponseCodeNoMoreBlocks,
		message.ResponseCodeNoMoreBlocks.String(), msg.SessionID, 0, nil, nil)

	return handler.respond(response, initiator, msg.SessionID)
}

func (handler *blocksRequestHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(handler.SelfID(), m)
}

func (handler *blocksRequestHandler) respond(msg *message.BlocksResponseMessage, to peer.ID, sessionID int) error {
	if msg.ResponseCode == message.ResponseCodeRejected {
		handler.logger.Warn("rejecting block request message", "message", msg, "to", to, "reason", msg.Reason)
	} else {
		handler.logger.Info("responding block request message", "message", msg, "to", to)
	}

	if err := handler.sendTo(msg, to); err != nil {
		// TODO: test me
		// Let's close the session with this peer because we couldn't establish a connection.
		// This helps to free sessions and ask blocks from other peers.
		handler.peerSet.CloseSession(sessionID)
		return err
	}

	return nil
}
