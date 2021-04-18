package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type latestBlocksResponseHandler struct {
	*synchronizer
}

func newLatestBlocksResponseHandler(sync *synchronizer) payloadHandler {
	return &latestBlocksResponseHandler{
		sync,
	}
}

func (handler *latestBlocksResponseHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.LatestBlocksResponsePayload)
	handler.logger.Trace("Parsing latest blocks response payload", "pld", pld)

	// TODO: to increase performance we can watch block stream
	// And add them into cache if they are close to our height
	if pld.Target != handler.SelfID() {
		return nil
	}

	handler.cache.AddCertificate(pld.LastCertificate)
	handler.cache.AddBlocks(pld.From, pld.Blocks)
	handler.cache.AddTransactions(pld.Transactions)
	handler.tryCommitBlocks()

	handler.updateSession(pld.ResponseCode, pld.SessionID, initiator, pld.Target)

	return nil
}

func (handler *latestBlocksResponseHandler) PrepareMessage(p payload.Payload) *message.Message {
	msg := message.NewMessage(handler.SelfID(), p)
	msg.CompressIt()

	return msg
}
