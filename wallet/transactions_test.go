package wallet

import (
	"errors"
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTransactionsListOptions(t *testing.T) {
	tests := []struct {
		opts     []ListTransactionsOption
		expDir   types.TxDirection
		expAddr  string
		expCount int
		expSkip  int
	}{
		{
			opts:   []ListTransactionsOption{},
			expDir: types.TxDirectionAny, expAddr: "*", expCount: 10, expSkip: 0,
		},
		{
			opts: []ListTransactionsOption{
				WithDirection(types.TxDirectionAny),
				WithAddress(""),
				WithCount(-1),
				WithSkip(-1),
			},
			expDir: types.TxDirectionAny, expAddr: "*", expCount: 10, expSkip: 0,
		},
		{
			opts: []ListTransactionsOption{
				WithDirection(types.TxDirectionOutgoing),
				WithAddress("*"),
				WithCount(0),
				WithSkip(0),
			},
			expDir: types.TxDirectionOutgoing, expAddr: "*", expCount: 10, expSkip: 0,
		},
		{
			opts: []ListTransactionsOption{
				WithDirection(types.TxDirectionIncoming),
				WithAddress("addr1"),
				WithCount(5),
				WithSkip(2),
			},
			expDir: types.TxDirectionIncoming, expAddr: "addr1", expCount: 5, expSkip: 2,
		},
	}

	for _, tt := range tests {
		cfg := defaultListTransactionsConfig
		for _, opt := range tt.opts {
			opt(&cfg)
		}

		assert.Equal(t, tt.expDir, cfg.direction)
		assert.Equal(t, tt.expAddr, cfg.address)
		assert.Equal(t, tt.expCount, cfg.count)
		assert.Equal(t, tt.expSkip, cfg.skip)
	}
}

func TestAddTransaction(t *testing.T) {
	td := setup(t)

	t.Run("add non-existent transaction returns error", func(t *testing.T) {
		txID := td.RandHash()

		td.mockStorage.EXPECT().IsLegacy().Return(false)
		td.mockStorage.EXPECT().HasTransaction(txID.String()).Return(false)
		td.mockProvider.EXPECT().
			GetTransaction(txID.String()).
			Return(nil, block.Height(0), errors.New("not exists"))

		err := td.wallet.AddTransaction(txID)
		require.Error(t, err)
	})

	t.Run("transaction already exists in storage", func(t *testing.T) {
		txID := td.RandHash()

		td.mockStorage.EXPECT().IsLegacy().Return(false)
		td.mockStorage.EXPECT().HasTransaction(txID.String()).Return(true)

		err := td.wallet.AddTransaction(txID)
		require.ErrorIs(t, err, ErrTransactionExists)
	})

	t.Run("add new transaction successfully", func(t *testing.T) {
		trx := td.GenerateTestTransferTx()

		td.mockStorage.EXPECT().IsLegacy().Return(false)
		td.mockStorage.EXPECT().HasTransaction(trx.ID().String()).Return(false)
		td.mockProvider.EXPECT().
			GetTransaction(trx.ID().String()).
			Return(trx, block.Height(td.RandHeight()), nil)
		td.mockStorage.EXPECT().HasAddress(gomock.Any()).Return(true)
		td.mockStorage.EXPECT().InsertTransaction(gomock.Any()).Return(nil)

		err := td.wallet.AddTransaction(trx.ID())
		require.NoError(t, err)
	})
}
