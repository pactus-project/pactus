package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/tx"
)

func TestTransactionsType(t *testing.T) {
	m := &TransactionsMessage{}
	assert.Equal(t, m.Type(), MessageTypeTransactions)
}

func TestTransactionsMessage(t *testing.T) {
	t.Run("No transactions", func(t *testing.T) {
		m := NewTransactionsMessage(nil)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		m := NewTransactionsMessage([]*tx.Tx{trx})

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), trx.ID().Fingerprint())
	})
}
