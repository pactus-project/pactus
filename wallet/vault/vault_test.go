package vault

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tPassword = "super_secret_password"

var mnemonic = GenerateMnemonic(128)
var _, importedPrv = bls.GenerateTestKeyPair()

// testVault return an instances of vault fo testing
func testVault(t *testing.T) *Vault {
	vault, err := CreateVaultFromMnemonic(mnemonic, 21888)
	assert.NoError(t, err)

	for _, p := range vault.Keystore.Purposes {
		key, _ := hdkeychain.NewKeyFromString(p.XPub)
		assert.False(t, key.IsPrivate())
	}

	// Create some test address
	_, err = vault.DeriveNewAddress("addr-1", PurposeBLS12381)
	assert.NoError(t, err)
	_, err = vault.DeriveNewAddress("addr-2", PurposeBLS12381)
	assert.NoError(t, err)
	_, err = vault.DeriveNewAddress("addr-3", PurposeBLS12381)
	assert.NoError(t, err)

	assert.NoError(t, vault.ImportPrivateKey("", importedPrv))
	assert.False(t, vault.IsEncrypted())

	opts := []encrypter.Option{
		encrypter.OptionIteration(1),
		encrypter.OptionMemory(1),
		encrypter.OptionParallelism(1),
	}

	err = vault.UpdatePassword("", tPassword, opts...)
	assert.NoError(t, err)
	assert.True(t, vault.IsEncrypted())
	return vault
}

func TestAddressInfo(t *testing.T) {
	vault := testVault(t)

	assert.Equal(t, vault.AddressCount(), 4)
	infos := vault.AddressLabels()
	blsIndex := 0
	importedIndex := 0
	for _, i := range infos {
		info := vault.AddressInfo(i.Address)
		assert.Equal(t, info.Address, info.Address)
		if !info.Imported {
			assert.Equal(t, info.Path.String(), fmt.Sprintf("m/12381'/21888'/%d/0", blsIndex))
			blsIndex++
		} else {
			assert.True(t, info.Imported)
			assert.Equal(t, info.ImportedIndex, importedIndex)
			importedIndex++
		}
		assert.Equal(t, info.Pub.Address().String(), info.Address)
		assert.Equal(t, info.Address, i.Address)
	}

	// Neutered
	neutered := vault.Neuter()
	assert.Equal(t, neutered.AddressCount(), 3)
	infos = neutered.AddressLabels()
	blsIndex = 0
	for _, i := range infos {
		info := vault.AddressInfo(i.Address)
		assert.Equal(t, info.Address, info.Address)
		if !info.Imported {
			assert.Equal(t, info.Path.String(), fmt.Sprintf("m/12381'/21888'/%d/0", blsIndex))
			blsIndex++
		} else {
			assert.Error(t, ErrNeutered)
		}
	}
}

func TestDeriveNewAddress(t *testing.T) {
	vault := testVault(t)

	t.Run("Invalid purpose", func(t *testing.T) {
		_, err := vault.DeriveNewAddress("", 0)
		assert.ErrorIs(t, err, ErrInvalidPath)
	})

	t.Run("Ok", func(t *testing.T) {
		addr, err := vault.DeriveNewAddress("new-addr", PurposeBLS12381)
		assert.NoError(t, err)
		assert.True(t, vault.Contains(addr))
		assert.Equal(t, vault.Label(addr), "new-addr")
	})
}

func TestRecover(t *testing.T) {
	vault := testVault(t)

	t.Run("Invalid mnemonic", func(t *testing.T) {
		_, err := CreateVaultFromMnemonic("invalid mnemonic phrase seed", 21888)
		assert.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		recovered, err := CreateVaultFromMnemonic(mnemonic, 21888)
		assert.NoError(t, err)

		// Recover addresses
		_, err = recovered.DeriveNewAddress("addr-1", PurposeBLS12381)
		assert.NoError(t, err)
		_, err = recovered.DeriveNewAddress("addr-2", PurposeBLS12381)
		assert.NoError(t, err)
		_, err = recovered.DeriveNewAddress("addr-3", PurposeBLS12381)
		assert.NoError(t, err)

		assert.Equal(t, recovered.Keystore.Purposes, vault.Keystore.Purposes)
	})
}

