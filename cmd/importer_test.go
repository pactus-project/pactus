package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewImporterTempDirFailure(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("the cache directory is taken from XDG_CACHE_HOME only on Linux")
	}

	// A cache directory that cannot be created makes the importer fall back
	// to the temporary directory, which is missing too.
	storeDir := filepath.Join(util.TempDirPath(), "store")
	notDir := filepath.Join(util.TempDirPath(), "file")
	require.NoError(t, os.WriteFile(notDir, nil, 0o600))
	t.Setenv("XDG_CACHE_HOME", filepath.Join(notDir, "cache"))
	t.Setenv("TMPDIR", filepath.Join(util.TempDirPath(), "missing"))

	importer, err := NewImporter(genesis.Mainnet, "https://example.com", storeDir)
	require.Error(t, err)
	assert.Nil(t, importer)
}

func TestNewImporterUsesCacheDir(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("the cache directory is taken from XDG_CACHE_HOME only on Linux")
	}

	cacheDir := filepath.Join(util.TempDirPath(), "cache")
	t.Setenv("XDG_CACHE_HOME", cacheDir)

	storeDir := filepath.Join(util.TempDirPath(), "store")
	importer, err := NewImporter(genesis.Mainnet, "https://example.com", storeDir)
	require.NoError(t, err)
	defer func() { _ = importer.Cleanup() }()

	assert.Equal(t, cacheDir, filepath.Dir(importer.tempDir))
}
