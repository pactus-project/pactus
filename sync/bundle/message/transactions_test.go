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
	assert.Equal(t, m.Type(), TypeTransactions)
}

func TestTransactionsMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("No transactions", func(t *testing.T) {
		m := NewTransactionsMessage(nil)

		assert.Equal(t, errors.Code(m.BasicCheck()), errors.ErrInvalidMessage)
	})

	t.Run("OK", func(t *testing.T) {
		trx, _ := ts.GenerateTestTransferTx()
		m := NewTransactionsMessage([]*tx.Tx{trx})

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), trx.ID().ShortString())
	})
}
