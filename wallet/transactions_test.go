package wallet_test

// func TestGetHistory(t *testing.T) {
// 	td := setup(t)

// 	history := td.wallet.History(td.RandAccAddress().String())
// 	assert.Empty(t, history)
// }

// func TestAddDuplicatedTrx(t *testing.T) {
// 	td := setup(t)

// 	trx := td.GenerateTestTransferTx()
// 	id, err := td.wallet.BroadcastTransaction(trx)
// 	assert.NoError(t, err)
// 	assert.Equal(t, trx.ID().String(), id)

// 	history := td.wallet.History(trx.Payload().Signer().String())
// 	assert.Equal(t, id, history[0].TxID)
// }
