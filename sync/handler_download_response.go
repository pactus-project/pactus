package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

type downloadResponseHandler struct {
	*synchronizer
}

func newDownloadResponseHandler(sync *synchronizer) payloadHandler {
	return &downloadResponseHandler{
		sync,
	}
}

func (handler *downloadResponseHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.DownloadResponsePayload)
	handler.logger.Trace("Parsing download response payload", "pld", pld)

	// TODO: to increase performance we can watch block stream
	// And add them into cache if they are close to our height
	if pld.Target != handler.SelfID() {
		return nil
	}

	if pld.IsRequestNotProcessed() {
		handler.logger.Warn("Download blocks request is rejected", "pid", util.FingerprintPeerID(initiator), "response", pld.ResponseCode)
	} else {
		handler.cache.AddBlocks(pld.From, pld.Blocks)
		handler.cache.AddTransactions(pld.Transactions)
		handler.tryCommitBlocks()
	}
	handler.updateSession(pld.ResponseCode, pld.SessionID, initiator, pld.Target)

	return nil
}

func (handler *downloadResponseHandler) PrepareMessage(p payload.Payload) *message.Message {
	msg := message.NewMessage(handler.SelfID(), p)
	msg.CompressIt()

	return msg
}
