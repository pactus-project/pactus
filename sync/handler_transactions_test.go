package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/stretchr/testify/assert"
)

func TestHandlerTransactionsParsingMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("Parsing transactions message", func(t *testing.T) {
		trx1 := td.GenerateTestBondTx()
		msg := message.NewTransactionsMessage([]*tx.Tx{trx1})
		pid := td.RandPeerID()

		td.receivingNewMessage(td.sync, msg, pid)

		assert.NotNil(t, td.sync.state.PendingTx(trx1.ID()))
	})
}

func TestHandlerTransactionsBroadcastingMessages(t *testing.T) {
	td := setup(t, nil)

	trx1 := td.GenerateTestBondTx()
	msg := message.NewTransactionsMessage([]*tx.Tx{trx1})
	td.sync.broadcast(msg)

	td.shouldPublishMessageWithThisType(t, message.TypeTransaction)
}
