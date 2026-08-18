package manager_test

import (
	"path/filepath"
	"testing"

	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/manager"
	"github.com/pactus-project/pactus/wallet/provider"
	"github.com/pactus-project/pactus/wallet/storage/jsonstorage"
	"github.com/pactus-project/pactus/wallet/vault"
	"github.com/stretchr/testify/require"
)

func TestWalletManager(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conf := &manager.Config{
		WalletsDir: util.TempDirPath(),
		ChainType:  genesis.Mainnet,
	}
	provider := provider.NewMockWalletProvider(ts.MockController())
	eventPipe := pipeline.New[any](t.Context())
	mgr, err := manager.NewManager(t.Context(), conf, provider, eventPipe)
	require.NoError(t, err)
	testWalletName := "test"

	t.Run("invalid wallet path", func(t *testing.T) {
		_, err := mgr.CreateWallet("../evil-path", "")
		require.ErrorContains(t, err, "illegal file path")
	})

	password := "password"

	t.Run("create wallet", func(t *testing.T) {
		_, err := mgr.CreateWallet(testWalletName, password)
		require.NoError(t, err)
	})

	t.Run("list wallet", func(t *testing.T) {
		wallets, err := mgr.ListWallets()
		require.NoError(t, err)
		require.Equal(t, []string{testWalletName}, wallets)
	})

	t.Run("get mnemonic", func(t *testing.T) {
		mnemonic, err := mgr.Mnemonic(testWalletName, password)
		require.NoError(t, err)
		require.NoError(t, vault.CheckMnemonic(mnemonic))
	})
}

func TestMigrateWallet(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	walletDir := util.TempDirPath()
	jsonPath := filepath.Join(walletDir, "legacy.json")

	mnemonic, err := wallet.GenerateMnemonic(128)
	require.NoError(t, err)

	vlt, err := vault.CreateVaultFromMnemonic(
		mnemonic, addresspath.CoinTypePactusMainnet, "password",
	)
	require.NoError(t, err)

	_, err = jsonstorage.Create(jsonPath, genesis.Mainnet, vlt)
	require.NoError(t, err)

	conf := &manager.Config{
		WalletsDir: walletDir,
		ChainType:  genesis.Mainnet,
	}
	provider := provider.NewMockWalletProvider(ts.MockController())
	eventPipe := pipeline.New[any](t.Context())

	mgr, err := manager.NewManager(t.Context(), conf, provider, eventPipe)
	require.NoError(t, err)

	t.Run("Migrate loaded legacy JSON wallet", func(t *testing.T) {
		info, err := mgr.WalletInfo("legacy.json")
		require.NoError(t, err)
		require.Equal(t, "JSON (legacy)", info.Driver)

		err = mgr.MigrateWallet("legacy.json")
		require.NoError(t, err)

		info, err = mgr.WalletInfo("legacy.json")
		require.NoError(t, err)
		require.Equal(t, "SQLite", info.Driver)

		wallets, err := mgr.ListWallets()
		require.NoError(t, err)
		require.Contains(t, wallets, "legacy.json")
	})

	t.Run("Migrate non-existent wallet", func(t *testing.T) {
		err = mgr.MigrateWallet("no-such-wallet.json")
		require.Error(t, err)
	})
}
