package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/tx"
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
