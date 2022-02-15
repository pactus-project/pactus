package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type blocksRequestHandler struct {
	*synchronizer
}

func newBlocksRequestHandler(sync *synchronizer) payloadHandler {
	return &blocksRequestHandler{
		sync,
	}
}

func (handler *blocksRequestHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.BlocksRequestPayload)
	handler.logger.Trace("Parsing blocks request payload", "pld", pld)

	if handler.peerSet.NumberOfOpenSessions() > handler.config.MaximumOpenSessions {
		handler.logger.Warn("We are busy", "pld", pld, "pid", initiator)
		response := payload.NewBlocksResponsePayload(payload.ResponseCodeBusy, pld.SessionID, 0, nil, nil, nil)
		handler.sendTo(response, initiator)

		return nil
	}

	peer := handler.peerSet.MustGetPeer(initiator)
	if !peer.IsKnownOrTrusted() {
		response := payload.NewBlocksResponsePayload(payload.ResponseCodeRejected, pld.SessionID, 0, nil, nil, nil)
		handler.sendTo(response, initiator)

		return errors.Errorf(errors.ErrInvalidMessage, "Peer status is %v", peer.Status())
	}

	if peer.Height() > pld.From {
		response := payload.NewBlocksResponsePayload(payload.ResponseCodeRejected, pld.SessionID, 0, nil, nil, nil)
		handler.sendTo(response, initiator)

		return errors.Errorf(errors.ErrInvalidMessage, "Peer request for blocks that already has: %v", pld.From)
	}

	if !handler.config.InitialBlockDownload {
		ourHeight := handler.state.LastBlockHeight()
		if pld.From < ourHeight-LatestBlockInterval {
			response := payload.NewBlocksResponsePayload(payload.ResponseCodeRejected, pld.SessionID, 0, nil, nil, nil)
			handler.sendTo(response, initiator)

			return errors.Errorf(errors.ErrInvalidMessage, "the request height is not acceptable: %v", pld.From)
		}
	}
	height := pld.From
	count := handler.config.BlockPerMessage

	// Help peer to catch up
	for {
		blocks, trxs := handler.prepareBlocksAndTransactions(height, count)
		if len(blocks) == 0 {
			break
		}

		response := payload.NewBlocksResponsePayload(payload.ResponseCodeMoreBlocks, pld.SessionID, height, blocks, trxs, nil)
		handler.sendTo(response, initiator)

		height += len(blocks)
		if height >= pld.To {
			break
		}
	}
	// To avoid sending blocks again, we update height for this peer
	peer.UpdateHeight(height - 1)

	if pld.To >= handler.state.LastBlockHeight() {
		lastCertificate := handler.state.LastCertificate()
		response := payload.NewBlocksResponsePayload(payload.ResponseCodeSynced, pld.SessionID, handler.state.LastBlockHeight(), nil, nil, lastCertificate)
		handler.sendTo(response, initiator)
	} else {
		response := payload.NewBlocksResponsePayload(payload.ResponseCodeNoMoreBlocks, pld.SessionID, 0, nil, nil, nil)
		handler.sendTo(response, initiator)
	}

	return nil
}

func (handler *blocksRequestHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}
