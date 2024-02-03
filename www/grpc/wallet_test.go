package grpc

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDisableWallet(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = false

	td := setup(t, conf)
	conn, client := td.walletClient(t)

	mnemonic, _ := wallet.GenerateMnemonic(128)
	res, err := client.CreateWallet(context.Background(),
		&pactus.CreateWalletRequest{
			WalletName: "TestWallet",
			Mnemonic:   mnemonic,
		})
	assert.ErrorIs(t, err, status.Error(codes.Unimplemented, "unknown service pactus.Wallet"))
	assert.Nil(t, res)

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestCreateWallet(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	conn, client := td.walletClient(t)

	t.Run("Invalid mnemonic", func(t *testing.T) {
		res, err := client.CreateWallet(context.Background(),
			&pactus.CreateWalletRequest{
				WalletName: "test",
				Mnemonic:   "not valid",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("No name, should return an error", func(t *testing.T) {
		mnemonic, _ := wallet.GenerateMnemonic(128)
		res, err := client.CreateWallet(context.Background(),
			&pactus.CreateWalletRequest{
				WalletName: "",
				Mnemonic:   mnemonic,
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Bad name, should return an error", func(t *testing.T) {
		mnemonic, _ := wallet.GenerateMnemonic(128)
		res, err := client.CreateWallet(context.Background(),
			&pactus.CreateWalletRequest{
				WalletName: "..",
				Mnemonic:   mnemonic,
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should create wallet", func(t *testing.T) {
		mnemonic, _ := wallet.GenerateMnemonic(128)
		res, err := client.CreateWallet(context.Background(),
			&pactus.CreateWalletRequest{
				WalletName: "test",
				Mnemonic:   mnemonic,
			})
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestLoadWallet(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	conn, client := td.walletClient(t)

	wltName := td.RandString(16)
	wltPath := filepath.Join(conf.WalletsDir, wltName)
	mnemonic, _ := wallet.GenerateMnemonic(128)
	wlt, err := wallet.Create(wltPath, mnemonic, "", genesis.Mainnet)
	require.NoError(t, err)

	wltAddr, err := wlt.NewBLSAccountAddress("test")
	assert.NoError(t, err)
	require.NoError(t, wlt.Save())

	t.Run("Load non-existing wallet", func(t *testing.T) {
		res, err := client.LoadWallet(context.Background(),
			&pactus.LoadWalletRequest{
				WalletName: "non-existing",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Load existing wallet", func(t *testing.T) {
		res, err := client.LoadWallet(context.Background(),
			&pactus.LoadWalletRequest{
				WalletName: wltName,
			})
		require.NoError(t, err)
		assert.Equal(t, wltName, res.WalletName)
	})

	t.Run("Load wallet again", func(t *testing.T) {
		res, err := client.LoadWallet(context.Background(),
			&pactus.LoadWalletRequest{
				WalletName: wltName,
			})
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Unload unknown wallet", func(t *testing.T) {
		res, err := client.UnloadWallet(context.Background(),
			&pactus.UnloadWalletRequest{
				WalletName: "not-loade",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Sign raw transaction, OK", func(t *testing.T) {
		wltAddr, _ := crypto.AddressFromString(wltAddr)
		bondTx := tx.NewBondTx(td.RandHeight(), wltAddr, td.RandValAddress(), nil, td.RandAmount(), td.RandAmount(), "memo")
		rawTx, _ := bondTx.Bytes()
		res, err := client.SignRawTransaction(context.Background(),
			&pactus.SignRawTransactionRequest{
				WalletName:     wltName,
				RawTransaction: rawTx,
				Password:       "",
			})
		assert.NoError(t, err)
		assert.Equal(t, bondTx.ID().Bytes(), res.TransactionId)

		signedTx, err := tx.FromBytes(res.SignedRawTransaction)
		assert.NoError(t, err)
		assert.NotNil(t, signedTx.Signature())
		assert.Nil(t, signedTx.BasicCheck())
	})

	t.Run("Sign raw transaction using not loaded wallet", func(t *testing.T) {
		wltAddr, _ := crypto.AddressFromString(wltAddr)
		bondTx := tx.NewBondTx(td.RandHeight(), wltAddr, td.RandValAddress(), nil, td.RandAmount(), td.RandAmount(), "memo")
		rawTx, _ := bondTx.Bytes()
		res, err := client.SignRawTransaction(context.Background(),
			&pactus.SignRawTransactionRequest{
				WalletName:     "not-loaded-wallet",
				RawTransaction: rawTx,
				Password:       "",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Sign invalid raw transaction", func(t *testing.T) {
		invRawData := td.DecodingHex("bad0")
		res, err := client.SignRawTransaction(context.Background(),
			&pactus.SignRawTransactionRequest{
				WalletName:     wltName,
				RawTransaction: invRawData,
				Password:       "",
			})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Unload wallet", func(t *testing.T) {
		res, err := client.UnloadWallet(context.Background(),
			&pactus.UnloadWalletRequest{
				WalletName: wltName,
			})
		require.NoError(t, err)
		assert.Equal(t, wltName, res.WalletName)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestGetValidatorAddress(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true

	td := setup(t, conf)
	conn, client := td.walletClient(t)

	t.Run("Invalid public key", func(t *testing.T) {
		res, err := client.GetValidatorAddress(context.Background(),
			&pactus.GetValidatorAddressRequest{PublicKey: "something"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("OK", func(t *testing.T) {
		valKey := td.RandValKey()
		pubKey := valKey.PublicKey()

		res, err := client.GetValidatorAddress(context.Background(),
			&pactus.GetValidatorAddressRequest{PublicKey: pubKey.String()})

		assert.Nil(t, err)
		assert.Equal(t, pubKey.ValidatorAddress().String(), res.Address)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}
