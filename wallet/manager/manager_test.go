package manager_test

import (
	"testing"

	"github.com/ezex-io/gopkg/pipeline"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet/manager"
	"github.com/pactus-project/pactus/wallet/provider"
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

	t.Run("create wallet", func(t *testing.T) {
		_, err := mgr.CreateWallet(testWalletName, "")
		require.NoError(t, err)
	})

	t.Run("list wallet", func(t *testing.T) {
		wallets, err := mgr.ListWallets()
		require.NoError(t, err)
		require.Equal(t, []string{testWalletName}, wallets)
	})

	t.Run("get mnemonic", func(t *testing.T) {
		mnemonic, err := mgr.Mnemonic(testWalletName, "")
		require.NoError(t, err)
		require.NoError(t, vault.CheckMnemonic(mnemonic))
	})
}
