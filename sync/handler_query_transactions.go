package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type queryTransactionsHandler struct {
	*synchronizer
}

func newQueryTransactionsHandler(sync *synchronizer) payloadHandler {
	return &queryTransactionsHandler{
		sync,
	}
}

func (handler *queryTransactionsHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.QueryTransactionsPayload)
	handler.logger.Trace("Parsing query transactions payload", "pld", pld)

	if !handler.peerIsInTheCommittee(initiator) {
		return errors.Errorf(errors.ErrInvalidMessage, "peers is not in the commmittee")
	}

	trxs := handler.prepareTransactions(pld.IDs)
	if len(trxs) > 0 {
		response := payload.NewTransactionsPayload(trxs)
		handler.broadcast(response)
	}

	return nil
}

func (handler *queryTransactionsHandler) PrepareMessage(p payload.Payload) *message.Message {
	pld := p.(*payload.QueryTransactionsPayload)
	missed := []hash.Hash{}
	for _, id := range pld.IDs {
		trx := handler.cache.GetTransaction(id)
		if trx != nil {
			if err := handler.state.AddPendingTx(trx); err != nil {
				handler.logger.Warn("Query for an invalid transaction", "tx", trx, "err", err)
			}
		} else {
			missed = append(missed, id)
		}
	}

	if !handler.weAreInTheCommittee() {
		handler.logger.Debug("Sending QueryTransactions ignored. We are not in the committee")
		return nil
	}

	if len(missed) == 0 {
		return nil
	}
	pld.IDs = missed
	return message.NewMessage(handler.SelfID(), p)
}
