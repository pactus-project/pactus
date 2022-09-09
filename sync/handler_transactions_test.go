package sync

import (
	"testing"

	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/stretchr/testify/assert"
)

func TestParsingTransactionsMessages(t *testing.T) {
	setup(t)

	t.Run("Parsing transactions message", func(t *testing.T) {
		trx1, _ := tx.GenerateTestBondTx()
		msg := message.NewTransactionsMessage([]*tx.Tx{trx1})

		assert.NoError(t, testReceiveingNewMessage(tSync, msg, network.TestRandomPeerID()))

		assert.NotNil(t, tSync.state.PendingTx(trx1.ID()))
	})
}
