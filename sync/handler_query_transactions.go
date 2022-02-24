package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/bundle"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
)

type queryTransactionsHandler struct {
	*synchronizer
}

func newQueryTransactionsHandler(sync *synchronizer) messageHandler {
	return &queryTransactionsHandler{
		sync,
	}
}

func (handler *queryTransactionsHandler) ParsMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.QueryTransactionsMessage)
	handler.logger.Trace("parsing QueryTransactions message", "msg", msg)

	if !handler.peerIsInTheCommittee(initiator) {
		return errors.Errorf(errors.ErrInvalidMessage, "peers is not in the commmittee")
	}

	trxs := handler.prepareTransactions(msg.IDs)
	if len(trxs) > 0 {
		response := message.NewTransactionsMessage(trxs)
		handler.broadcast(response)
	}

	return nil
}

func (handler *queryTransactionsHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	msg := m.(*message.QueryTransactionsMessage)
	missed := []hash.Hash{}
	for _, id := range msg.IDs {
		trx := handler.cache.GetTransaction(id)
		if trx != nil {
			if err := handler.state.AddPendingTx(trx); err != nil {
				handler.logger.Warn("query for an invalid transaction", "tx", trx, "err", err)
			}
		} else {
			missed = append(missed, id)
		}
	}

	if !handler.weAreInTheCommittee() {
		handler.logger.Debug("sending QueryTransactions ignored. We are not in the committee")
		return nil
	}

	if len(missed) == 0 {
		return nil
	}
	msg.IDs = missed
	return bundle.NewBundle(handler.SelfID(), m)
}
