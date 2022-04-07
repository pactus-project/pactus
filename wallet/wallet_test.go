package wallet

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/util"
)

var tWallet *Wallet
var tPassphrase string

func setup(t *testing.T) {
	passphrase := ""
	path := util.TempFilePath()
	w, err := CreateWallet(path, passphrase, 0) // 2 for testing
	assert.NoError(t, err)
	assert.False(t, w.IsEncrypted())
	assert.Equal(t, w.Path(), path)

	// create some test addresses
	_, err = w.NewAddress(passphrase, "addr-1")
	assert.NoError(t, err)
	_, err = w.NewAddress(passphrase, "addr-2")
	assert.NoError(t, err)
	_, err = w.NewAddress(passphrase, "addr-3")
	assert.NoError(t, err)

	tPassphrase = passphrase
	tWallet = w
}

func reopenWallet(t *testing.T) {
	w, err := OpenWallet(tWallet.path)
	assert.NoError(t, err)
	assert.Equal(t, tWallet.store.UUID, w.store.UUID, "UUID is changed")
	tWallet = w
}

func TestCreateWallet(t *testing.T) {
	setup(t)

	t.Run("Wallet exists", func(t *testing.T) {
		_, err := CreateWallet(tWallet.path, "", 0)
		assert.Error(t, err)
	})

	t.Run("Invalid network", func(t *testing.T) {
		_, err := CreateWallet(util.TempFilePath(), "", 3)
		assert.Error(t, err)
	})

	t.Run("OK", func(t *testing.T) {
		w, err := CreateWallet(util.TempFilePath(), "super_secret_password", 0)
		assert.NoError(t, err)
		assert.True(t, w.IsEncrypted())
	})
}

func TestOpenWallet(t *testing.T) {
	if os.Getenv("INVALID_WALLET") == "1" {
		w, _ := OpenWallet(os.Getenv("WALLET_PATH"))
		assert.Nil(t, w, "should exit before")
	}

	setup(t)
	t.Run("Invalid wallet path", func(t *testing.T) {
		_, err := OpenWallet(util.TempFilePath())
		assert.Error(t, err)
	})

	t.Run("Invalid crc", func(t *testing.T) {
		tWallet.store.VaultCRC = 0
		bs, _ := json.Marshal(tWallet.store)
		assert.NoError(t, util.WriteFile(tWallet.path, bs))

		cmd := exec.Command(os.Args[0], "-test.run=TestOpenWallet", "-covermode=atomic") //nolint:gosec
		cmd.Env = append(os.Environ(), "INVALID_WALLET=1")
		cmd.Env = append(cmd.Env, "WALLET_PATH="+tWallet.path)
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	})

	t.Run("Invalid json", func(t *testing.T) {
		assert.NoError(t, util.WriteFile(tWallet.path, []byte("invalid_json")))

		cmd := exec.Command(os.Args[0], "-test.run=TestOpenWallet", "-covermode=atomic") //nolint:gosec
		cmd.Env = append(os.Environ(), "INVALID_WALLET=1")
		cmd.Env = append(cmd.Env, "WALLET_PATH="+tWallet.path)
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	})

}

func TestRecoverWallet(t *testing.T) {
	setup(t)

	mnemonic, _ := tWallet.Mnemonic(tPassphrase)
	t.Run("Wallet exists", func(t *testing.T) {
		_, err := RecoverWallet(tWallet.path, mnemonic, 0)
		assert.Error(t, err)
	})

	t.Run("Invalid mnemonic", func(t *testing.T) {
		_, err := RecoverWallet(util.TempFilePath(), "invali mnemonic phrase seed", 0)
		assert.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		recovered, err := RecoverWallet(util.TempFilePath(), mnemonic, 0)
		assert.NoError(t, err)

		reopenWallet(t)
		_, err = recovered.NewAddress("", "addr-1")
		assert.NoError(t, err)
		_, err = recovered.NewAddress("", "addr-2")
		assert.NoError(t, err)
		_, err = recovered.NewAddress("", "addr-3")
		assert.NoError(t, err)

		assert.Equal(t, tWallet.Addresses(), recovered.Addresses())
		assert.Equal(t, tWallet.store.VaultCRC, recovered.store.VaultCRC)
	})
}

func TestGetPrivateKey(t *testing.T) {
	setup(t)

	addrs := tWallet.Addresses()
	assert.NotEmpty(t, addrs)
	for addr := range addrs {
		prvStr, err := tWallet.PrivateKey(tPassphrase, addr)
		assert.NoError(t, err)
		pubStr, err := tWallet.PublicKey(tPassphrase, addr)
		assert.NoError(t, err)
		prv, _ := bls.PrivateKeyFromString(prvStr)
		pub, _ := bls.PublicKeyFromString(pubStr)
		assert.True(t, prv.PublicKey().EqualsTo(pub))
		assert.Equal(t, pub.Address().String(), addr)
	}
}

func TestInvalidAddress(t *testing.T) {
	setup(t)

	_, err := tWallet.PrivateKey(tPassphrase, crypto.GenerateTestAddress().String())
	assert.Error(t, err)
}

func TestImportPrivateKey(t *testing.T) {
	setup(t)

	_, prv1 := bls.GenerateTestKeyPair()
	assert.NoError(t, tWallet.ImportPrivateKey(tPassphrase, prv1.String()))
	reopenWallet(t)

	assert.True(t, tWallet.store.Contains(prv1.PublicKey().Address()))
	prv2, err := tWallet.PrivateKey(tPassphrase, prv1.PublicKey().Address().String())
	assert.NoError(t, err)
	assert.Equal(t, prv1.String(), prv2)

	// Import again
	assert.Error(t, tWallet.ImportPrivateKey(tPassphrase, prv1.String()))
}
