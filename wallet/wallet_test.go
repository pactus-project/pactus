package wallet

import (
	"encoding/json"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/util"
)

var tWallet *Wallet
var tPassword string

func setup(t *testing.T) {
	tPassword := ""
	walletPath := util.TempFilePath()
	mnemonic := GenerateMnemonic()
	w, err := FromMnemonic(walletPath, mnemonic, tPassword, NetworkMainNet)
	assert.NoError(t, err)
	assert.False(t, w.IsEncrypted())
	assert.Equal(t, w.Path(), walletPath)
	assert.Equal(t, w.Name(), path.Base(walletPath))

	// create some test addresses
	_, err = w.MakeNewAddress("", "addr-1")
	assert.NoError(t, err)
	_, err = w.MakeNewAddress("", "addr-2")
	assert.NoError(t, err)
	assert.False(t, w.IsEncrypted())
	tWallet = w
}

func TestOpenWallet(t *testing.T) {
	setup(t)

	t.Run("Ok", func(t *testing.T) {
		assert.NoError(t, tWallet.Save())
		_, err := OpenWallet(tWallet.path)
		assert.NoError(t, err)
	})

	t.Run("Invalid wallet path", func(t *testing.T) {
		_, err := OpenWallet(util.TempFilePath())
		assert.Error(t, err)
	})

	t.Run("Invalid crc", func(t *testing.T) {
		tWallet.store.VaultCRC = 0
		bs, _ := json.Marshal(tWallet.store)
		assert.NoError(t, util.WriteFile(tWallet.path, bs))

		_, err := OpenWallet(tWallet.path)
		assert.ErrorIs(t, err, ErrInvalidCRC)
	})

	t.Run("Invalid json", func(t *testing.T) {
		assert.NoError(t, util.WriteFile(tWallet.path, []byte("invalid_json")))

		_, err := OpenWallet(tWallet.path)
		assert.Error(t, err)
	})
}

func TestRecoverWallet(t *testing.T) {
	setup(t)

	mnemonic, _ := tWallet.Mnemonic(tPassword)
	password := ""
	t.Run("Wallet exists", func(t *testing.T) {
		// Save the test wallet first then
		// try to recover a wallet at the same place
		assert.NoError(t, tWallet.Save())

		_, err := FromMnemonic(tWallet.path, mnemonic, password, 0)
		assert.ErrorIs(t, err, NewErrWalletExits(tWallet.path))
	})

	t.Run("Invalid mnemonic", func(t *testing.T) {
		_, err := FromMnemonic(util.TempFilePath(),
			"invali mnemonic phrase seed", password, 0)
		assert.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		path := util.TempFilePath()
		recovered, err := FromMnemonic(path, mnemonic, password, 0)
		assert.NoError(t, err)

		addr1, err := recovered.MakeNewAddress("", "addr-1")
		assert.NoError(t, err)

		assert.NoFileExists(t, path)
		assert.NoError(t, recovered.Save())

		assert.FileExists(t, path)
		assert.True(t, tWallet.Contains(addr1))
	})
}

func TestSaveWallet(t *testing.T) {
	setup(t)

	t.Run("Invalid path", func(t *testing.T) {
		tWallet.path = "/"
		assert.Error(t, tWallet.Save())
	})
}
func TestInvalidAddress(t *testing.T) {
	setup(t)

	addr := crypto.GenerateTestAddress().String()
	_, err := tWallet.PrivateKey(tPassword, addr)
	assert.Error(t, err)
}

func TestImportPrivateKey(t *testing.T) {
	setup(t)

	_, prv := bls.GenerateTestKeyPair()
	assert.NoError(t, tWallet.ImportPrivateKey(tPassword, prv.String()))

	addr := prv.PublicKey().Address().String()
	assert.True(t, tWallet.Contains(addr))
	pub, err := tWallet.PublicKey(tPassword, addr)
	assert.NoError(t, err)
	assert.Equal(t, pub, prv.PublicKey().String())
}
func TestTestKeyInfo(t *testing.T) {
	mnemonic := GenerateMnemonic()
	w1, err := FromMnemonic(util.TempFilePath(), mnemonic, tPassword,
		NetworkMainNet)
	assert.NoError(t, err)
	addrStr1, _ := w1.MakeNewAddress("", "")
	prvStr1, _ := w1.PrivateKey("", addrStr1)
	prv1, _ := bls.PrivateKeyFromString(prvStr1)

	w2, err := FromMnemonic(util.TempFilePath(), mnemonic, tPassword,
		NetworkTestNet)
	assert.NoError(t, err)
	addrStr2, _ := w2.MakeNewAddress("", "")
	prvStr2, _ := w2.PrivateKey("", addrStr2)
	prv2, _ := bls.PrivateKeyFromString(prvStr2)

	assert.NotEqual(t, prv1.Bytes(), prv2.Bytes(),
		"Should generate different private key for the testnet")
}

func TestMakeSendTx(t *testing.T) {
	setup(t)

	//TODO
}
