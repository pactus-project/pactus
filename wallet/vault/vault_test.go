package vault

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

var tVault *Vault
var tPassword string

func setup(t *testing.T) {
	password := "super_secret_password"
	mnemonic := GenerateMnemonic()
	vault, err := CreateVaultFromMnemonic(mnemonic, "")
	assert.NoError(t, err)
	assert.False(t, vault.IsEncrypted())

	// create some test addresses
	_, err = vault.MakeNewAddress("", "addr-1")
	assert.NoError(t, err)
	_, err = vault.MakeNewAddress("", "addr-2")
	assert.NoError(t, err)
	_, err = vault.MakeNewAddress("", "addr-3")
	assert.NoError(t, err)

	// Create some keys
	_, prv1 := bls.GenerateTestKeyPair()
	_, prv2 := bls.GenerateTestKeyPair()

	assert.NoError(t, vault.ImportPrivateKey("", prv1.String()))
	assert.NoError(t, vault.ImportPrivateKey("", prv2.String()))

	assert.NoError(t, vault.UpdatePassword("", password))
	assert.True(t, vault.IsEncrypted())
	tPassword = password
	tVault = vault
}

func TestAddressInfo(t *testing.T) {
	setup(t)

	assert.Equal(t, tVault.AddressCount(), 5)
	infos := tVault.AddressInfos()
	for i, info := range infos {
		if tVault.Addresses[i].Method == "IMPORTED" {
			assert.True(t, info.Imported)
		} else {
			assert.Equal(t, info.Label, fmt.Sprintf("addr-%v", i+1))
			assert.False(t, info.Imported)
		}
	}
}

func TestMakeNewAddress(t *testing.T) {
	setup(t)

	t.Run("Invalid password", func(t *testing.T) {
		_, err := tVault.MakeNewAddress("invalid-password", "label")
		assert.ErrorIs(t, err, ErrInvalidPassword)
	})

	t.Run("No password", func(t *testing.T) {
		_, err := tVault.MakeNewAddress("", "label")
		assert.ErrorIs(t, err, ErrInvalidPassword)
	})
}

func TestRecover(t *testing.T) {
	setup(t)

	mnemonic, _ := tVault.Mnemonic(tPassword)
	password := ""

	t.Run("Invalid mnemonic", func(t *testing.T) {
		_, err := CreateVaultFromMnemonic("invali mnemonic phrase seed", password)
		assert.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		recovered, err := CreateVaultFromMnemonic(mnemonic, password)
		assert.NoError(t, err)

		addr1, err := recovered.MakeNewAddress("", "addr-1")
		assert.NoError(t, err)
		addr2, err := recovered.MakeNewAddress("", "addr-2")
		assert.NoError(t, err)
		addr3, err := recovered.MakeNewAddress("", "addr-3")
		assert.NoError(t, err)

		assert.True(t, recovered.Contains(addr1))
		assert.True(t, recovered.Contains(addr2))
		assert.True(t, recovered.Contains(addr3))
	})
}

func TestGetPrivateKey(t *testing.T) {
	setup(t)

	assert.NotEmpty(t, tVault.Addresses)

	_, prv := bls.GenerateTestKeyPair()
	assert.NoError(t, tVault.ImportPrivateKey(tPassword, prv.String()))

	t.Run("Unknown adddress", func(t *testing.T) {
		addr := crypto.GenerateTestAddress()
		_, err := tVault.PrivateKey(tPassword, addr.String())
		assert.Equal(t, err.Error(), NewErrAddressNotFound(addr.String()).Error())
	})

	t.Run("No password", func(t *testing.T) {
		for _, info := range tVault.Addresses {
			_, err := tVault.PrivateKey("", info.Address)
			assert.ErrorIs(t, err, ErrInvalidPassword)
			_, err = tVault.PublicKey("", info.Address)
			assert.ErrorIs(t, err, ErrInvalidPassword)
		}
	})

	t.Run("Invalid password", func(t *testing.T) {
		for _, info := range tVault.Addresses {
			_, err := tVault.PrivateKey("wrong_password", info.Address)
			assert.ErrorIs(t, err, ErrInvalidPassword)
			_, err = tVault.PublicKey("wrong_password", info.Address)
			assert.ErrorIs(t, err, ErrInvalidPassword)
		}
	})

	t.Run("Check all the private keys", func(t *testing.T) {
		for _, info := range tVault.Addresses {
			prvStr, err := tVault.PrivateKey(tPassword, info.Address)
			assert.NoError(t, err)
			pubStr, err := tVault.PublicKey(tPassword, info.Address)
			assert.NoError(t, err)
			prv, _ := bls.PrivateKeyFromString(prvStr)
			pub, _ := bls.PublicKeyFromString(pubStr)
			assert.True(t, prv.PublicKey().EqualsTo(pub))
			assert.Equal(t, pub.Address().String(), info.Address)
		}
	})
}

