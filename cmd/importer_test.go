package cmd

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewImporterTempDirFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the temporary directory is not taken from TMPDIR on Windows")
	}

	storeDir := filepath.Join(util.TempDirPath(), "store")
	t.Setenv("TMPDIR", filepath.Join(util.TempDirPath(), "missing"))

	importer, err := NewImporter(genesis.Mainnet, "https://example.com", storeDir)
	require.Error(t, err)
	assert.Nil(t, importer)
}
