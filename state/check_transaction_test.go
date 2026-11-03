package state

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/stretchr/testify/require"
)

func TestCheckTransactionRunsBasicCheck(t *testing.T) {
	td := setup(t)

	sender := td.genAccKey.PublicKeyNative().AccountAddress()
	lockTime := td.state.LastBlockHeight()

	t.Run("signed transaction is accepted", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, sender, td.RandAccAddress(), 10e9, 1e7)
		td.HelperSignTransaction(td.genAccKey, trx)

		require.NoError(t, td.state.CheckTransaction(trx))
	})

	t.Run("unsigned transaction", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, sender, td.RandAccAddress(), 10e9, 1e7)

		require.Error(t, td.state.CheckTransaction(trx))
	})

	t.Run("signed by a key that does not own the sender", func(t *testing.T) {
		_, otherPrv := td.RandEd25519KeyPair()
		trx := tx.NewTransferTx(lockTime, sender, td.RandAccAddress(), 10e9, 1e7)
		td.HelperSignTransaction(otherPrv, trx)

		require.Error(t, td.state.CheckTransaction(trx))
	})
}
