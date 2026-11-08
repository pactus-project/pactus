package wallet

import (
	"errors"
	"testing"

	"github.com/pactus-project/pactus/types"
	wtypes "github.com/pactus-project/pactus/wallet/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTransactionsListOptions(t *testing.T) {
	tests := []struct {
		opts     []ListTransactionsOption
		expDir   wtypes.TxDirection
		expAddr  string
		expCount int
		expSkip  int
	}{
		{
			opts:   []ListTransactionsOption{},
			expDir: wtypes.TxDirectionAny, expAddr: "*", expCount: 10, expSkip: 0,
		},
		{
			opts: []ListTransactionsOption{
				WithDirection(wtypes.TxDirectionAny),
				WithAddress(""),
				WithCount(-1),
				WithSkip(-1),
			},
			expDir: wtypes.TxDirectionAny, expAddr: "*", expCount: 10, expSkip: 0,
		},
		{
			opts: []ListTransactionsOption{
				WithDirection(wtypes.TxDirectionOutgoing),
				WithAddress("*"),
				WithCount(0),
				WithSkip(0),
			},
			expDir: wtypes.TxDirectionOutgoing, expAddr: "*", expCount: 10, expSkip: 0,
		},
		{
			opts: []ListTransactionsOption{
				WithDirection(wtypes.TxDirectionIncoming),
				WithAddress("addr1"),
				WithCount(5),
				WithSkip(2),
			},
			expDir: wtypes.TxDirectionIncoming, expAddr: "addr1", expCount: 5, expSkip: 2,
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
			Return(nil, types.Height(0), errors.New("not exists"))

		err := td.wallet.AddTransactionByID(txID)
		require.Error(t, err)
	})

	t.Run("transaction already exists in storage", func(t *testing.T) {
		txID := td.RandHash()

		td.mockStorage.EXPECT().IsLegacy().Return(false)
		td.mockStorage.EXPECT().HasTransaction(txID.String()).Return(true)

		err := td.wallet.AddTransactionByID(txID)
		require.ErrorIs(t, err, ErrTransactionExists)
	})

	t.Run("add new transaction successfully", func(t *testing.T) {
		trx := td.GenerateTestTransferTx()

		td.mockStorage.EXPECT().IsLegacy().Return(false)
		td.mockStorage.EXPECT().HasTransaction(trx.ID().String()).Return(false)
		td.mockProvider.EXPECT().
			GetTransaction(trx.ID().String()).
			Return(trx, td.RandHeight(), nil)
		td.mockStorage.EXPECT().HasAddress(gomock.Any()).Return(true)
		td.mockStorage.EXPECT().InsertTransaction(gomock.Any()).Return(nil)

		err := td.wallet.AddTransactionByID(trx.ID())
		require.NoError(t, err)
	})
}

func TestProcessBlockPendingTransactions(t *testing.T) {
	td := setup(t)

	trx := td.GenerateTestTransferTx()
	data, _ := trx.Bytes()
	pendingInfo := &wtypes.TransactionInfo{No: 1, TxID: trx.ID().String(), Data: data}
	pendingTxs := map[string]*wtypes.TransactionInfo{pendingInfo.TxID: pendingInfo}
	blk, _ := td.GenerateTestBlock(td.RandHeight())

	t.Run("committed transaction is confirmed, not checked again", func(t *testing.T) {
		height := td.RandHeight()

		td.mockStorage.EXPECT().IsLegacy().Return(false)
		td.mockStorage.EXPECT().GetPendingTransactions().Return(pendingTxs, nil)
		td.mockProvider.EXPECT().GetTransaction(pendingInfo.TxID).Return(trx, height, nil)
		td.mockStorage.EXPECT().
			UpdateTransactionStatus(pendingInfo.No, wtypes.TransactionStatusConfirmed, height).
			Return(nil)
		td.mockStorage.EXPECT().HasAddress(gomock.Any()).Return(false).AnyTimes()

		td.wallet.processEvent(blk)
	})

	t.Run("uncommitted transaction is checked and re-broadcast", func(t *testing.T) {
		td.mockStorage.EXPECT().IsLegacy().Return(false)
		td.mockStorage.EXPECT().GetPendingTransactions().Return(pendingTxs, nil)
		td.mockProvider.EXPECT().
			GetTransaction(pendingInfo.TxID).
			Return(nil, types.Height(0), errors.New("not found"))
		td.mockProvider.EXPECT().CheckTransaction(data).Return(nil)
		td.mockProvider.EXPECT().SendTx(gomock.Any()).Return(pendingInfo.TxID, nil)
		td.mockStorage.EXPECT().HasAddress(gomock.Any()).Return(false).AnyTimes()

		td.wallet.processEvent(blk)
	})
}
