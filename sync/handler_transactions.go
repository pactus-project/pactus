package sync

import (
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
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
		}
	}

	return nil
}

func (handler *transactionsHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	return bundle.NewBundle(m)
}
