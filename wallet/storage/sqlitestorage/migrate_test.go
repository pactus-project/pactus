package sqlitestorage

import (
	"testing"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet/storage/jsonstorage"
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
	password := ts.RandString(12)
	vlt, err := vault.CreateVaultFromMnemonic(testMnemonic,
		addresspath.CoinTypePactusMainnet, password, opts...)
	require.NoError(t, err)

	jsonPath := util.TempFilePath()
	jsonStrg, err := jsonstorage.Create(jsonPath, genesis.Mainnet, vlt)
	require.NoError(t, err)

	// Add some addresses to the legacy wallet.
	addr1, _ := vlt.NewValidatorAddress("addr 1")
	addr2, _ := vlt.NewBLSAccountAddress("addr 2")
	addr3, _ := vlt.NewEd25519AccountAddress("addr 3", password)
	addr4, _ := vlt.NewSecp256k1AccountAddress("addr 4", password)

	_ = jsonStrg.InsertAddress(addr1)
	_ = jsonStrg.InsertAddress(addr2)
	_ = jsonStrg.InsertAddress(addr3)
	_ = jsonStrg.InsertAddress(addr4)

	defFee := ts.RandFee()
	_ = jsonStrg.SetDefaultFee(defFee)

	_ = jsonStrg.UpdateVault(vlt)

	t.Run("Migration to the same path (invalid)", func(t *testing.T) {
		err := Migrate(t.Context(), jsonPath, jsonStrg)
		require.Error(t, err)
	})

	t.Run("Successful migration", func(t *testing.T) {
		path := util.TempFilePath()
		err := Migrate(t.Context(), path, jsonStrg)
		require.NoError(t, err)

		strg, err := Open(t.Context(), path)
		require.NoError(t, err)

		assert.Equal(t, vlt, strg.Vault())
		assert.Equal(t, "SQLite", strg.WalletInfo().Driver)
		assert.Equal(t, defFee, strg.WalletInfo().DefaultFee)
		assert.Equal(t, genesis.Mainnet, strg.WalletInfo().Network)
		assert.Equal(t, true, strg.WalletInfo().Encrypted)
		assert.Equal(t, path, strg.WalletInfo().Path)
		assert.Equal(t, VersionLatest, strg.WalletInfo().Version)

		assert.True(t, strg.HasAddress(addr1.Address))
		assert.True(t, strg.HasAddress(addr2.Address))
		assert.True(t, strg.HasAddress(addr3.Address))
		assert.True(t, strg.HasAddress(addr4.Address))
	})

	t.Run("Neutered migration", func(t *testing.T) {
		vlt.Neuter()
		jsonStrg.UpdateVault(vlt)

		path := util.TempFilePath()
		err := Migrate(t.Context(), path, jsonStrg)
		require.NoError(t, err)

		strg, err := Open(t.Context(), path)
		require.NoError(t, err)

		assert.Equal(t, vlt, strg.Vault())
		assert.Equal(t, false, strg.WalletInfo().Encrypted)
		assert.Equal(t, true, strg.WalletInfo().Neutered)
	})
}
