package grpc

import (
	"errors"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/types"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestWalletServiceIsDisabled(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = false

	td := setup(t, conf)
	client := td.walletClient(t)

	res, err := client.CreateWallet(t.Context(),
		&pactus.CreateWalletRequest{
			WalletName: "TestWallet",
		})
	assert.ErrorIs(t, err, status.Error(codes.Unimplemented, "unknown service pactus.Wallet"))
	assert.Nil(t, res)
}

func TestCreateWallet(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("No wallet name, should return an error", func(t *testing.T) {
		res, err := client.CreateWallet(t.Context(),
			&pactus.CreateWalletRequest{
				WalletName: "",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Error on creating wallet", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			CreateWallet("test", "password").
			Return("", errors.New("error on creating wallet"))

		res, err := client.CreateWallet(t.Context(),
			&pactus.CreateWalletRequest{
				WalletName: "test",
				Password:   "password",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Create wallet successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			CreateWallet("test", "password").
			Return("mnemonic", nil)

		res, err := client.CreateWallet(t.Context(),
			&pactus.CreateWalletRequest{
				WalletName: "test",
				Password:   "password",
			})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "test", res.WalletName)
		assert.Equal(t, "mnemonic", res.Mnemonic)
	})
}

func TestRestoreWallet(t *testing.T) {
	config := testConfig()
	config.EnableWallet = true

	td := setup(t, config)
	client := td.walletClient(t)

	t.Run("No wallet name, should return an error", func(t *testing.T) {
		res, err := client.RestoreWallet(t.Context(),
			&pactus.RestoreWalletRequest{})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("No mnemonic, should return an error", func(t *testing.T) {
		res, err := client.RestoreWallet(t.Context(),
			&pactus.RestoreWalletRequest{
				WalletName: "test",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Restore wallet successfully", func(t *testing.T) {
		mnemonic, err := wallet.GenerateMnemonic(128)
		assert.NoError(t, err)

		td.mockWalletMgr.EXPECT().
			RestoreWallet("test", mnemonic, "password").
			Return(nil)

		res, err := client.RestoreWallet(t.Context(),
			&pactus.RestoreWalletRequest{
				WalletName: "test",
				Mnemonic:   mnemonic,
				Password:   "password",
			})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "test", res.WalletName)
	})
}

func TestLoadWallet(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on loading wallet", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			LoadWallet("test", gomock.Any()).
			Return(errors.New("error on loading wallet"))

		res, err := client.LoadWallet(t.Context(),
			&pactus.LoadWalletRequest{
				WalletName: "test",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Load wallet successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			LoadWallet("test", gomock.Any()).
			Return(nil)

		res, err := client.LoadWallet(t.Context(),
			&pactus.LoadWalletRequest{
				WalletName: "test",
			})
		require.NoError(t, err)
		assert.Equal(t, "test", res.WalletName)
	})
}

func TestUnloadWallet(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on unloading wallet", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			UnloadWallet("test").
			Return(errors.New("error on unloading wallet"))

		res, err := client.UnloadWallet(t.Context(),
			&pactus.UnloadWalletRequest{
				WalletName: "test",
			})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Unload wallet successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			UnloadWallet("test").
			Return(nil)

		res, err := client.UnloadWallet(t.Context(),
			&pactus.UnloadWalletRequest{
				WalletName: "test",
			})
		require.NoError(t, err)
		assert.Equal(t, "test", res.WalletName)
	})
}

func TestGetTotalBalance(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on getting total balance", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			TotalBalance("test").
			Return(amount.Amount(0), errors.New("error on getting total balance"))

		res, err := client.GetTotalBalance(t.Context(),
			&pactus.GetTotalBalanceRequest{
				WalletName: "test",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Get total balance successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			TotalBalance("test").
			Return(amount.Amount(123), nil)

		res, err := client.GetTotalBalance(t.Context(),
			&pactus.GetTotalBalanceRequest{
				WalletName: "test",
			})
		require.NoError(t, err)
		assert.Equal(t, "test", res.WalletName)
		assert.Equal(t, int64(123), res.TotalBalance)
	})
}

func TestGetTotalStake(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on getting total stake", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			TotalStake("test").
			Return(amount.Amount(0), errors.New("error on getting total stake"))

		res, err := client.GetTotalStake(t.Context(),
			&pactus.GetTotalStakeRequest{
				WalletName: "test",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Get total stake successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			TotalStake("test").
			Return(amount.Amount(123), nil)

		res, err := client.GetTotalStake(t.Context(),
			&pactus.GetTotalStakeRequest{
				WalletName: "test",
			})
		require.NoError(t, err)
		assert.Equal(t, "test", res.WalletName)
		assert.Equal(t, int64(123), res.TotalStake)
	})
}

func TestSignRawTransaction(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Bad traansaction data", func(t *testing.T) {
		res, err := client.SignRawTransaction(t.Context(),
			&pactus.SignRawTransactionRequest{
				WalletName:     "test",
				RawTransaction: "invalid-hex",
			})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Error on signing raw transaction", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			SignRawTransaction("test", "", gomock.Any()).
			Return(nil, nil, errors.New("error on signing raw transaction"))

		res, err := client.SignRawTransaction(t.Context(),
			&pactus.SignRawTransactionRequest{
				WalletName:     "test",
				RawTransaction: "1a2b3c4d",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Sign raw transaction successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			SignRawTransaction("test", "", gomock.Any()).
			Return(nil, nil, nil)

		res, err := client.SignRawTransaction(t.Context(),
			&pactus.SignRawTransactionRequest{
				WalletName:     "test",
				RawTransaction: "1a2b3c4d",
			})
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestGetValidatorAddress(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on getting validator address", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			GetValidatorAddress("pubKey").
			Return("", errors.New("error on getting validator address"))

		res, err := client.GetValidatorAddress(t.Context(),
			&pactus.GetValidatorAddressRequest{
				PublicKey: "pubKey",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Get validator address successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			GetValidatorAddress("pubKey").
			Return("valAddr", nil)

		res, err := client.GetValidatorAddress(t.Context(),
			&pactus.GetValidatorAddressRequest{
				PublicKey: "pubKey",
			})
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "valAddr", res.Address)
	})
}

func TestSignMessage(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on signing message", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			SignMessage("test", "password", "addr", "hello").
			Return("", errors.New("error on signing message"))

		res, err := client.SignMessage(t.Context(),
			&pactus.SignMessageRequest{
				WalletName: "test",
				Password:   "password",
				Address:    "addr",
				Message:    "hello",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Sign message successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			SignMessage("test", "password", "addr", "hello").
			Return("signature", nil)

		res, err := client.SignMessage(t.Context(),
			&pactus.SignMessageRequest{
				WalletName: "test",
				Password:   "password",
				Address:    "addr",
				Message:    "hello",
			})
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "signature", res.Signature)
	})
}

func TestNewAddress(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on getting new address", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			NewAddress("test", "label", "password", crypto.AddressTypeBLSAccount).
			Return(nil, errors.New("error on getting new address"))

		res, err := client.GetNewAddress(t.Context(),
			&pactus.GetNewAddressRequest{
				WalletName:  "test",
				AddressType: pactus.AddressType_ADDRESS_TYPE_BLS_ACCOUNT,
				Label:       "label",
				Password:    "password",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Get new address successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			NewAddress("test", "label", "password", crypto.AddressTypeBLSAccount).
			Return(&types.AddressInfo{
				Address:   "addr",
				Label:     "label",
				PublicKey: "pub",
				Path:      "m/44'/0'/0'/0/0",
			}, nil)

		res, err := client.GetNewAddress(t.Context(),
			&pactus.GetNewAddressRequest{
				WalletName:  "test",
				AddressType: pactus.AddressType_ADDRESS_TYPE_BLS_ACCOUNT,
				Label:       "label",
				Password:    "password",
			})
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "test", res.WalletName)
		require.NotNil(t, res.AddressInfo)
		assert.Equal(t, "addr", res.AddressInfo.Address)
		assert.Equal(t, "label", res.AddressInfo.Label)
		assert.Equal(t, "pub", res.AddressInfo.PublicKey)
		assert.Equal(t, "m/44'/0'/0'/0/0", res.AddressInfo.Path)
	})
}

func TestAddressInfo(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on getting address info", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			AddressInfo("test", "addr").
			Return(nil, errors.New("error on getting address info"))

		res, err := client.GetAddressInfo(t.Context(),
			&pactus.GetAddressInfoRequest{
				WalletName: "test",
				Address:    "addr",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Get address info successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			AddressInfo("test", "addr").
			Return(&types.AddressInfo{
				Address:   "addr",
				Label:     "label",
				PublicKey: "pub",
				Path:      "path",
			}, nil)

		res, err := client.GetAddressInfo(t.Context(),
			&pactus.GetAddressInfoRequest{
				WalletName: "test",
				Address:    "addr",
			})
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "test", res.WalletName)
		assert.Equal(t, "addr", res.Address)
		assert.Equal(t, "label", res.Label)
		assert.Equal(t, "pub", res.PublicKey)
		assert.Equal(t, "path", res.Path)
	})
}

func TestSetAddressLabel(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on setting address label", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			SetAddressLabel("test", "addr", "label").
			Return(errors.New("error on setting address label"))

		res, err := client.SetAddressLabel(t.Context(),
			&pactus.SetAddressLabelRequest{
				WalletName: "test",
				Password:   "password",
				Address:    "addr",
				Label:      "label",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Set address label successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			SetAddressLabel("test", "addr", "label").
			Return(nil)

		res, err := client.SetAddressLabel(t.Context(),
			&pactus.SetAddressLabelRequest{
				WalletName: "test",
				Password:   "password",
				Address:    "addr",
				Label:      "label",
			})
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "test", res.WalletName)
		assert.Equal(t, "addr", res.Address)
		assert.Equal(t, "label", res.Label)
	})
}

func TestListWallet(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on listing wallets", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			ListWallets().
			Return(nil, errors.New("error on listing wallets"))

		res, err := client.ListWallets(t.Context(), &pactus.ListWalletsRequest{})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("List wallets successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			ListWallets().
			Return([]string{"w1", "w2"}, nil)

		res, err := client.ListWallets(t.Context(), &pactus.ListWalletsRequest{})
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, []string{"w1", "w2"}, res.Wallets)
	})
}

