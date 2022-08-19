package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/bundle"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/util/errors"
)

type blocksRequestHandler struct {
	*synchronizer
}

func newBlocksRequestHandler(sync *synchronizer) messageHandler {
	return &blocksRequestHandler{
		sync,
	}
}

func (handler *blocksRequestHandler) ParsMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.BlocksRequestMessage)
	handler.logger.Trace("parsing BlocksRequest message", "message", msg)

	if handler.peerSet.NumberOfOpenSessions() > handler.config.MaxOpenSessions {
		handler.logger.Warn("we are busy", "message", msg, "pid", initiator)
		response := message.NewBlocksResponseMessage(message.ResponseCodeBusy,
			msg.SessionID, 0, nil, nil)
		handler.sendTo(response, initiator)

		return nil
	}

	peer := handler.peerSet.GetPeer(initiator)
	if !peer.IsKnownOrTrusty() {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			msg.SessionID, 0, nil, nil)
		handler.sendTo(response, initiator)

		return errors.Errorf(errors.ErrInvalidMessage, "peer status is %v", peer.Status)
	}

	if peer.Height > msg.From {
		response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
			msg.SessionID, 0, nil, nil)
		handler.sendTo(response, initiator)

		return errors.Errorf(errors.ErrInvalidMessage, "peer request for blocks that already has: %v", msg.From)
	}

	if !handler.config.NodeNetwork {
		ourHeight := handler.state.LastBlockHeight()
		if msg.From < ourHeight-LatestBlockInterval {
			response := message.NewBlocksResponseMessage(message.ResponseCodeRejected,
				msg.SessionID, 0, nil, nil)
			handler.sendTo(response, initiator)

			return errors.Errorf(errors.ErrInvalidMessage, "the request height is not acceptable: %v", msg.From)
		}
	}
	height := msg.From
	count := handler.config.BlockPerMessage

	// Help peer to catch up
	for {
		blocks := handler.prepareBlocks(height, count)
		if len(blocks) == 0 {
			break
		}

		response := message.NewBlocksResponseMessage(message.ResponseCodeMoreBlocks,
			msg.SessionID, height, blocks, nil)
		handler.sendTo(response, initiator)

		height += uint32(len(blocks))
		if height >= msg.To {
			break
		}
	}
	// To avoid sending blocks again, we update height for this peer
	handler.peerSet.UpdateHeight(initiator, height-1)

	if msg.To >= handler.state.LastBlockHeight() {
		lastCertificate := handler.state.LastCertificate()
		response := message.NewBlocksResponseMessage(message.ResponseCodeSynced,
			msg.SessionID, handler.state.LastBlockHeight(), nil, lastCertificate)
		handler.sendTo(response, initiator)
	} else {
		response := message.NewBlocksResponseMessage(message.ResponseCodeNoMoreBlocks,
			msg.SessionID, 0, nil, nil)
		handler.sendTo(response, initiator)
	}

	return nil
}

func (handler *blocksRequestHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(handler.SelfID(), m)
}
