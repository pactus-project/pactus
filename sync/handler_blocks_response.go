package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type blocksResponseHandler struct {
	*synchronizer
}

func newBlocksResponseHandler(sync *synchronizer) payloadHandler {
	return &blocksResponseHandler{
		sync,
	}
}

func (handler *blocksResponseHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.BlocksResponsePayload)
	handler.logger.Trace("Parsing blocks response payload", "pld", pld)

	if pld.IsRequestNotProcessed() {
		handler.logger.Warn("Blocks request is rejected", "pid", initiator, "response", pld.ResponseCode)
	} else {
		handler.cache.AddCertificate(pld.LastCertificate)
		handler.cache.AddBlocks(pld.From, pld.Blocks)
		handler.cache.AddTransactions(pld.Transactions)
		handler.tryCommitBlocks()
	}
	handler.updateSession(pld.SessionID, initiator, pld.ResponseCode)

	return nil
}

func (handler *blocksResponseHandler) PrepareMessage(p payload.Payload) *message.Message {
	msg := message.NewMessage(handler.SelfID(), p)
	msg.CompressIt()

	return msg
}
