package sqlitestorage

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConcurrentAccess(t *testing.T) {
	td := setup(t)

	// The exclusive locking mode is set on one connection.
	// Concurrent calls must not open a second connection that can never get the lock.
	infos := make([]*types.TransactionInfo, 8)
	for i := range infos {
		infos[i] = td.RandomTransactionInfo(t, td.RandomDirection(t))
	}

	errs := make(chan error, 2*len(infos))
	var wg sync.WaitGroup
	for _, info := range infos {
		wg.Go(func() {
			errs <- td.storage.InsertTransaction(info)
		})
		wg.Go(func() {
			_, err := td.storage.GetPendingTransactions()
			errs <- err
		})
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}

	pendings, err := td.storage.GetPendingTransactions()
	require.NoError(t, err)
	assert.Len(t, pendings, len(infos))
}

func TestDatabaseFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permissions are not supported on Windows")
	}

	td := setup(t)

	stat, err := os.Stat(dbPath(td.storage.path))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), stat.Mode().Perm())

	clonePath := util.TempDirPath()
	cloned, err := td.storage.Clone(clonePath)
	require.NoError(t, err)
	require.NoError(t, cloned.Close())

	stat, err = os.Stat(dbPath(clonePath))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), stat.Mode().Perm())
}

func TestCloneToPathWithQuote(t *testing.T) {
	td := setup(t)

	clonePath := filepath.Join(util.TempDirPath(), "o'brien")
	cloned, err := td.storage.Clone(clonePath)
	require.NoError(t, err)
	require.NoError(t, cloned.Close())
}

func TestNeuteredCloneHasNoKeyStore(t *testing.T) {
	td := setup(t)

	keyStore := td.storage.Vault().KeyStore
	require.NotEmpty(t, keyStore)

	clonePath := util.TempDirPath()
	cloned, err := td.storage.Clone(clonePath)
	require.NoError(t, err)

	vlt := cloned.Vault()
	vlt.Neuter()
	require.NoError(t, cloned.UpdateVault(vlt))
	require.NoError(t, cloned.Close())

	// The key store is encrypted, but it must not be in a wallet that is given to others.
	// Even a part of it must not stay in the free pages of the file or in the write-ahead log.
	entries, err := os.ReadDir(clonePath)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	for _, entry := range entries {
		data, err := os.ReadFile(filepath.Join(clonePath, entry.Name()))
		require.NoError(t, err)
		for i := 0; i+16 <= len(keyStore); i += 16 {
			assert.NotContains(t, string(data), keyStore[i:i+16],
				"file %s contains a part of the key store", entry.Name())
		}
	}

	reopened, err := Open(t.Context(), clonePath)
	require.NoError(t, err)
	assert.True(t, reopened.WalletInfo().Neutered)
	assert.Empty(t, reopened.Vault().KeyStore)
	require.NoError(t, reopened.Close())
}