func TestImportPrivateKey(t *testing.T) {
	setup(t)

	t.Run("Reimporting private key", func(t *testing.T) {
		addr := ""
		// Get first key (address)
		for _, info := range tVault.Addresses {
			addr = info.Address
			break
		}
		prv, err := tVault.PrivateKey(tPassword, addr)
		assert.NoError(t, err)

		// Import again
		err = tVault.ImportPrivateKey(tPassword, prv)
		assert.ErrorIs(t, err, ErrAddressExists)
	})

	t.Run("Invalid private key", func(t *testing.T) {
		err := tVault.ImportPrivateKey(tPassword, "invalid-private-key-string")
		assert.Error(t, err)
	})

	t.Run("Invalid password", func(t *testing.T) {
		_, prv := bls.GenerateTestKeyPair()
		err := tVault.ImportPrivateKey("invalid-password", prv.String())
		assert.ErrorIs(t, err, ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		_, prv := bls.GenerateTestKeyPair()
		assert.NoError(t, tVault.ImportPrivateKey(tPassword, prv.String()))
	})
}

func TestGetMnemonic(t *testing.T) {
	t.Run("Invalid password", func(t *testing.T) {
		_, err := tVault.Mnemonic("invalid-password")
		assert.ErrorIs(t, err, ErrInvalidPassword)
	})

	t.Run("No password", func(t *testing.T) {
		_, err := tVault.Mnemonic("")
		assert.ErrorIs(t, err, ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		_, err := tVault.Mnemonic(tPassword)
		assert.NoError(t, err)

	})
}

func TestUpdatePassword(t *testing.T) {
	setup(t)

	pubs := make([]string, 0, len(tVault.Addresses))
	for _, info := range tVault.Addresses {
		pub, _ := tVault.PublicKey(tPassword, info.Address)
		pubs = append(pubs, pub)
	}

	newPassword := "new-password"
	t.Run("Change password", func(t *testing.T) {
		err := tVault.UpdatePassword("", newPassword)
		assert.ErrorIs(t, err, ErrInvalidPassword)
		err = tVault.UpdatePassword("invalid-password", newPassword)
		assert.ErrorIs(t, err, ErrInvalidPassword)
		assert.NoError(t, tVault.UpdatePassword(tPassword, newPassword))
		assert.True(t, tVault.IsEncrypted())
		for _, info := range tVault.Addresses {
			pub, _ := tVault.PublicKey(newPassword, info.Address)
			assert.Contains(t, pubs, pub)
		}

	})

	t.Run("Set empty password for the vault", func(t *testing.T) {
		err := tVault.UpdatePassword("invalid-password", newPassword)
		assert.ErrorIs(t, err, ErrInvalidPassword)
		assert.NoError(t, tVault.UpdatePassword(newPassword, ""))
		assert.False(t, tVault.IsEncrypted())
		for _, info := range tVault.Addresses {
			pub, _ := tVault.PublicKey("", info.Address)
			assert.Contains(t, pubs, pub)
		}
	})
}
