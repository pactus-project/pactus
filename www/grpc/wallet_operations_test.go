package grpc_test

import (
	"errors"
	"testing"
	"time"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/wallet/types"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadAndUnloadWallet(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	loaded, err := client.LoadWallet(t.Context(), &pactus.LoadWalletRequest{WalletName: "test"})
	require.NoError(t, err)
	assert.Equal(t, "test", loaded.WalletName)

	unloaded, err := client.UnloadWallet(t.Context(), &pactus.UnloadWalletRequest{WalletName: "test"})
	require.NoError(t, err)
	assert.Equal(t, "test", unloaded.WalletName)
}

func TestListTransactions(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on listing transactions", func(t *testing.T) {
		td.FakeWalletMgr.EXPECT().
			ListTransactions("test", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errors.New("error on listing transactions"))

		res, err := client.ListTransactions(t.Context(),
			&pactus.ListTransactionsRequest{WalletName: "test"})
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("List transactions successfully", func(t *testing.T) {
		created := time.Unix(1_700_000_000, 0)
		updated := time.Unix(1_700_000_100, 0)
		info := &types.TransactionInfo{
			No:          7,
			TxID:        td.RandHash().String(),
			Sender:      td.RandAccAddress().String(),
			Receiver:    td.RandAccAddress().String(),
			Direction:   types.TxDirectionOutgoing,
			Amount:      amount.Amount(5e9),
			Fee:         amount.Amount(1e7),
			Memo:        "memo",
			Status:      types.TransactionStatusConfirmed,
			BlockHeight: 42,
			PayloadType: payload.TypeTransfer,
			Data:        []byte{1, 2, 3},
			Comment:     "comment",
			CreatedAt:   created,
			UpdatedAt:   updated,
		}

		td.FakeWalletMgr.EXPECT().
			ListTransactions("test", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return([]*types.TransactionInfo{info}, nil)

		res, err := client.ListTransactions(t.Context(),
			&pactus.ListTransactionsRequest{
				WalletName: "test",
				Address:    info.Sender,
				Direction:  pactus.TxDirection_TX_DIRECTION_OUTGOING,
				Count:      10,
				Skip:       0,
			})
		require.NoError(t, err)
		assert.Equal(t, "test", res.WalletName)
		require.Len(t, res.Txs, 1)

		got := res.Txs[0]
		assert.Equal(t, int64(7), got.No)
		assert.Equal(t, info.TxID, got.TxId)
		assert.Equal(t, info.Sender, got.Sender)
		assert.Equal(t, info.Receiver, got.Receiver)
		assert.Equal(t, pactus.TxDirection_TX_DIRECTION_OUTGOING, got.Direction)
		assert.Equal(t, int64(5e9), got.Amount)
		assert.Equal(t, int64(1e7), got.Fee)
		assert.Equal(t, "memo", got.Memo)
		assert.Equal(t, pactus.TransactionStatus_TRANSACTION_STATUS_CONFIRMED, got.Status)
		assert.Equal(t, uint32(42), got.BlockHeight)
		assert.Equal(t, pactus.PayloadType_PAYLOAD_TYPE_TRANSFER, got.PayloadType)
		assert.Equal(t, []byte{1, 2, 3}, got.Data)
		assert.Equal(t, "comment", got.Comment)
		assert.Equal(t, created.Unix(), got.CreatedAt)
		assert.Equal(t, updated.Unix(), got.UpdatedAt)
	})

	t.Run("No transaction", func(t *testing.T) {
		td.FakeWalletMgr.EXPECT().
			ListTransactions("test", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil)

		res, err := client.ListTransactions(t.Context(),
			&pactus.ListTransactionsRequest{WalletName: "test"})
		require.NoError(t, err)
		assert.Empty(t, res.Txs)
	})
}

func TestUpdatePassword(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on updating the password", func(t *testing.T) {
		td.FakeWalletMgr.EXPECT().
			UpdatePassword("test", "wrong", "new").
			Return(errors.New("invalid password"))

		res, err := client.UpdatePassword(t.Context(),
			&pactus.UpdatePasswordRequest{WalletName: "test", OldPassword: "wrong", NewPassword: "new"})
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Update the password successfully", func(t *testing.T) {
		td.FakeWalletMgr.EXPECT().
			UpdatePassword("test", "old", "new").
			Return(nil)

		res, err := client.UpdatePassword(t.Context(),
			&pactus.UpdatePasswordRequest{WalletName: "test", OldPassword: "old", NewPassword: "new"})
		require.NoError(t, err)
		assert.Equal(t, "test", res.WalletName)
	})
}

func TestSetDefaultFee(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on setting the default fee", func(t *testing.T) {
		td.FakeWalletMgr.EXPECT().
			SetDefaultFee("test", amount.Amount(1e7)).
			Return(errors.New("error on setting the default fee"))

		res, err := client.SetDefaultFee(t.Context(),
			&pactus.SetDefaultFeeRequest{WalletName: "test", Amount: 1e7})
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Set the default fee successfully", func(t *testing.T) {
		td.FakeWalletMgr.EXPECT().
			SetDefaultFee("test", amount.Amount(1e7)).
			Return(nil)

		res, err := client.SetDefaultFee(t.Context(),
			&pactus.SetDefaultFeeRequest{WalletName: "test", Amount: 1e7})
		require.NoError(t, err)
		assert.Equal(t, "test", res.WalletName)
	})
}

func TestGetMnemonic(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on getting the mnemonic", func(t *testing.T) {
		td.FakeWalletMgr.EXPECT().
			Mnemonic("test", "wrong").
			Return("", errors.New("invalid password"))

		res, err := client.GetMnemonic(t.Context(),
			&pactus.GetMnemonicRequest{WalletName: "test", Password: "wrong"})
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Get the mnemonic successfully", func(t *testing.T) {
		td.FakeWalletMgr.EXPECT().
			Mnemonic("test", "password").
			Return("seed phrase", nil)

		res, err := client.GetMnemonic(t.Context(),
			&pactus.GetMnemonicRequest{WalletName: "test", Password: "password"})
		require.NoError(t, err)
		assert.Equal(t, "seed phrase", res.Mnemonic)
	})
}

func TestGetPrivateKey(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	_, prv := td.RandBLSKeyPair()
	addr := td.RandAccAddress().String()

	t.Run("Error on getting the private key", func(t *testing.T) {
		td.FakeWalletMgr.EXPECT().
			PrivateKey("test", "wrong", addr).
			Return(nil, errors.New("invalid password"))

		res, err := client.GetPrivateKey(t.Context(),
			&pactus.GetPrivateKeyRequest{WalletName: "test", Password: "wrong", Address: addr})
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Get the private key successfully", func(t *testing.T) {
		td.FakeWalletMgr.EXPECT().
			PrivateKey("test", "password", addr).
			Return(prv, nil)

		res, err := client.GetPrivateKey(t.Context(),
			&pactus.GetPrivateKeyRequest{WalletName: "test", Password: "password", Address: addr})
		require.NoError(t, err)
		assert.Equal(t, prv.String(), res.PrivateKey)
	})
}
