package jsonstorage

import (
	"testing"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:dupword // duplicated seed phrase words
var testMnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon cactus"

func TestUnsupportedWallet(t *testing.T) {
	err := Upgrade("./testdata/unsupported_wallet")
	require.ErrorIs(t, err, UnsupportedVersionError{
		WalletVersion:    6,
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
		data, err := util.ReadFile("./testdata/wallet_version_1")
		require.NoError(t, err)

		tempPath := util.TempFilePath()
		err = util.WriteFile(tempPath, data)
		require.NoError(t, err)

		err = Upgrade(tempPath)
		require.NoError(t, err)

		strg, err := Open(tempPath)
		require.NoError(t, err)

		assert.Equal(t, VersionLatest, strg.WalletInfo().Version)

		infos, err := strg.AllAddresses()
		require.NoError(t, err)
		for _, info := range infos {
			assert.NotEmpty(t, info.PublicKey)
		}
	})

	t.Run("Upgrade Wallet From Version 2", func(t *testing.T) {
		// In this upgrade, the IV for AES is generated from derived bytes using the Argon2id hasher.
		// This test ensures that wallet decryption works after the upgrade.
		data, err := util.ReadFile("./testdata/wallet_version_2")
		require.NoError(t, err)

		tempPath := util.TempFilePath()
		err = util.WriteFile(tempPath, data)
		require.NoError(t, err)

		err = Upgrade(tempPath)
		require.NoError(t, err)

		strg, err := Open(tempPath)
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
		data, err := util.ReadFile("./testdata/wallet_version_3")
		require.NoError(t, err)

		tempPath := util.TempFilePath()
		err = util.WriteFile(tempPath, data)
		require.NoError(t, err)

		err = Upgrade(tempPath)
		require.NoError(t, err)

		strg, err := Open(tempPath)
		require.NoError(t, err)

		assert.Equal(t, VersionLatest, strg.WalletInfo().Version)

		mnemonic, err := strg.Vault().Mnemonic(password)
		require.NoError(t, err)
		assert.Equal(t, "ARGON2ID-AES_256_CBC-MACV1", strg.Vault().Encrypter.Method)
		assert.Equal(t, testMnemonic, mnemonic)
	})

	t.Run("Upgrade Wallet From Version 4", func(t *testing.T) {
		data, err := util.ReadFile("./testdata/wallet_version_4")
		require.NoError(t, err)

		tempPath := util.TempFilePath()
		err = util.WriteFile(tempPath, data)
		require.NoError(t, err)

		err = Upgrade(tempPath)
		require.NoError(t, err)

		strg, err := Open(tempPath)
		require.NoError(t, err)

		assert.Equal(t, VersionLatest, strg.WalletInfo().Version)
		assert.Equal(t, genesis.Mainnet, strg.WalletInfo().Network)
		assert.Equal(t, amount.Amount(2e7), strg.WalletInfo().DefaultFee) // 0.02 PAC

		infos, err := strg.AllAddresses()
		require.NoError(t, err)
		assert.Len(t, infos, 5)
	})
}

func TestUpgradeTestnet(t *testing.T) {
	data, err := util.ReadFile("./testdata/testnet_wallet")
	require.NoError(t, err)

	tempPath := util.TempFilePath()
	err = util.WriteFile(tempPath, data)
	require.NoError(t, err)

	err = Upgrade(tempPath)
	require.NoError(t, err)

	strg, err := Open(tempPath)
	require.NoError(t, err)

	assert.Equal(t, VersionLatest, strg.WalletInfo().Version)
	assert.Equal(t, genesis.Testnet, strg.WalletInfo().Network)
}
