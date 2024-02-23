package vault

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet2/addresspath"
	"github.com/pactus-project/pactus/wallet2/db"
	"github.com/pactus-project/pactus/wallet2/encrypter"
	"github.com/stretchr/testify/assert"
)

const tPassword = "super_secret_password"

type testData struct {
	*testsuite.TestSuite

	vault       *Vault
	mnemonic    string
	importedPrv crypto.PrivateKey
}

// setup returns an instances of vault fo testing.
func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	mnemonic, _ := GenerateMnemonic(128)
	_, importedPrv := ts.RandBLSKeyPair()

	dbInstance, _ := db.NewDB(":memory:")
	err := dbInstance.CreateTables()
	assert.NoError(t, err)

	vault, err := CreateVaultFromMnemonic(mnemonic, 21888, dbInstance)
	assert.NoError(t, err)

	key, _ := hdkeychain.NewKeyFromString(vault.Purposes.PurposeBLS.XPubAccount)
	assert.False(t, key.IsPrivate())

	// Create some test address
	_, err = vault.NewBLSAccountAddress("addr-1")
	assert.NoError(t, err)
	_, err = vault.NewBLSAccountAddress("addr-2")
	assert.NoError(t, err)
	_, err = vault.NewValidatorAddress("addr-3")
	assert.NoError(t, err)
	_, err = vault.NewValidatorAddress("addr-4")
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

func TestAddresses(t *testing.T) {
	td := setup(t)

	addressCount, _ := td.vault.AddressCount()
	assert.Equal(t, 6, addressCount)
	infos := td.vault.Addresses()
	for _, i := range infos {
		info := td.vault.Address(i.Address)
		assert.Equal(t, i.Address, info.Address)

		cryptoAddr, err := crypto.AddressFromString(i.Address)
		assert.NoError(t, err)

		var xPub string
		if cryptoAddr.IsAccountAddress() {
			xPub = td.vault.Purposes.PurposeBLS.XPubAccount
		} else if cryptoAddr.IsValidatorAddress() {
			xPub = td.vault.Purposes.PurposeBLS.XPubValidator
		}

		ext, err := hdkeychain.NewKeyFromString(xPub)
		assert.NoError(t, err)

		p, err := addresspath.NewPathFromString(i.Path)
		assert.NoError(t, err)

		if p.IsBLSPurpose() {
			extendedKey, err := ext.Derive(p.AddressIndex())
			assert.NoError(t, err)

			blsPubKey, err := bls.PublicKeyFromBytes(extendedKey.RawPublicKey())
			assert.NoError(t, err)

			assert.Equal(t, blsPubKey.String(), info.PublicKey)
		}

		addr, _ := crypto.AddressFromString(info.Address)
		path, _ := addresspath.NewPathFromString(info.Path)

		switch path.Purpose() {
		case H(PurposeBLS12381):
			if addr.IsValidatorAddress() {
				assert.Equal(t, info.Path, fmt.Sprintf("m/%d'/%d'/1'/%d",
					PurposeBLS12381, td.vault.CoinType, path.AddressIndex()))
			}

			if addr.IsAccountAddress() {
				assert.Equal(t, info.Path, fmt.Sprintf("m/%d'/%d'/2'/%d",
					PurposeBLS12381, td.vault.CoinType, path.AddressIndex()))
			}
		case H(PurposeImportPrivateKey):
			if addr.IsValidatorAddress() {
				assert.Equal(t, info.Path, fmt.Sprintf("m/%d'/%d'/1'/%d'",
					PurposeImportPrivateKey, td.vault.CoinType, path.AddressIndex()-hdkeychain.HardenedKeyStart))
			}

			if addr.IsAccountAddress() {
				assert.Equal(t, info.Path, fmt.Sprintf("m/%d'/%d'/2'/%d'",
					PurposeImportPrivateKey, td.vault.CoinType, path.AddressIndex()-hdkeychain.HardenedKeyStart))
			}
		}
	}

	// Neutered
	neutered := td.vault.Neuter()
	neuteredAddressCount, err := neutered.AddressCount()
	assert.Nil(t, err)
	assert.Equal(t, 6, neuteredAddressCount)
}

func TestSortAddresses(t *testing.T) {
	td := setup(t)

	addressCount, err := td.vault.AddressCount()
	assert.NoError(t, err)

	assert.Equal(t, 6, addressCount)
	infos := td.vault.Addresses()

	assert.Equal(t, "m/12381'/21888'/1'/0", infos[0].Path)
	assert.Equal(t, "m/65535'/21888'/2'/0'", infos[len(infos)-1].Path)
}
