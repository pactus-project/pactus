package vault

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tPassword = "super_secret_password"

type testData struct {
	*testsuite.TestSuite

	vault              *Vault
	mnemonic           string
	importedBLSPrv     *bls.PrivateKey
	importedEd25519Prv *ed25519.PrivateKey
}

// setup returns an instances of vault fo testing.
func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	mnemonic, _ := GenerateMnemonic(128)
	vault, err := CreateVaultFromMnemonic(mnemonic, 21888)
	assert.NoError(t, err)

	key, _ := hdkeychain.NewKeyFromString(vault.Purposes.PurposeBLS.XPubAccount)
	assert.False(t, key.IsPrivate())

	// Create some test address
	_, err = vault.NewBLSAccountAddress("bls-account-address")
	assert.NoError(t, err)
	_, err = vault.NewEd25519AccountAddress("ed25519-account-address", "")
	assert.NoError(t, err)
	_, err = vault.NewValidatorAddress("validator-address")
	assert.NoError(t, err)

	_, importedBLSPrv := ts.RandBLSKeyPair()
	assert.NoError(t, vault.ImportBLSPrivateKey("", importedBLSPrv))

	_, importedEd25519Prv := ts.RandEd25519KeyPair()
	assert.NoError(t, vault.ImportEd25519PrivateKey("", importedEd25519Prv))

	assert.False(t, vault.IsEncrypted())

	// Set encryption options to minimal values for faster test execution.
	opts := []encrypter.Option{
		encrypter.OptionIteration(1),
		encrypter.OptionMemory(1),
		encrypter.OptionParallelism(1),
	}

	err = vault.UpdatePassword("", tPassword, opts...)
	assert.NoError(t, err)
	assert.True(t, vault.IsEncrypted())

	return &testData{
		TestSuite:          ts,
		vault:              vault,
		mnemonic:           mnemonic,
		importedBLSPrv:     importedBLSPrv,
		importedEd25519Prv: importedEd25519Prv,
	}
}

func TestAddressCount(t *testing.T) {
	td := setup(t)

	assert.Equal(t, 6, td.vault.AddressCount())

	// Neutered
	neutered := td.vault.Neuter()
	assert.Equal(t, 6, neutered.AddressCount())
}

func TestContains(t *testing.T) {
	td := setup(t)

	t.Run("Vault should contain all known addresses", func(t *testing.T) {
		infos := td.vault.AddressInfos()
		for _, i := range infos {
			assert.True(t, td.vault.Contains(i.Address))
		}
	})

	t.Run("Vault should not contain unknown address", func(t *testing.T) {
		unknownAddr := td.RandAccAddress().String()
		assert.False(t, td.vault.Contains(unknownAddr))
	})
}

func TestSortAddressInfo(t *testing.T) {
	td := setup(t)

	infos := td.vault.AddressInfos()

	// Ed25519 Keys
	assert.Equal(t, "m/44'/21888'/3'/0'", infos[0].Path)
	// BLS Keys
	assert.Equal(t, "m/12381'/21888'/1'/0", infos[1].Path)
	assert.Equal(t, "m/12381'/21888'/2'/0", infos[2].Path)
	// Imported Keys
	assert.Equal(t, "m/65535'/21888'/1'/0'", infos[3].Path)
	assert.Equal(t, "m/65535'/21888'/2'/0'", infos[4].Path)
	assert.Equal(t, "m/65535'/21888'/3'/1'", infos[5].Path)
}

func TestAllAccountAddresses(t *testing.T) {
	td := setup(t)

	accountAddrs := td.vault.AllAccountAddresses()
	for _, i := range accountAddrs {
		path, err := addresspath.FromString(i.Path)
		assert.NoError(t, err)

		assert.NotEqual(t, _H(crypto.AddressTypeValidator), path.AddressType())
	}
}

func TestAllValidatorAddresses(t *testing.T) {
	td := setup(t)

	validatorAddrs := td.vault.AllValidatorAddresses()
	for _, i := range validatorAddrs {
		info := td.vault.AddressInfo(i.Address)
		assert.Equal(t, i.Address, info.Address)

		path, _ := addresspath.FromString(info.Path)

		switch path.Purpose() {
		case _H(PurposeBLS12381):
			assert.Equal(t, fmt.Sprintf("m/%d'/%d'/1'/%d",
				PurposeBLS12381, td.vault.CoinType, path.AddressIndex()), info.Path)
		case _H(PurposeImportPrivateKey):
			assert.Equal(t, fmt.Sprintf("m/%d'/%d'/1'/%d'",
				PurposeImportPrivateKey, td.vault.CoinType, _N(path.AddressIndex())), info.Path)
		default:
			assert.Fail(t, "not supported")
		}
	}
}

func TestSortAllValidatorAddresses(t *testing.T) {
	td := setup(t)

	validatorAddrs := td.vault.AllValidatorAddresses()

	assert.Equal(t, "m/12381'/21888'/1'/0", validatorAddrs[0].Path)
	assert.Equal(t, "m/65535'/21888'/1'/0'", validatorAddrs[len(validatorAddrs)-1].Path)
}

