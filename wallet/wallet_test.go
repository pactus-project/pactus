package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

var tWallet *Wallet
var tPassphrase string

func setup(t *testing.T) {
	tPassphrase = "super_secret_password"
	w, err := NewWallet(util.TempFilePath(), tPassphrase)
	assert.NoError(t, err)

	tWallet = w
}
func TestOpenWallet(t *testing.T) {
	setup(t)

	_, err := OpenWallet(tWallet.path)
	assert.NoError(t, err)
}

func TestRecoverWallet(t *testing.T) {
	setup(t)

	mnemonic := tWallet.Mnemonic(tPassphrase)
	recovered, err := RecoverWallet(util.TempFilePath(), mnemonic)
	assert.NoError(t, err)

	assert.Equal(t, tWallet.store.Vault.Seed.parentKey(tPassphrase).Bytes(), recovered.store.Vault.Seed.parentKey("").Bytes())
}

func TestGetPrivateKey(t *testing.T) {
	setup(t)

	addrs := tWallet.Addresses()
	for _, addr := range addrs {
		prv, err := tWallet.PrivateKey(tPassphrase, addr.String())
		assert.NoError(t, err)
		assert.Equal(t, prv.PublicKey().Address().String(), addr.String())
	}
}

func TestInvalidAddress(t *testing.T) {
	setup(t)

	_, err := tWallet.PrivateKey(tPassphrase, crypto.GenerateTestAddress().String())
	assert.Error(t, err)
}
