package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/stretchr/testify/assert"
)

func TestParsingTransactionsMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("Parsing transactions message", func(t *testing.T) {
		trx1, _ := td.GenerateTestBondTx()
		msg := message.NewTransactionsMessage([]*tx.Tx{trx1})

		assert.NoError(t, td.receivingNewMessage(td.sync, msg, td.RandPeerID()))

		assert.NotNil(t, td.sync.state.PendingTx(trx1.ID()))
	})
}
