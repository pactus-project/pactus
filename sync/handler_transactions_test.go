package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/tx"
)

func TestHandlerTransactionsParsingMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("Parsing transactions message", func(*testing.T) {
		trx1 := td.GenerateTestBondTx()
		msg := message.NewTransactionsMessage([]*tx.Tx{trx1})
		pid := td.RandPeerID()

		td.state.MockTxPool.EXPECT().AppendTx(trx1).Return(nil).Times(1)
		td.receivingNewMessage(td.sync, msg, pid)
	})
}

func TestHandlerTransactionsBroadcastingMessages(t *testing.T) {
	td := setup(t, nil)

	trx1 := td.GenerateTestBondTx()
	msg := message.NewTransactionsMessage([]*tx.Tx{trx1})
	td.sync.broadcast(msg)

	td.shouldPublishMessageWithThisType(t, message.TypeTransaction)
}
