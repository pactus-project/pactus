package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
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

	peer := handler.peerSet.MustGetPeer(initiator)
	peer.UpdateHeight(pld.From)

	if pld.Target != handler.SelfID() {
		return nil
	}
	if pld.To-pld.From > handler.config.RequestBlockInterval {
		return errors.Errorf(errors.ErrInvalidMessage, "Peer request interval is not acceptable: %v", pld.To-pld.From)
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

	response := payload.NewDownloadResponsePayload(payload.ResponseCodeNoMoreBlocks, pld.SessionID, initiator, 0, nil, nil)
	handler.broadcast(response)

	return nil
}

func (handler *downloadRequestHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}
