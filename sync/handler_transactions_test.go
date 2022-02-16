package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingTransactionsMessages(t *testing.T) {
	setup(t)

	t.Run("Parsing transaction message", func(t *testing.T) {
		trx1, _ := tx.GenerateTestBondTx()
		pld := payload.NewTransactionsPayload([]*tx.Tx{trx1})

		assert.NoError(t, testReceiveingNewMessage(tSync, pld, util.RandomPeerID()))

		assert.NotNil(t, tSync.cache.GetTransaction(trx1.ID()))
		assert.NotNil(t, tSync.state.PendingTx(trx1.ID()))
	})
}
