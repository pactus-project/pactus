package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/tx"
)

func TestTransactionsType(t *testing.T) {
	p := &TransactionsPayload{}
	assert.Equal(t, p.Type(), PayloadTypeTransactions)
}

func TestTransactionsPayload(t *testing.T) {
	t.Run("No transactions", func(t *testing.T) {
		p := NewTransactionsPayload(nil)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		p := NewTransactionsPayload([]*tx.Tx{trx})

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), trx.ID().Fingerprint())
	})
}
