package vault

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/wallet/encrypter"
	"github.com/zarbchain/zarb-go/wallet/hdkeychain"
)

const tPassword = "super_secret_password"

var mnemonic = GenerateMnemonic()
var _, importedPrv = bls.GenerateTestKeyPair()

// testVault return an instances of vault fo testing
func testVault(t *testing.T) *Vault {
	vault, err := CreateVaultFromMnemonic(mnemonic, 21888)
	assert.NoError(t, err)

	for _, p := range vault.Keystore.Purposes {
		key, _ := hdkeychain.NewKeyFromString(p.XPub)
		assert.False(t, key.IsPrivate())
	}

	assert.NoError(t, vault.ImportPrivateKey("", importedPrv))
	assert.False(t, vault.IsEncrypted())

	opts := []encrypter.Option{
		encrypter.OptionIteration(1),
		encrypter.OptionMemory(1),
		encrypter.OptionParallelism(1),
	}

	vault.UpdatePassword("", tPassword, opts...)
	assert.True(t, vault.IsEncrypted())
	return vault
}

func TestAddressInfo(t *testing.T) {
	vault := testVault(t)

	assert.Equal(t, vault.AddressCount(), 22)
	infos := vault.AddressLabels()
	blsIndex := 0
	importedIndex := 0
	for _, i := range infos {
		info := vault.GetAddressInfo(i.Address)
		assert.Equal(t, info.Address, info.Address)
		if !info.Imported {
			assert.Equal(t, info.Path.String(), fmt.Sprintf("m/12381'/21888'/%d/0", blsIndex))
			blsIndex++
		} else {
			assert.True(t, info.Imported)
			assert.Equal(t, info.ImportedIndex, importedIndex)
			importedIndex++
		}
	}

	// Neutered
	neutered := vault.Neuter()
	assert.Equal(t, neutered.AddressCount(), 21)
	infos = neutered.AddressLabels()
	blsIndex = 0
	for _, i := range infos {
		info := vault.GetAddressInfo(i.Address)
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

		assert.Equal(t, recovered.Keystore.Purposes, vault.Keystore.Purposes)
	})
}

// func TestGetPrivateKey(t *testing.T) {
// 	setup(t)

// 	assert.NotEmpty(t, tVault.Addresses)

// 	_, prv := bls.GenerateTestKeyPair()
// 	assert.NoError(t, tVault.ImportPrivateKey(tPassword, prv.String()))

// 	t.Run("Unknown adddress", func(t *testing.T) {
// 		addr := crypto.GenerateTestAddress()
// 		_, err := tVault.PrivateKey(tPassword, addr.String())
// 		assert.ErrorIs(t, err, NewErrAddressNotFound(addr.String()))
// 	})

// 	t.Run("No password", func(t *testing.T) {
// 		for _, info := range tVault.Addresses {
// 			_, err := tVault.PrivateKey("", info.Address)
// 			assert.ErrorIs(t, err, ErrInvalidPassword)
// 			_, err = tVault.PublicKey("", info.Address)
// 			assert.ErrorIs(t, err, ErrInvalidPassword)
// 		}
// 	})

// 	t.Run("Invalid password", func(t *testing.T) {
// 		for _, info := range tVault.Addresses {
// 			_, err := tVault.PrivateKey("wrong_password", info.Address)
// 			assert.ErrorIs(t, err, ErrInvalidPassword)
// 			_, err = tVault.PublicKey("wrong_password", info.Address)
// 			assert.ErrorIs(t, err, ErrInvalidPassword)
// 		}
// 	})

// 	t.Run("Check all the private keys", func(t *testing.T) {
// 		for _, info := range tVault.Addresses {
// 			prvStr, err := tVault.PrivateKey(tPassword, info.Address)
// 			assert.NoError(t, err)
// 			pubStr, err := tVault.PublicKey(tPassword, info.Address)
// 			assert.NoError(t, err)
// 			prv, _ := bls.PrivateKeyFromString(prvStr)
// 			pub, _ := bls.PublicKeyFromString(pubStr)
// 			assert.True(t, prv.PublicKey().EqualsTo(pub))
// 			assert.Equal(t, pub.Address().String(), info.Address)
// 		}
// 	})

// 	t.Run("Invalid method", func(t *testing.T) {
// 		tVault.Addresses[0].Method = "UNKNOWN"
// 		_, err := tVault.PrivateKey(tPassword, tVault.Addresses[0].Address)
// 		assert.ErrorIs(t, err, NewErrUnknownMethod("UNKNOWN"))
// 	})
// }

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
		info := vault.GetAddressInfo(info.Address)
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
			assert.Equal(t, info, vault.GetAddressInfo(info.Address))
		}
	})

	t.Run("Set empty password for the vault", func(t *testing.T) {
		err := vault.UpdatePassword("invalid-password", newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
		assert.NoError(t, vault.UpdatePassword(newPassword, ""))
		assert.False(t, vault.IsEncrypted())
		for _, info := range infos {
			assert.Equal(t, info, vault.GetAddressInfo(info.Address))
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
}

func TestNeuter(t *testing.T) {
	vault := testVault(t)

	neutered := vault.Neuter()
	_, err := neutered.Mnemonic(tPassword)
	assert.ErrorIs(t, err, ErrNeutered)

	_, err = neutered.PrivateKey(tPassword, "any address")
	assert.ErrorIs(t, err, ErrNeutered)

	err = neutered.ImportPrivateKey("any", importedPrv)
	assert.ErrorIs(t, err, ErrNeutered)

	err = vault.Neuter().UpdatePassword("any", "any")
	assert.ErrorIs(t, err, ErrNeutered)
}