func TestGetPrivateKeys(t *testing.T) {
	vault := testVault(t)

	t.Run("Unknown address", func(t *testing.T) {
		addr := crypto.GenerateTestAddress()
		_, err := vault.PrivateKeys(tPassword, []string{addr.String()})
		assert.ErrorIs(t, err, NewErrAddressNotFound(addr.String()))
	})

	t.Run("No password", func(t *testing.T) {
		addr := vault.AddressLabels()[0].Address
		_, err := vault.PrivateKeys("", []string{addr})
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Invalid password", func(t *testing.T) {
		addr := vault.AddressLabels()[0].Address
		_, err := vault.PrivateKeys("wrong_password", []string{addr})
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Check all the private keys", func(t *testing.T) {
		for _, info := range vault.AddressLabels() {
			prv, err := vault.PrivateKeys(tPassword, []string{info.Address})
			assert.NoError(t, err)
			i := vault.AddressInfo(info.Address)
			require.True(t, prv[0].PublicKey().EqualsTo(i.Pub))
			require.Equal(t, prv[0].PublicKey().Address().String(), info.Address)
		}
	})
}

func TestImportPrivateKey(t *testing.T) {
	vault := testVault(t)

	t.Run("Reimporting private key", func(t *testing.T) {
		err := vault.ImportPrivateKey(tPassword, importedPrv)
		assert.ErrorIs(t, err, ErrAddressExists)
	})

	t.Run("Invalid password", func(t *testing.T) {
		_, prv := bls.GenerateTestKeyPair()
		err := vault.ImportPrivateKey("invalid-password", prv)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		_, prv := bls.GenerateTestKeyPair()
		assert.NoError(t, vault.ImportPrivateKey(tPassword, prv))
		assert.True(t, vault.Contains(prv.PublicKey().Address().String()))
	})
}

func TestGetMnemonic(t *testing.T) {
	vault := testVault(t)

	t.Run("Invalid password", func(t *testing.T) {
		_, err := vault.Mnemonic("invalid-password")
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("No password", func(t *testing.T) {
		_, err := vault.Mnemonic("")
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		m, err := vault.Mnemonic(tPassword)
		assert.NoError(t, err)
		assert.Equal(t, m, mnemonic)
	})

	t.Run("Neutered wallet", func(t *testing.T) {
		_, err := vault.Neuter().Mnemonic("")
		assert.ErrorIs(t, err, ErrNeutered)
	})
}

func TestUpdatePassword(t *testing.T) {
	vault := testVault(t)

	infos := make([]*AddressInfo, 0, vault.AddressCount())
	for _, info := range vault.AddressLabels() {
		info := vault.AddressInfo(info.Address)
		infos = append(infos, info)
	}

	newPassword := "new-password"
	t.Run("Change password", func(t *testing.T) {
		opts := []encrypter.Option{
			encrypter.OptionIteration(1),
			encrypter.OptionMemory(1),
			encrypter.OptionParallelism(1),
		}

		err := vault.UpdatePassword("", newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
		err = vault.UpdatePassword("invalid-password", newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
		assert.NoError(t, vault.UpdatePassword(tPassword, newPassword, opts...))
		assert.True(t, vault.IsEncrypted())
		for _, info := range infos {
			assert.Equal(t, info, vault.AddressInfo(info.Address))
		}
	})

	t.Run("Set empty password for the vault", func(t *testing.T) {
		err := vault.UpdatePassword("invalid-password", newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
		assert.NoError(t, vault.UpdatePassword(newPassword, ""))
		assert.False(t, vault.IsEncrypted())
		for _, info := range infos {
			assert.Equal(t, info, vault.AddressInfo(info.Address))
		}
	})
}

func TestSetLabel(t *testing.T) {
	vault := testVault(t)

	t.Run("Set label for unknown address", func(t *testing.T) {
		invAddr := crypto.GenerateTestAddress().String()
		err := vault.SetLabel(invAddr, "i have label")
		assert.ErrorIs(t, err, NewErrAddressNotFound(invAddr))
		assert.Equal(t, vault.Label(invAddr), "")
	})

	t.Run("Update label", func(t *testing.T) {
		testAddr := vault.AddressLabels()[0].Address
		err := vault.SetLabel(testAddr, "i have label")
		assert.NoError(t, err)
		assert.Equal(t, vault.Label(testAddr), "i have label")
	})

	t.Run("Remove label", func(t *testing.T) {
		testAddr := vault.AddressLabels()[0].Address
		err := vault.SetLabel(testAddr, "")
		assert.NoError(t, err)
		assert.Empty(t, vault.Label(testAddr))
		_, ok := vault.Labels[testAddr]
		assert.False(t, ok)
	})
}

func TestNeuter(t *testing.T) {
	vault := testVault(t)

	neutered := vault.Neuter()
	_, err := neutered.Mnemonic(tPassword)
	assert.ErrorIs(t, err, ErrNeutered)

	_, err = neutered.PrivateKeys(tPassword, []string{
		crypto.GenerateTestAddress().String()})
	assert.ErrorIs(t, err, ErrNeutered)

	err = neutered.ImportPrivateKey("any", importedPrv)
	assert.ErrorIs(t, err, ErrNeutered)

	err = vault.Neuter().UpdatePassword("any", "any")
	assert.ErrorIs(t, err, ErrNeutered)
}
