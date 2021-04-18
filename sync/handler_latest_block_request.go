package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
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

	peer := handler.peerSet.MustGetPeer(initiator)
	peer.UpdateHeight(pld.From)

	if pld.Target != handler.SelfID() {
		return nil
	}
	ourHeight := handler.state.LastBlockHeight()
	if pld.From < ourHeight-handler.config.RequestBlockInterval {
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

	lastCertificate := handler.state.LastCertificate()
	response := payload.NewLatestBlocksResponsePayload(payload.ResponseCodeSynced, pld.SessionID, initiator, from, nil, nil, lastCertificate)
	handler.broadcast(response)

	return nil
}

func (handler *latestBlocksRequestHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}
