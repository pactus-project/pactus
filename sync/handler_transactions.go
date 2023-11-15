package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
)

type transactionsHandler struct {
	*synchronizer
}

func newTransactionsHandler(sync *synchronizer) messageHandler {
	return &transactionsHandler{
		sync,
	}
}

func (handler *transactionsHandler) ParseMessage(m message.Message, _ peer.ID) error {
	msg := m.(*message.TransactionsMessage)
	handler.logger.Trace("parsing Transactions message", "msg", msg)

	for _, trx := range msg.Transactions {
		if err := handler.state.AddPendingTx(trx); err != nil {
			handler.logger.Debug("cannot append transaction", "tx", trx, "error", err)

			// TODO: set peer as bad peer?
		}
	}

	return nil
}

func (handler *transactionsHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(m)
}
