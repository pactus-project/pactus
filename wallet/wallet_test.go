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
	passphrase := "super_secret_password"
	path := util.TempFilePath()
	mnemonic := GenerateMnemonic()
	w, err := FromMnemonic(path, mnemonic, "", 0)
	assert.NoError(t, err)
	assert.False(t, w.IsEncrypted())
	assert.Equal(t, w.Path(), path)

	// create some test addresses
	_, err = w.NewAddress("", "addr-1")
	assert.NoError(t, err)
	_, err = w.NewAddress("", "addr-2")
	assert.NoError(t, err)
	_, err = w.NewAddress("", "addr-3")
	assert.NoError(t, err)

	// Create some keys
	_, prv1 := bls.GenerateTestKeyPair()
	_, prv2 := bls.GenerateTestKeyPair()

	assert.NoError(t, w.ImportPrivateKey("", prv1.String()))
	assert.NoError(t, w.ImportPrivateKey("", prv2.String()))

	assert.NoError(t, w.UpdatePassword("", passphrase))
	assert.True(t, w.IsEncrypted())
	tPassphrase = passphrase
	tWallet = w
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
	password := ""
	t.Run("Wallet exists", func(t *testing.T) {
		// Save the test wallet first then
		// try to recover a wallet at the same place
		assert.NoError(t, tWallet.Save())

		_, err := FromMnemonic(tWallet.path, mnemonic, password, 0)
		assert.Error(t, err)
	})

	t.Run("Invalid mnemonic", func(t *testing.T) {
		_, err := FromMnemonic(util.TempFilePath(), "invali mnemonic phrase seed", password, 0)
		assert.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		path := util.TempFilePath()
		recovered, err := FromMnemonic(path, mnemonic, password, 0)
		assert.NoError(t, err)

		addr1, err := recovered.NewAddress("", "addr-1")
		assert.NoError(t, err)
		addr2, err := recovered.NewAddress("", "addr-2")
		assert.NoError(t, err)
		addr3, err := recovered.NewAddress("", "addr-3")
		assert.NoError(t, err)

		assert.NoFileExists(t, path)
		assert.NoError(t, recovered.Save())

		assert.FileExists(t, path)
		assert.True(t, tWallet.Contains(addr1))
		assert.True(t, tWallet.Contains(addr2))
		assert.True(t, tWallet.Contains(addr3))
	})
}

func TestGetPrivateKey(t *testing.T) {
	setup(t)

	_, prv := bls.GenerateTestKeyPair()
	assert.NoError(t, tWallet.ImportPrivateKey(tPassphrase, prv.String()))

	t.Run("Check all private keys", func(t *testing.T) {
		addrs := tWallet.Addresses()
		assert.NotEmpty(t, addrs)
		for _, addr := range addrs {
			prvStr, err := tWallet.PrivateKey(tPassphrase, addr.Address)
			assert.NoError(t, err)
			pubStr, err := tWallet.PublicKey(tPassphrase, addr.Address)
			assert.NoError(t, err)
			prv, _ := bls.PrivateKeyFromString(prvStr)
			pub, _ := bls.PublicKeyFromString(pubStr)
			assert.True(t, prv.PublicKey().EqualsTo(pub))
			assert.Equal(t, pub.Address().String(), addr.Address)
		}
	})

	t.Run("Empty password", func(t *testing.T) {
		addrs := tWallet.Addresses()
		assert.NotEmpty(t, addrs)
		for _, addr := range addrs {
			_, err := tWallet.PrivateKey("", addr.Address)
			assert.Error(t, err)
			_, err = tWallet.PublicKey("", addr.Address)
			assert.Error(t, err)
		}
	})

	t.Run("Invalid password", func(t *testing.T) {
		addrs := tWallet.Addresses()
		assert.NotEmpty(t, addrs)
		for _, addr := range addrs {
			_, err := tWallet.PrivateKey("wrong_password", addr.Address)
			assert.Error(t, err)
			_, err = tWallet.PublicKey("wrong_password", addr.Address)
			assert.Error(t, err)
		}
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

	_, err := tWallet.PrivateKey(tPassphrase, crypto.GenerateTestAddress().String())
	assert.Error(t, err)
}

func TestImportPrivateKey(t *testing.T) {
	setup(t)

	_, prv1 := bls.GenerateTestKeyPair()
	assert.NoError(t, tWallet.ImportPrivateKey(tPassphrase, prv1.String()))

	assert.True(t, tWallet.store.Contains(prv1.PublicKey().Address().String()))
	prv2, err := tWallet.PrivateKey(tPassphrase, prv1.PublicKey().Address().String())
	assert.NoError(t, err)
	assert.Equal(t, prv1.String(), prv2)

	// Import again
	assert.Error(t, tWallet.ImportPrivateKey(tPassphrase, prv1.String()))
}

func TestUpdatePassphrase(t *testing.T) {
	setup(t)

	addrs := tWallet.Addresses()
	newPassphrase := "new-passphrase"
	invalidPassphrase := "invalid-passphrase"
	assert.Error(t, tWallet.UpdatePassword("", newPassphrase))
	assert.Error(t, tWallet.UpdatePassword(invalidPassphrase, newPassphrase))
	assert.NoError(t, tWallet.UpdatePassword(tPassphrase, newPassphrase))
	assert.True(t, tWallet.IsEncrypted())
	for _, addr := range addrs {
		assert.True(t, tWallet.Contains(addr.Address))
	}

	assert.Error(t, tWallet.UpdatePassword(invalidPassphrase, newPassphrase))
	assert.NoError(t, tWallet.UpdatePassword(newPassphrase, ""))
	assert.False(t, tWallet.IsEncrypted())
	for _, addr := range addrs {
		assert.True(t, tWallet.Contains(addr.Address))
	}

	assert.Error(t, tWallet.UpdatePassword(invalidPassphrase, newPassphrase))
}

func TestMakeSendTx(t *testing.T) {
	setup(t)
}
