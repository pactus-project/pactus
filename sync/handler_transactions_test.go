package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
)

func TestParsingTransactionsMessages(t *testing.T) {
	setup(t)
	disableHeartbeat(t)

	t.Run("Bob should add all transactions to cache", func(t *testing.T) {
		trx1, _ := tx.GenerateTestBondTx()

		// Alice send transaction to bob, bob should cache it
		pld := payload.NewTransactionsPayload([]*tx.Tx{trx1})
		tAliceSync.broadcast(pld)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeTransactions)
		assert.NotNil(t, tBobSync.cache.GetTransaction(trx1.ID()))
		assert.NotNil(t, tBobSync.state.PendingTx(trx1.ID()))
	})
}
