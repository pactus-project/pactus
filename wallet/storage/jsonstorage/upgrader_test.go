package jsonstorage

import (
	"testing"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openTestStorage(t *testing.T, path string) (*Storage, error) {
	t.Helper()

	data, err := util.ReadFile(path)
	require.NoError(t, err)

	tempPath := util.TempFilePath()
	err = util.WriteFile(tempPath, data)
	require.NoError(t, err)

	return Open(tempPath)
}

func TestUnsupportedWallet(t *testing.T) {
	_, err := openTestStorage(t, "./testdata/unsupported_wallet")

	require.ErrorIs(t, err, UnsupportedVersionError{
		WalletVersion:    7,
		SupportedVersion: VersionLatest,
	})
}

// TestUpgrade ensures that old JSON wallets can be safely upgraded to the latest version.
// Encryption parameters are intentionally reduced to speed up the test.
func TestUpgrade(t *testing.T) {
	password := "password"

	t.Run("Upgrade Wallet From Version 1", func(t *testing.T) {
		// In this upgrade, some addresses may not have a public key.
		// This test ensures that after the upgrade, all addresses have a public key.
		strg, err := openTestStorage(t, "./testdata/wallet_version_1")
		require.NoError(t, err)

		assert.Equal(t, VersionLatest, strg.WalletInfo().Version)

		infos := strg.AllAddresses()
		for _, info := range infos {
			assert.NotEmpty(t, info.PublicKey)
		}
	})

	t.Run("Upgrade Wallet From Version 2", func(t *testing.T) {
		// In this upgrade, the IV for AES is generated from derived bytes using the Argon2id hasher.
		// This test ensures that wallet decryption works after the upgrade.
		strg, err := openTestStorage(t, "./testdata/wallet_version_2")
		require.NoError(t, err)

		assert.Equal(t, VersionLatest, strg.WalletInfo().Version)

		assert.Equal(t, "ARGON2ID-AES_256_CTR-MACV1", strg.Vault().Encrypter.Method)
		assert.Equal(t, uint32(32), strg.Vault().Encrypter.Params.GetUint32("keylen"))

		err = strg.Vault().UpdatePassword(password, password)
		require.NoError(t, err)
		assert.Equal(t, "ARGON2ID-AES_256_CTR-MACV1", strg.Vault().Encrypter.Method)
		assert.Equal(t, uint32(48), strg.Vault().Encrypter.Params.GetUint32("keylen"))
	})

	t.Run("Upgrade Wallet From Version 3", func(t *testing.T) {
		// This test ensures that AES encryption supports CBC method.
		strg, err := openTestStorage(t, "./testdata/wallet_version_3")
		require.NoError(t, err)

		assert.Equal(t, VersionLatest, strg.WalletInfo().Version)

		mnemonic, err := strg.Vault().Mnemonic(password)
		require.NoError(t, err)
		assert.Equal(t, "ARGON2ID-AES_256_CBC-MACV1", strg.Vault().Encrypter.Method)
		assert.Equal(t, testMnemonic, mnemonic)
	})

	t.Run("Upgrade Wallet From Version 4", func(t *testing.T) {
		// This test ensures that wallet keeps the default fee.
		strg, err := openTestStorage(t, "./testdata/wallet_version_4")
		require.NoError(t, err)

		assert.Equal(t, VersionLatest, strg.WalletInfo().Version)
		assert.Equal(t, genesis.Mainnet, strg.WalletInfo().Network)
		assert.Equal(t, amount.Amount(2e7), strg.WalletInfo().DefaultFee) // 0.02 PAC

		infos := strg.AllAddresses()
		assert.Len(t, infos, 5)
	})

	t.Run("Upgrade Wallet From Version 5", func(t *testing.T) {
		// This test ensures that wallet supports Secp256k1.
		strg, err := openTestStorage(t, "./testdata/wallet_version_5")
		require.NoError(t, err)

		assert.Equal(t, VersionLatest, strg.WalletInfo().Version)
		assert.Zero(t, strg.Vault().Purposes.PurposeBIP44.NextSexp256k1Index)
	})

	t.Run("Upgrade Wallet From Version 6", func(t *testing.T) {
		// Latest version.
		strg, err := openTestStorage(t, "./testdata/wallet_version_6")
		require.NoError(t, err)
		assert.Equal(t, VersionLatest, strg.WalletInfo().Version)
	})
}

func TestUpgradeTestnet(t *testing.T) {
	strg, err := openTestStorage(t, "./testdata/testnet_wallet")
	require.NoError(t, err)

	assert.Equal(t, VersionLatest, strg.WalletInfo().Version)
	assert.Equal(t, genesis.Testnet, strg.WalletInfo().Network)
}
