package grpc

import (
	"context"
	"testing"

	"github.com/pactus-project/pactus/wallet"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
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
