package local_test

import (
	"testing"

	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet/provider/local"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLocalProvider(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	mockState := state.NewMockState(ts.MockController())
	prov := local.NewLocalBlockchainProvider(mockState)

	t.Run("last block height", func(t *testing.T) {
		mockState.EXPECT().LastBlockHeight().Return(types.Height(100))

		height, err := prov.LastBlockHeight()
		require.NoError(t, err)
		assert.Equal(t, types.Height(100), height)
	})

	t.Run("account", func(t *testing.T) {
		addr, acc := ts.GenerateTestAccount()
		mockState.EXPECT().AccountByAddress(addr).Return(acc, nil)

		got, err := prov.GetAccount(addr.String())
		require.NoError(t, err)
		assert.Equal(t, acc, got)

		_, err = prov.GetAccount("invalid")
		require.Error(t, err)

		unknown := ts.RandAccAddress()
		mockState.EXPECT().AccountByAddress(unknown).Return(nil, assert.AnError)
		_, err = prov.GetAccount(unknown.String())
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("validator", func(t *testing.T) {
		val := ts.GenerateTestValidator()
		mockState.EXPECT().ValidatorByAddress(val.Address()).Return(val, nil)

		got, err := prov.GetValidator(val.Address().String())
		require.NoError(t, err)
		assert.Equal(t, val, got)

		_, err = prov.GetValidator("invalid")
		require.Error(t, err)
	})

	t.Run("transaction", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		data, err := trx.Bytes()
		require.NoError(t, err)

		mockState.EXPECT().CommittedTx(trx.ID()).Return(
			&store.CommittedTx{TxID: trx.ID(), Height: 42, Data: data}, nil)

		got, height, err := prov.GetTransaction(trx.ID().String())
		require.NoError(t, err)
		assert.Equal(t, trx.ID(), got.ID())
		assert.Equal(t, types.Height(42), height)

		_, _, err = prov.GetTransaction("invalid")
		require.Error(t, err)

		// A transaction that is not committed.
		other := ts.GenerateTestTransferTx()
		mockState.EXPECT().CommittedTx(other.ID()).Return(nil, assert.AnError)
		_, _, err = prov.GetTransaction(other.ID().String())
		require.ErrorIs(t, err, assert.AnError)

		// A stored transaction that can't be decoded.
		mockState.EXPECT().CommittedTx(other.ID()).Return(
			&store.CommittedTx{TxID: other.ID(), Height: 42, Data: []byte("invalid")}, nil)
		_, _, err = prov.GetTransaction(other.ID().String())
		require.Error(t, err)
	})

	t.Run("send transaction", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		mockState.EXPECT().AddPendingTxAndBroadcast(trx).Return(nil)

		txID, err := prov.SendTx(trx)
		require.NoError(t, err)
		assert.Equal(t, trx.ID().String(), txID)

		mockState.EXPECT().AddPendingTxAndBroadcast(trx).Return(assert.AnError)
		_, err = prov.SendTx(trx)
		require.ErrorIs(t, err, assert.AnError)
	})

	t.Run("check transaction", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		data, err := trx.Bytes()
		require.NoError(t, err)

		mockState.EXPECT().CheckTransaction(gomock.Any()).Return(nil)
		require.NoError(t, prov.CheckTransaction(data))

		mockState.EXPECT().CheckTransaction(gomock.Any()).Return(assert.AnError)
		require.ErrorIs(t, prov.CheckTransaction(data), assert.AnError)

		require.Error(t, prov.CheckTransaction([]byte("invalid")))
	})

	require.NoError(t, prov.Close())
}
