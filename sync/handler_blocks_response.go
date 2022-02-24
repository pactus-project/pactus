package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/bundle"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
)

type blocksResponseHandler struct {
	*synchronizer
}

func newBlocksResponseHandler(sync *synchronizer) messageHandler {
	return &blocksResponseHandler{
		sync,
	}
}

func (handler *blocksResponseHandler) ParsMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.BlocksResponseMessage)
	handler.logger.Trace("parsing BlocksResponse message", "msg", msg)

	if msg.IsRequestRejected() {
		handler.logger.Warn("blocks request is rejected", "pid", initiator, "response", msg.ResponseCode)
	} else {
		handler.cache.AddCertificate(msg.LastCertificate)
		handler.cache.AddBlocks(msg.From, msg.Blocks)
		handler.cache.AddTransactions(msg.Transactions)
		handler.tryCommitBlocks()
	}
	handler.updateSession(msg.SessionID, initiator, msg.ResponseCode)

	return nil
}

func (handler *blocksResponseHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	msg := bundle.NewBundle(handler.SelfID(), m)
	msg.CompressIt()

	return msg
}