func TestAddressFromPath(t *testing.T) {
	td := setup(t)

	t.Run("Could not find address from path", func(t *testing.T) {
		path := "m/12381'/26888'/983'/0"
		assert.Nil(t, td.vault.AddressFromPath(path))
	})

	t.Run("Ok", func(t *testing.T) {
		var address string
		var addrInfo AddressInfo

		for addr, ai := range td.vault.Addresses {
			address = addr
			addrInfo = ai

			break
		}

		assert.Equal(t, address, td.vault.AddressFromPath(addrInfo.Path).Address)
	})
}

func TestNewValidatorAddress(t *testing.T) {
	td := setup(t)

	label := td.RandString(16)
	addressInfo, err := td.vault.NewValidatorAddress(label)
	assert.NoError(t, err)
	assert.NotEmpty(t, addressInfo.Address)
	assert.NotEmpty(t, addressInfo.PublicKey)
	assert.Contains(t, addressInfo.Path, "m/12381'/21888'/1'")
	assert.Equal(t, label, addressInfo.Label)

	pub, _ := bls.PublicKeyFromString(addressInfo.PublicKey)
	assert.Equal(t, pub.ValidatorAddress().String(), addressInfo.Address)
}

func TestNewBLSAccountAddress(t *testing.T) {
	td := setup(t)

	label := td.RandString(16)
	addressInfo, err := td.vault.NewBLSAccountAddress(label)
	assert.NoError(t, err)
	assert.NotEmpty(t, addressInfo.Address)
	assert.NotEmpty(t, addressInfo.PublicKey)
	assert.Contains(t, addressInfo.Path, "m/12381'/21888'/2'")
	assert.Equal(t, label, addressInfo.Label)

	pub, _ := bls.PublicKeyFromString(addressInfo.PublicKey)
	assert.Equal(t, pub.AccountAddress().String(), addressInfo.Address)
}

func TestNewE225519AccountAddress(t *testing.T) {
	td := setup(t)

	addressInfo, err := td.vault.NewEd25519AccountAddress("addr-2", tPassword)
	assert.NoError(t, err)
	assert.NotEmpty(t, addressInfo.Address)
	assert.NotEmpty(t, addressInfo.PublicKey)
	assert.Equal(t, "m/44'/21888'/3'/1'", addressInfo.Path)

	pub, _ := ed25519.PublicKeyFromString(addressInfo.PublicKey)
	assert.Equal(t, pub.AccountAddress().String(), addressInfo.Address)
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
		_, err = recovered.NewBLSAccountAddress("bls-account-address")
		assert.NoError(t, err)
		_, err = recovered.NewEd25519AccountAddress("ed25519-account-address", "")
		assert.NoError(t, err)
		_, err = recovered.NewValidatorAddress("validator-address")
		assert.NoError(t, err)

		assert.Equal(t, recovered.Purposes, td.vault.Purposes)
	})
}

func TestGetPrivateKeys(t *testing.T) {
	td := setup(t)

	t.Run("Unknown address", func(t *testing.T) {
		addr := td.RandAccAddress()
		_, err := td.vault.PrivateKeys(tPassword, []string{addr.String()})
		assert.ErrorIs(t, err, NewErrAddressNotFound(addr.String()))
	})

	t.Run("No password", func(t *testing.T) {
		addr := td.vault.AddressInfos()[0].Address
		_, err := td.vault.PrivateKeys("", []string{addr})
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Invalid password", func(t *testing.T) {
		addr := td.vault.AddressInfos()[0].Address
		_, err := td.vault.PrivateKeys("wrong_password", []string{addr})
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Check all the private keys", func(t *testing.T) {
		for _, info := range td.vault.AddressInfos() {
			prv, err := td.vault.PrivateKeys(tPassword, []string{info.Address})
			assert.NoError(t, err)
			addrInfo := td.vault.AddressInfo(info.Address)
			path, _ := addresspath.FromString(info.Path)

			switch _N(path.AddressType()) {
			case uint32(crypto.AddressTypeBLSAccount),
				uint32(crypto.AddressTypeValidator):
				pub, _ := bls.PublicKeyFromString(addrInfo.PublicKey)
				require.True(t, prv[0].PublicKey().EqualsTo(pub))
			case uint32(crypto.AddressTypeEd25519Account):
				pub, _ := ed25519.PublicKeyFromString(addrInfo.PublicKey)
				require.True(t, prv[0].PublicKey().EqualsTo(pub))
			default:
				assert.Fail(t, "not supported")
			}
		}
	})
}

