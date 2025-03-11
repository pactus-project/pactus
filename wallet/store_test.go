package wallet

import (
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpgradeWallet(t *testing.T) {
	// password is: "password"
	tests := []struct {
		walletPath      string
		upgradedVersion int
	}{
		{"./testdata/wallet_version_1", 2},
		{"./testdata/wallet_version_2", 2},
		{"./testdata/wallet_version_3", 3},
	}

	for _, tt := range tests {
		data, err := util.ReadFile(tt.walletPath)
		require.NoError(t, err)

		tempPath := util.TempFilePath()
		err = util.WriteFile(tempPath, data)
		require.NoError(t, err)

		wlt, err := Open(tempPath, true)
		require.NoError(t, err)

		assert.Equal(t, 4, wlt.AddressCount())
		assert.Equal(t, tt.upgradedVersion, wlt.store.Version)

		infos := wlt.AddressInfos()
		for _, info := range infos {
			assert.NotEmpty(t, info.PublicKey)
		}
	}
}

func TestUnsupportedWallet(t *testing.T) {
	_, err := Open("./testdata/unsupported_wallet", true)
	require.ErrorIs(t, err, UnsupportedVersionError{
		WalletVersion:    4,
		SupportedVersion: VersionLatest,
	})
}
