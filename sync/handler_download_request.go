package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
)

type downloadRequestHandler struct {
	*synchronizer
}

func newDownloadRequestHandler(sync *synchronizer) payloadHandler {
	return &downloadRequestHandler{
		sync,
	}
}

func (handler *downloadRequestHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.DownloadRequestPayload)
	handler.logger.Trace("Parsing download request payload", "pld", pld)

	if pld.Target != handler.SelfID() {
		return nil
	}

	if handler.peerSet.NumberOfOpenSessions() > handler.config.MaximumOpenSessions {
		handler.logger.Warn("We are busy", "pld", pld, "pid", initiator)
		response := payload.NewDownloadResponsePayload(payload.ResponseCodeBusy, pld.SessionID, initiator, 0, nil, nil)
		handler.broadcast(response)

		return nil
	}

	peer := handler.peerSet.MustGetPeer(initiator)
	if peer.Status() != peerset.StatusCodeOK {
		response := payload.NewDownloadResponsePayload(payload.ResponseCodeRejected, pld.SessionID, initiator, 0, nil, nil)
		handler.broadcast(response)

		return errors.Errorf(errors.ErrInvalidMessage, "Peer status is not ok: %v", peer.Status())
	}

	if peer.Height() > pld.From {
		response := payload.NewDownloadResponsePayload(payload.ResponseCodeRejected, pld.SessionID, initiator, 0, nil, nil)
		handler.broadcast(response)

		return errors.Errorf(errors.ErrInvalidMessage, "Peer request for blocks that already has: %v", pld.From)
	}

	if pld.To-pld.From > handler.config.RequestBlockInterval {
		response := payload.NewDownloadResponsePayload(payload.ResponseCodeRejected, pld.SessionID, initiator, 0, nil, nil)
		handler.broadcast(response)

		return errors.Errorf(errors.ErrInvalidMessage, "peer request interval is not acceptable: %v", pld.To-pld.From)
	}

	from := pld.From
	count := handler.config.BlockPerMessage

	for {
		if from+count >= pld.To {
			// Last packet has one extra block, for confirming last block
			count++
		}
		blocks, trxs := handler.prepareBlocksAndTransactions(from, count)
		if len(blocks) == 0 {
			break
		}

		response := payload.NewDownloadResponsePayload(payload.ResponseCodeMoreBlocks, pld.SessionID, initiator, from, blocks, trxs)
		handler.broadcast(response)

		from += len(blocks)
		if from >= pld.To {
			break
		}
	}
	// To avoid sending blocks again, we update height for this peer
	peer.UpdateHeight(from - 1)

	response := payload.NewDownloadResponsePayload(payload.ResponseCodeNoMoreBlocks, pld.SessionID, initiator, 0, nil, nil)
	handler.broadcast(response)

	return nil
}

func (handler *downloadRequestHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}
