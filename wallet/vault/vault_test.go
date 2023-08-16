package vault

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tPassword = "super_secret_password"

type testData struct {
	*testsuite.TestSuite

	vault       *Vault
	mnemonic    string
	importedPrv crypto.PrivateKey
}

// setup return an instances of vault fo testing.
func setup(t *testing.T) *testData {
	ts := testsuite.NewTestSuite(t)

	mnemonic := GenerateMnemonic(128)
	_, importedPrv := ts.RandomBLSKeyPair()
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
	return &testData{
		TestSuite:   ts,
		vault:       vault,
		mnemonic:    mnemonic,
		importedPrv: importedPrv,
	}
}

func TestAddressInfo(t *testing.T) {
	td := setup(t)

	assert.Equal(t, td.vault.AddressCount(), 4)
	infos := td.vault.AddressLabels()
	blsIndex := 0
	importedIndex := 0
	for _, i := range infos {
		info := td.vault.AddressInfo(i.Address)
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
	neutered := td.vault.Neuter()
	assert.Equal(t, neutered.AddressCount(), 3)
	infos = neutered.AddressLabels()
	blsIndex = 0
	for _, i := range infos {
		info := td.vault.AddressInfo(i.Address)
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
	td := setup(t)

	t.Run("Invalid purpose", func(t *testing.T) {
		_, err := td.vault.DeriveNewAddress("", 0)
		assert.ErrorIs(t, err, ErrInvalidPath)
	})

	t.Run("Ok", func(t *testing.T) {
		addr, err := td.vault.DeriveNewAddress("new-addr", PurposeBLS12381)
		assert.NoError(t, err)
		assert.True(t, td.vault.Contains(addr))
		assert.Equal(t, td.vault.Label(addr), "new-addr")
	})
}

func TestRecover(t *testing.T) {
	td := setup(t)

	t.Run("Invalid mnemonic", func(t *testing.T) {
		_, err := CreateVaultFromMnemonic("invalid mnemonic phrase seed", 21888)
		assert.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		recovered, err := CreateVaultFromMnemonic(td.mnemonic, 21888)
		assert.NoError(t, err)

		// Recover addresses
		_, err = recovered.DeriveNewAddress("addr-1", PurposeBLS12381)
		assert.NoError(t, err)
		_, err = recovered.DeriveNewAddress("addr-2", PurposeBLS12381)
		assert.NoError(t, err)
		_, err = recovered.DeriveNewAddress("addr-3", PurposeBLS12381)
		assert.NoError(t, err)

		assert.Equal(t, recovered.Keystore.Purposes, td.vault.Keystore.Purposes)
	})
}

func TestGetPrivateKeys(t *testing.T) {
	td := setup(t)

	t.Run("Unknown address", func(t *testing.T) {
		addr := td.RandomAddress()
		_, err := td.vault.PrivateKeys(tPassword, []string{addr.String()})
		assert.ErrorIs(t, err, NewErrAddressNotFound(addr.String()))
	})

	t.Run("No password", func(t *testing.T) {
		addr := td.vault.AddressLabels()[0].Address
		_, err := td.vault.PrivateKeys("", []string{addr})
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Invalid password", func(t *testing.T) {
		addr := td.vault.AddressLabels()[0].Address
		_, err := td.vault.PrivateKeys("wrong_password", []string{addr})
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Check all the private keys", func(t *testing.T) {
		for _, info := range td.vault.AddressLabels() {
			prv, err := td.vault.PrivateKeys(tPassword, []string{info.Address})
			assert.NoError(t, err)
			i := td.vault.AddressInfo(info.Address)
			require.True(t, prv[0].PublicKey().EqualsTo(i.Pub))
			require.Equal(t, prv[0].PublicKey().Address().String(), info.Address)
		}
	})
}

func TestImportPrivateKey(t *testing.T) {
	td := setup(t)

	t.Run("Reimporting private key", func(t *testing.T) {
		err := td.vault.ImportPrivateKey(tPassword, td.importedPrv)
		assert.ErrorIs(t, err, ErrAddressExists)
	})

	t.Run("Invalid password", func(t *testing.T) {
		_, prv := td.RandomBLSKeyPair()
		err := td.vault.ImportPrivateKey("invalid-password", prv)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		_, prv := td.RandomBLSKeyPair()
		assert.NoError(t, td.vault.ImportPrivateKey(tPassword, prv))
		assert.True(t, td.vault.Contains(prv.PublicKey().Address().String()))
	})
}

func TestGetMnemonic(t *testing.T) {
	td := setup(t)

	t.Run("Invalid password", func(t *testing.T) {
		_, err := td.vault.Mnemonic("invalid-password")
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("No password", func(t *testing.T) {
		_, err := td.vault.Mnemonic("")
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		m, err := td.vault.Mnemonic(tPassword)
		assert.NoError(t, err)
		assert.Equal(t, m, td.mnemonic)
	})

	t.Run("Neutered wallet", func(t *testing.T) {
		_, err := td.vault.Neuter().Mnemonic("")
		assert.ErrorIs(t, err, ErrNeutered)
	})
}

func TestUpdatePassword(t *testing.T) {
	td := setup(t)

	infos := make([]*AddressInfo, 0, td.vault.AddressCount())
	for _, info := range td.vault.AddressLabels() {
		info := td.vault.AddressInfo(info.Address)
		infos = append(infos, info)
	}

	newPassword := "new-password"
	t.Run("Change password", func(t *testing.T) {
		opts := []encrypter.Option{
			encrypter.OptionIteration(1),
			encrypter.OptionMemory(1),
			encrypter.OptionParallelism(1),
		}

		err := td.vault.UpdatePassword("", newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
		err = td.vault.UpdatePassword("invalid-password", newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
		assert.NoError(t, td.vault.UpdatePassword(tPassword, newPassword, opts...))
		assert.True(t, td.vault.IsEncrypted())
		for _, info := range infos {
			assert.Equal(t, info, td.vault.AddressInfo(info.Address))
		}
	})

	t.Run("Set empty password for the vault", func(t *testing.T) {
		err := td.vault.UpdatePassword("invalid-password", newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
		assert.NoError(t, td.vault.UpdatePassword(newPassword, ""))
		assert.False(t, td.vault.IsEncrypted())
		for _, info := range infos {
			assert.Equal(t, info, td.vault.AddressInfo(info.Address))
		}
	})
}

func TestSetLabel(t *testing.T) {
	td := setup(t)

	t.Run("Set label for unknown address", func(t *testing.T) {
		invAddr := td.RandomAddress().String()
		err := td.vault.SetLabel(invAddr, "i have label")
		assert.ErrorIs(t, err, NewErrAddressNotFound(invAddr))
		assert.Equal(t, td.vault.Label(invAddr), "")
	})

	t.Run("Update label", func(t *testing.T) {
		testAddr := td.vault.AddressLabels()[0].Address
		err := td.vault.SetLabel(testAddr, "i have label")
		assert.NoError(t, err)
		assert.Equal(t, td.vault.Label(testAddr), "i have label")
	})

	t.Run("Remove label", func(t *testing.T) {
		testAddr := td.vault.AddressLabels()[0].Address
		err := td.vault.SetLabel(testAddr, "")
		assert.NoError(t, err)
		assert.Empty(t, td.vault.Label(testAddr))
		_, ok := td.vault.Labels[testAddr]
		assert.False(t, ok)
	})
}

func TestNeuter(t *testing.T) {
	td := setup(t)

	neutered := td.vault.Neuter()
	_, err := neutered.Mnemonic(tPassword)
	assert.ErrorIs(t, err, ErrNeutered)

	_, err = neutered.PrivateKeys(tPassword, []string{
		td.RandomAddress().String()})
	assert.ErrorIs(t, err, ErrNeutered)

	err = neutered.ImportPrivateKey("any", td.importedPrv)
	assert.ErrorIs(t, err, ErrNeutered)

	err = td.vault.Neuter().UpdatePassword("any", "any")
	assert.ErrorIs(t, err, ErrNeutered)
}

func TestValidateMnemonic(t *testing.T) {
	tests := []struct {
		mnenomic string
		errStr   string
	}{
		{
			"",
			"Invalid mnenomic",
		},
		{
			"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon",
			"Invalid mnenomic",
		},
		{
			"bandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon",
			"word `bandon` not found in reverse map",
		},
		{
			"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon",
			"Checksum incorrect",
		},
		{
			"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon cactus",
			"",
		},
	}
	for i, test := range tests {
		err := CheckMnemonic(test.mnenomic)
		if err != nil {
			assert.ErrorContains(t, err, test.errStr, "test %v failed", i)
		}
	}
}
