package message

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestTransactionsType(t *testing.T) {
	m := &TransactionsMessage{}
	assert.Equal(t, TypeTransaction, m.Type())
}

func TestTransactionsMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("No transactions", func(t *testing.T) {
		m := NewTransactionsMessage(nil)

		err := m.BasicCheck()
		assert.ErrorIs(t, err, BasicCheckError{Reason: "no transaction"})
	})

	t.Run("OK", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		m := NewTransactionsMessage([]*tx.Tx{trx})

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), trx.ID().ShortString())
	})
}
