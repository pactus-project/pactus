package sqlitestorage

import (
	"path"
	"testing"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openTestStorage(t *testing.T, dir string) (*Storage, error) {
	t.Helper()

	data, err := util.ReadFile(path.Join(dir, "wallet.db"))
	require.NoError(t, err)

	tempPath := util.TempDirPath()
	err = util.WriteFile(path.Join(tempPath, "wallet.db"), data)
	require.NoError(t, err)

	return Open(t.Context(), tempPath)
}

func TestUnsupportedWallet(t *testing.T) {
	_, err := openTestStorage(t, "./testdata/unsupported_wallet")
	require.ErrorIs(t, err, UnsupportedVersionError{
		WalletVersion:    3,
		SupportedVersion: VersionLatest,
	})
}

// TestUpgrade ensures that old JSON wallets can be safely upgraded to the latest version.
// Encryption parameters are intentionally reduced to speed up the test.
func TestUpgrade(t *testing.T) {

	t.Run("Upgrade Wallet From Version 1", func(t *testing.T) {
		// This test ensures that wallet supports Secp256k1.
		strg, err := openTestStorage(t, "./testdata/wallet_version_1")
		require.NoError(t, err)

		assert.Equal(t, VersionLatest, strg.WalletInfo().Version)

		infos := strg.AllAddresses()
		for _, info := range infos {
			assert.NotEmpty(t, info.PublicKey)
		}
	})

	t.Run("Upgrade Wallet From Version 2", func(t *testing.T) {
		// Latest version.
		strg, err := openTestStorage(t, "./testdata/wallet_version_2")
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
