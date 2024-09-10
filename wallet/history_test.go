package wallet_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHistory(t *testing.T) {
	td := setup(t)
	defer td.Close()

	history := td.wallet.History(td.RandAccAddress().String())
	assert.Empty(t, history)
}

func TestAddDuplicatedTrx(t *testing.T) {
	td := setup(t)
	defer td.Close()

	trx := td.GenerateTestTransferTx()
	id, err := td.wallet.BroadcastTransaction(trx)
	assert.NoError(t, err)
	assert.Equal(t, trx.ID().String(), id)

	history := td.wallet.History(trx.Payload().Signer().String())
	assert.Equal(t, id, history[0].TxID)
}
