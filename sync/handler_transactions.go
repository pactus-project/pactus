package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type transactionsHandler struct {
	*synchronizer
}

func newTransactionsHandler(sync *synchronizer) payloadHandler {
	return &transactionsHandler{
		sync,
	}
}

func (handler *transactionsHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.TransactionsPayload)
	handler.logger.Trace("Parsing transactions payload", "pld", pld)

	handler.cache.AddTransactions(pld.Transactions)

	for _, trx := range pld.Transactions {
		if err := handler.state.AddPendingTx(trx); err != nil {
			handler.logger.Debug("Cannot append transaction", "tx", trx, "err", err)

			// TODO: set peer as bad peer?
		}
	}

	return nil
}

func (handler *transactionsHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}
