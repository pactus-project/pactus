package message

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionsType(t *testing.T) {
	msg := &TransactionsMessage{}
	assert.Equal(t, TypeTransaction, msg.Type())
}

func TestTransactionsMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("No transactions", func(t *testing.T) {
		msg := NewTransactionsMessage(nil)

		err := msg.BasicCheck()
		require.ErrorIs(t, err, BasicCheckError{Reason: "no transaction"})
		assert.Contains(t, msg.LogString(), "[]")
	})

	t.Run("Nil transaction ", func(t *testing.T) {
		msg := NewTransactionsMessage([]*tx.Tx{nil})
		err := msg.BasicCheck()

		require.ErrorIs(t, err, BasicCheckError{Reason: "nil transaction"})
		assert.Contains(t, msg.LogString(), "nil-tx")
	})

	t.Run("Invalid transactions", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		trx.SetSignature(nil)

		msg := NewTransactionsMessage([]*tx.Tx{trx})

		err := msg.BasicCheck()
		require.ErrorIs(t, err, tx.BasicCheckError{Reason: "no signature"})
	})

	t.Run("OK", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		msg := NewTransactionsMessage([]*tx.Tx{trx})

		require.NoError(t, msg.BasicCheck())
		assert.Zero(t, msg.ConsensusHeight())
		assert.Contains(t, msg.LogString(), trx.ID().LogString())
	})
}
