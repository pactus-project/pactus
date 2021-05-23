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
		p1 := NewTransactionsPayload(nil)
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		p2 := NewTransactionsPayload([]*tx.Tx{trx})
		assert.NoError(t, p2.SanityCheck())
	})
}