func TestGetWalletInfo(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on getting wallet info", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			WalletInfo("test").
			Return(nil, errors.New("error on getting wallet info"))

		res, err := client.GetWalletInfo(t.Context(),
			&pactus.GetWalletInfoRequest{WalletName: "test"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Get wallet info successfully", func(t *testing.T) {
		createdAt := time.Unix(123, 0).UTC()
		td.mockWalletMgr.EXPECT().
			WalletInfo("test").
			Return(&types.WalletInfo{
				Version:    7,
				Network:    "testnet",
				Encrypted:  true,
				UUID:       "uuid",
				CreatedAt:  createdAt,
				DefaultFee: amount.Amount(456),
			}, nil)

		res, err := client.GetWalletInfo(t.Context(),
			&pactus.GetWalletInfoRequest{WalletName: "test"})
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "test", res.WalletName)
		assert.Equal(t, int32(7), res.Version)
		assert.Equal(t, "testnet", res.Network)
		assert.True(t, res.Encrypted)
		assert.Equal(t, "uuid", res.Uuid)
		assert.Equal(t, createdAt.Unix(), res.CreatedAt)
		assert.Equal(t, int64(456), res.DefaultFee)
	})
}

func TestListAddress(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	client := td.walletClient(t)

	t.Run("Error on listing addresses", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			ListAddresses("test").
			Return(nil, errors.New("error on listing addresses"))

		res, err := client.ListAddresses(t.Context(),
			&pactus.ListAddressesRequest{WalletName: "test"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("List addresses successfully", func(t *testing.T) {
		td.mockWalletMgr.EXPECT().
			ListAddresses("test").
			Return([]types.AddressInfo{
				{
					Address:   "addr1",
					Label:     "label1",
					PublicKey: "pub1",
					Path:      "path1",
				},
				{
					Address:   "addr2",
					Label:     "label2",
					PublicKey: "pub2",
					Path:      "path2",
				},
			}, nil)

		res, err := client.ListAddresses(t.Context(),
			&pactus.ListAddressesRequest{WalletName: "test"})
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "test", res.WalletName)
		require.Len(t, res.Data, 2)
		assert.Equal(t, "addr1", res.Data[0].Address)
		assert.Equal(t, "label1", res.Data[0].Label)
		assert.Equal(t, "pub1", res.Data[0].PublicKey)
		assert.Equal(t, "path1", res.Data[0].Path)
	})
}
