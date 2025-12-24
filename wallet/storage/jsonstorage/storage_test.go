package jsonstorage

import (
	"testing"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	tempPath := util.TempFilePath()

	vlt, err := vault.CreateVaultFromMnemonic(testMnemonic, 21888)
	require.NoError(t, err)

	strg, err := Create(tempPath, genesis.Mainnet, *vlt)
	require.NoError(t, err)

	assert.Equal(t, VersionLatest, strg.WalletInfo().Version)
	assert.Equal(t, genesis.Mainnet, strg.WalletInfo().Network)
	assert.Equal(t, vlt, strg.Vault())
}

func TestOpen(t *testing.T) {
}

func TestOpenNeuterWallet(t *testing.T) {
	data, err := util.ReadFile("./testdata/neuter_wallet")
	require.NoError(t, err)

	tempPath := util.TempFilePath()
	err = util.WriteFile(tempPath, data)
	require.NoError(t, err)

	strg, err := Open(tempPath)
	require.NoError(t, err)

	assert.False(t, strg.Vault().IsNeutered())
}
