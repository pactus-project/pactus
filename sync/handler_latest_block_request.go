package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
)

type latestBlocksRequestHandler struct {
	*synchronizer
}

func newLatestBlocksRequestHandler(sync *synchronizer) payloadHandler {
	return &latestBlocksRequestHandler{
		sync,
	}
}

func (handler *latestBlocksRequestHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.LatestBlocksRequestPayload)
	handler.logger.Trace("Parsing latest blocks request payload", "pld", pld)

	if pld.Target != handler.SelfID() {
		return nil
	}

	if handler.peerSet.NumberOfOpenSessions() > handler.config.MaximumOpenSessions {
		handler.logger.Warn("We are busy", "pld", pld, "pid", initiator)
		response := payload.NewLatestBlocksResponsePayload(payload.ResponseCodeBusy, pld.SessionID, initiator, 0, nil, nil, nil)
		handler.broadcast(response)

		return nil
	}

	peer := handler.peerSet.MustGetPeer(initiator)
	if peer.Status() != peerset.StatusCodeOK {
		response := payload.NewLatestBlocksResponsePayload(payload.ResponseCodeRejected, pld.SessionID, initiator, 0, nil, nil, nil)
		handler.broadcast(response)

		return errors.Errorf(errors.ErrInvalidMessage, "Peer status is not ok: %v", peer.Status())
	}

	if peer.Height() > pld.From {
		response := payload.NewLatestBlocksResponsePayload(payload.ResponseCodeRejected, pld.SessionID, initiator, 0, nil, nil, nil)
		handler.broadcast(response)

		return errors.Errorf(errors.ErrInvalidMessage, "Peer request for blocks that already has: %v", pld.From)
	}

	ourHeight := handler.state.LastBlockHeight()
	if pld.From < ourHeight-LatestBlockInterval {
		response := payload.NewLatestBlocksResponsePayload(payload.ResponseCodeRejected, pld.SessionID, initiator, 0, nil, nil, nil)
		handler.broadcast(response)

		return errors.Errorf(errors.ErrInvalidMessage, "the request height is not acceptable: %v", pld.From)
	}
	from := pld.From
	count := handler.config.BlockPerMessage

	// Help peer to catch up
	for {
		blocks, trxs := handler.prepareBlocksAndTransactions(from, count)
		if len(blocks) == 0 {
			break
		}

		response := payload.NewLatestBlocksResponsePayload(payload.ResponseCodeMoreBlocks, pld.SessionID, initiator, from, blocks, trxs, nil)
		handler.broadcast(response)

		from += len(blocks)
	}
	// To avoid sending blocks again, we update height for this peer
	peer.UpdateHeight(from - 1)

	lastCertificate := handler.state.LastCertificate()
	response := payload.NewLatestBlocksResponsePayload(payload.ResponseCodeSynced, pld.SessionID, initiator, from, nil, nil, lastCertificate)
	handler.broadcast(response)

	return nil
}

func (handler *latestBlocksRequestHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}