func TestImportBLSPrivateKey(t *testing.T) {
	td := setup(t)

	_, prv := td.RandBLSKeyPair()

	t.Run("Invalid password", func(t *testing.T) {
		err := td.vault.ImportBLSPrivateKey("invalid-password", prv)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		err := td.vault.ImportBLSPrivateKey(tPassword, prv)
		assert.NoError(t, err)

		valAddr := prv.PublicKeyNative().ValidatorAddress().String()
		accAddr := prv.PublicKeyNative().AccountAddress().String()

		valAddrInfo := td.vault.AddressInfo(valAddr)
		accAddrInfo := td.vault.AddressInfo(accAddr)

		assert.True(t, td.vault.Contains(valAddr))
		assert.True(t, td.vault.Contains(accAddr))

		assert.Equal(t, valAddr, valAddrInfo.Address)
		assert.Equal(t, accAddr, accAddrInfo.Address)

		assert.Equal(t, prv.PublicKeyNative().String(), valAddrInfo.PublicKey)
		assert.Equal(t, prv.PublicKeyNative().String(), accAddrInfo.PublicKey)

		assert.Equal(t, "m/65535'/21888'/1'/2'", valAddrInfo.Path)
		assert.Equal(t, "m/65535'/21888'/2'/2'", accAddrInfo.Path)
	})

	t.Run("Reimporting private key", func(t *testing.T) {
		err := td.vault.ImportBLSPrivateKey(tPassword, prv)
		assert.ErrorIs(t, err, ErrAddressExists)
	})
}

func TestImportEd25519PrivateKey(t *testing.T) {
	td := setup(t)

	_, prv := td.RandEd25519KeyPair()

	t.Run("Invalid password", func(t *testing.T) {
		err := td.vault.ImportEd25519PrivateKey("invalid-password", prv)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Ok", func(t *testing.T) {
		err := td.vault.ImportEd25519PrivateKey(tPassword, prv)
		assert.NoError(t, err)

		accAddr := prv.PublicKeyNative().AccountAddress().String()

		accAddrInfo := td.vault.AddressInfo(accAddr)
		assert.True(t, td.vault.Contains(accAddr))
		assert.Equal(t, accAddr, accAddrInfo.Address)
		assert.Equal(t, prv.PublicKeyNative().String(), accAddrInfo.PublicKey)
		assert.Equal(t, "m/65535'/21888'/3'/2'", accAddrInfo.Path)
	})

	t.Run("Reimporting private key", func(t *testing.T) {
		err := td.vault.ImportEd25519PrivateKey(tPassword, td.importedEd25519Prv)
		assert.ErrorIs(t, err, ErrAddressExists)
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

	opts := []encrypter.Option{
		encrypter.OptionIteration(1),
		encrypter.OptionMemory(1),
		encrypter.OptionParallelism(1),
	}

	addrInfos := td.vault.AddressInfos()
	newPassword := "new-password"

	t.Run("Empty password", func(t *testing.T) {
		err := td.vault.UpdatePassword("", newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Incorrect password", func(t *testing.T) {
		err := td.vault.UpdatePassword("invalid-password", newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Valid password update", func(t *testing.T) {
		assert.NoError(t, td.vault.UpdatePassword(tPassword, newPassword, opts...))
		assert.True(t, td.vault.IsEncrypted())
		assert.Equal(t, addrInfos, td.vault.AddressInfos())
	})

	t.Run("Old password should no longer be valid", func(t *testing.T) {
		err := td.vault.UpdatePassword(tPassword, newPassword)
		assert.ErrorIs(t, err, encrypter.ErrInvalidPassword)
	})

	t.Run("Set vault password to empty", func(t *testing.T) {
		assert.NoError(t, td.vault.UpdatePassword(newPassword, ""))
		assert.False(t, td.vault.IsEncrypted())
		assert.Equal(t, addrInfos, td.vault.AddressInfos())
	})
}

func TestSetLabel(t *testing.T) {
	td := setup(t)

	t.Run("Set label for unknown address", func(t *testing.T) {
		invAddr := td.RandAccAddress().String()
		err := td.vault.SetLabel(invAddr, "i have label")
		assert.ErrorIs(t, err, NewErrAddressNotFound(invAddr))
		assert.Equal(t, "", td.vault.Label(invAddr))
	})

	t.Run("Update label", func(t *testing.T) {
		testAddr := td.vault.AddressInfos()[0].Address
		err := td.vault.SetLabel(testAddr, "I have a label")
		assert.NoError(t, err)
		assert.Equal(t, "I have a label", td.vault.Label(testAddr))
	})

	t.Run("Remove label", func(t *testing.T) {
		testAddr := td.vault.AddressInfos()[0].Address
		err := td.vault.SetLabel(testAddr, "")
		assert.NoError(t, err)
		var ok bool
		l := td.vault.Label(testAddr)
		if strings.TrimSpace(l) != "" {
			ok = true
		}
		assert.Empty(t, td.vault.Label(testAddr))
		assert.False(t, ok)
	})
}

func TestNeuter(t *testing.T) {
	td := setup(t)

	neutered := td.vault.Neuter()
	_, err := neutered.Mnemonic(tPassword)
	assert.ErrorIs(t, err, ErrNeutered)

	_, err = neutered.PrivateKeys(tPassword, []string{
		td.RandAccAddress().String(),
	})
	assert.ErrorIs(t, err, ErrNeutered)

	_, prv := td.RandBLSKeyPair()
	err = neutered.ImportBLSPrivateKey("any", prv)
	assert.ErrorIs(t, err, ErrNeutered)

	err = td.vault.Neuter().UpdatePassword("any", "any")
	assert.ErrorIs(t, err, ErrNeutered)
}
