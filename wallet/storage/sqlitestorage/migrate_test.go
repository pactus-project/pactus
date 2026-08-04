package sqlitestorage

import (
	"testing"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet/storage/jsonstorage"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	opts := []encrypter.Option{
		encrypter.OptionIteration(1),
		encrypter.OptionMemory(8),
		encrypter.OptionParallelism(1),
	}
	vlt, err := vault.CreateVaultFromMnemonic(testMnemonic, addresspath.CoinTypePactusTestnet, "password", opts...)
	require.NoError(t, err)

	jsonPath := util.TempFilePath()
	jsonStrg, err := jsonstorage.Create(jsonPath, genesis.Testnet, vlt)
	require.NoError(t, err)

	// Add some addresses to the legacy wallet.
	expectedAddrs := make([]*types.AddressInfo, 3)
	for i := 0; i < len(expectedAddrs); i++ {
		addrInfo := &types.AddressInfo{
			Address:   ts.RandAccAddress().String(),
			PublicKey: ts.RandString(32),
			Label:     ts.RandString(16),
			Path:      ts.RandString(16),
		}
		require.NoError(t, jsonStrg.InsertAddress(addrInfo))
		expectedAddrs[i] = addrInfo
	}

	jsonInfo := jsonStrg.WalletInfo()

	sqlitePath := util.TempDirPath()
	strg, err := Migrate(t.Context(), jsonPath, sqlitePath)
	require.NoError(t, err)
	defer func() { _ = strg.Close() }()

	// The legacy JSON wallet should be left untouched.
	assert.True(t, util.PathExists(jsonPath))

	// The wallet metadata should be preserved.
	info := strg.WalletInfo()
	assert.Equal(t, "SQLite", info.Driver)
	assert.Equal(t, jsonInfo.Network, info.Network)
	assert.Equal(t, jsonInfo.DefaultFee, info.DefaultFee)
	assert.Equal(t, jsonInfo.UUID, info.UUID)
	assert.True(t, jsonInfo.CreatedAt.Equal(info.CreatedAt))

	// The vault should be preserved.
	assert.Equal(t, jsonStrg.Vault(), strg.Vault())

	// The addresses should be preserved.
	require.Equal(t, len(expectedAddrs), strg.AddressCount())
	for _, expected := range expectedAddrs {
		got, err := strg.AddressInfo(expected.Address)
		require.NoError(t, err)
		assert.Equal(t, expected.Address, got.Address)
		assert.Equal(t, expected.PublicKey, got.PublicKey)
		assert.Equal(t, expected.Label, got.Label)
		assert.Equal(t, expected.Path, got.Path)
	}
}

func TestMigrateInvalidJSON(t *testing.T) {
	jsonPath := util.TempFilePath()
	require.NoError(t, util.WriteFile(jsonPath, []byte("invalid_data")))

	_, err := Migrate(t.Context(), jsonPath, util.TempDirPath())
	require.Error(t, err)

	// The invalid file should not be deleted.
	assert.True(t, util.PathExists(jsonPath))
}

// TestMigrateLegacyWallet migrates a real legacy (version 4) JSON wallet that
// needs to be upgraded to the latest version before the migration.
func TestMigrateLegacyWallet(t *testing.T) {
	data, err := util.ReadFile("../jsonstorage/testdata/wallet_version_4")
	require.NoError(t, err)

	jsonPath := util.TempFilePath()
	require.NoError(t, util.WriteFile(jsonPath, data))

	sqlitePath := util.TempDirPath()
	strg, err := Migrate(t.Context(), jsonPath, sqlitePath)
	require.NoError(t, err)
	defer func() { _ = strg.Close() }()

	// The legacy JSON wallet should be left untouched.
	assert.True(t, util.PathExists(jsonPath))

	// The wallet metadata should be preserved.
	info := strg.WalletInfo()
	assert.Equal(t, VersionLatest, info.Version)
	assert.Equal(t, genesis.Mainnet, info.Network)
	assert.Equal(t, amount.Amount(2e7), info.DefaultFee) // 0.02 PAC

	// The vault should be intact.
	assert.True(t, strg.Vault().IsEncrypted())

	// The addresses should be preserved (wallet_version_4 has 5 addresses).
	assert.Equal(t, 5, strg.AddressCount())
}
