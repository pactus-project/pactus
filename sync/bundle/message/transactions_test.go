package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/util/errors"
)

func TestTransactionsType(t *testing.T) {
	m := &TransactionsMessage{}
	assert.Equal(t, m.Type(), MessageTypeTransactions)
}

func TestTransactionsMessage(t *testing.T) {
	t.Run("No transactions", func(t *testing.T) {
		m := NewTransactionsMessage(nil)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidMessage)
	})

	t.Run("OK", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		m := NewTransactionsMessage([]*tx.Tx{trx})

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), trx.ID().Fingerprint())
	})
}
