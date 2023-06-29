package message

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestTransactionsType(t *testing.T) {
	m := &TransactionsMessage{}
	assert.Equal(t, m.Type(), MessageTypeTransactions)
}

func TestTransactionsMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("No transactions", func(t *testing.T) {
		m := NewTransactionsMessage(nil)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidMessage)
	})

	t.Run("OK", func(t *testing.T) {
		trx, _ := ts.GenerateTestSendTx()
		m := NewTransactionsMessage([]*tx.Tx{trx})

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), trx.ID().Fingerprint())
	})
}
