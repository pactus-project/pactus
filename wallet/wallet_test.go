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

func TestGetPrivateKey(t *testing.T) {
	setup(t)

	addrs := tWallet.Addresses(tPassphrase)
	prv, err := tWallet.PrivateKey(tPassphrase, addrs[0].String())
	assert.NoError(t, err)
	assert.Equal(t, prv.PublicKey().Address().String(), addrs[0].String())

	_, err = tWallet.PrivateKey(tPassphrase, crypto.GenerateTestAddress().String())
	assert.Error(t, err)
}
