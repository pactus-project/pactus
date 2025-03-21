package wallet

import (
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSupportedWallets(t *testing.T) {
	// password is: "password"
	tests := []struct {
		walletPath   string
		addressCount int
	}{
		{"./testdata/wallet_version_1", 4},
		{"./testdata/wallet_version_2", 5},
		{"./testdata/wallet_version_3", 5},
	}

	for _, tt := range tests {
		data, err := util.ReadFile(tt.walletPath)
		require.NoError(t, err)

		tempPath := util.TempFilePath()
		err = util.WriteFile(tempPath, data)
		require.NoError(t, err)

		wlt, err := Open(tempPath, true)
		require.NoError(t, err)

		// TODO: use public method to check version, like Wallet.Info()
		assert.Equal(t, VersionLatest, wlt.store.Version)
		assert.Equal(t, tt.addressCount, wlt.AddressCount())

		mnemonic, err := wlt.Mnemonic("password")
		require.NoError(t, err)
		//nolint:dupword // duplicated seed phrase words
		assert.Equal(t,
			"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon cactus", mnemonic)
	}
}

func TestUnsupportedWallet(t *testing.T) {
	_, err := Open("./testdata/unsupported_wallet", true)
	require.ErrorIs(t, err, UnsupportedVersionError{
		WalletVersion:    4,
		SupportedVersion: VersionLatest,
	})
}

// Tis test ensures that old wallets can be safely upgraded to the latest version.
// Encryption parameters are intentionally reduced to speed up the test.
func TestUpgradeWallet(t *testing.T) {
	t.Run("Upgrade Wallet Version 1", func(t *testing.T) {
		// In this upgrade, some addresses may not have a public key.
		// This test ensures that after the upgrade, all addresses have a public key.
		data, err := util.ReadFile("./testdata/wallet_version_1")
		require.NoError(t, err)

		tempPath := util.TempFilePath()
		err = util.WriteFile(tempPath, data)
		require.NoError(t, err)

		wlt, err := Open(tempPath, true)
		require.NoError(t, err)

		for _, info := range wlt.AddressInfos() {
			assert.NotEmpty(t, info.PublicKey)
		}
	})

	t.Run("Upgrade Wallet Version 2 (encrypted)", func(t *testing.T) {
		// In this upgrade, the IV for AES is generated from derived bytes using the Argon2id hasher.
		// This test ensures that wallet decryption works after the upgrade.
		data, err := util.ReadFile("./testdata/wallet_version_2")
		require.NoError(t, err)

		tempPath := util.TempFilePath()
		err = util.WriteFile(tempPath, data)
		require.NoError(t, err)

		wlt, err := Open(tempPath, true)
		require.NoError(t, err)

		password := "password"

		err = wlt.UpdatePassword(password, password)
		require.NoError(t, err)
		assert.Equal(t, "ARGON2ID-AES_256_CTR-MACV1", wlt.store.Vault.Encrypter.Method)
		assert.Equal(t, uint32(48), wlt.store.Vault.Encrypter.Params.GetUint32("keylen"))

		addrInfo, err := wlt.NewEd25519AccountAddress("", password)
		require.NoError(t, err)
		assert.Equal(t, "pc1r7aynw9urvh66ktr3fte2gskjjnxzruflkgde94", addrInfo.Address)
	})
}
