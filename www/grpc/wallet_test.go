package grpc

import (
	"testing"

	"github.com/pactus-project/pactus/wallet"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestGenerateMnemonic(t *testing.T) {
	conn, client := callWalletSerer(t)

	t.Run("Should return mnemonic", func(t *testing.T) {
		res, err := client.GenerateMnemonic(tCtx, &pactus.GenerateMnemonicRequest{Language: "english", Entropy: 128})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Mnemonic)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}
func TestCreateWallet(t *testing.T) {
	conn, client := callWalletSerer(t)

	t.Run("Invalid mnemonic", func(t *testing.T) {
		res, err := client.CreateWallet(tCtx, &pactus.CreateWalletRequest{
			Name:     "test",
			Mnemonic: "not valid",
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("No name, should return an error", func(t *testing.T) {
		mnemonic := wallet.GenerateMnemonic(128)
		res, err := client.CreateWallet(tCtx, &pactus.CreateWalletRequest{
			Name:     "",
			Mnemonic: mnemonic,
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Bad name, should return an error", func(t *testing.T) {
		mnemonic := wallet.GenerateMnemonic(128)
		res, err := client.CreateWallet(tCtx, &pactus.CreateWalletRequest{
			Name:     "..",
			Mnemonic: mnemonic,
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should create wallet", func(t *testing.T) {
		mnemonic := wallet.GenerateMnemonic(128)
		res, err := client.CreateWallet(tCtx, &pactus.CreateWalletRequest{
			Name:     "test",
			Mnemonic: mnemonic,
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
}
