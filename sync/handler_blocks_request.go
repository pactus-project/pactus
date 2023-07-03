package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
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

	if handler.peerSet.NumberOfOpenSessions() > handler.config.MaxOpenSessions {
		handler.logger.Warn("we are busy", "message", msg, "pid", initiator)
		response := message.NewBlocksResponseMessage(message.ResponseCodeBusy,
			msg.SessionID, 0, nil, nil)
		handler.sendTo(response, initiator, msg.SessionID)

		return nil
	}

	peer := handler.peerSet.GetPeer(initiator)
	if !peer.IsKnownOrTrusty() {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			msg.SessionID, 0, nil, nil)
		handler.sendTo(response, initiator, msg.SessionID)

		return errors.Errorf(errors.ErrInvalidMessage, "peer status is %v", peer.Status)
	}

	// TODO
	// This condition causes some troubles on testnet. Commenting it for further discussion.
	// Connections between the nodes can be interrupted or even manually restarted during sync.
	// When the number of "node networks" are less than of MaxOpenSessions, it can cause the node to never sync.
	// if peer.Height > msg.From {
	// 	response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
	// 		msg.SessionID, 0, nil, nil)
	// 	handler.sendTo(response, initiator)

	// 	return errors.Errorf(errors.ErrInvalidMessage, "peer request for blocks that already has: %v", msg.From)
	// }

	if !handler.config.NodeNetwork {
		ourHeight := handler.state.LastBlockHeight()
		if msg.From < ourHeight-LatestBlockInterval {
			response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
				msg.SessionID, 0, nil, nil)
			handler.sendTo(response, initiator, msg.SessionID)

			return errors.Errorf(errors.ErrInvalidMessage, "the request height is not acceptable: %v", msg.From)
		}
	}
	height := msg.From
	count := msg.Count

	if count > LatestBlockInterval {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			msg.SessionID, 0, nil, nil)
		handler.sendTo(response, initiator, msg.SessionID)

		return errors.Errorf(errors.ErrInvalidMessage, "too many blocks requested: %v-%v", msg.From, msg.Count)
	}

	// Help this peer to sync up
	for {
		blockToRead := util.MinU32(handler.config.BlockPerMessage, count)
		blocksData := handler.prepareBlocks(height, blockToRead)
		if len(blocksData) == 0 {
			break
		}

		response := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks,
			msg.SessionID, height, blocksData, nil)
		handler.sendTo(response, initiator, msg.SessionID)

		height += uint32(len(blocksData))
		count -= uint32(len(blocksData))
		if count <= 0 {
			break
		}
	}
	// To avoid sending blocks again, we update height for this peer
	// Height is always greater than zeo.
	peerHeight := height - 1
	handler.peerSet.UpdateHeight(initiator, peerHeight)

	if msg.To() >= handler.state.LastBlockHeight() {
		lastCertificate := handler.state.LastCertificate()
		response := message.NewBlocksResponseMessage(message.ResponseCodeSynced,
			msg.SessionID, peerHeight, nil, lastCertificate)
		handler.sendTo(response, initiator, msg.SessionID)
	} else {
		response := message.NewBlocksResponseMessage(message.ResponseCodeNoMoreBlocks,
			msg.SessionID, 0, nil, nil)
		handler.sendTo(response, initiator, msg.SessionID)
	}

	return nil
}

func (handler *blocksRequestHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(handler.SelfID(), m)
}
