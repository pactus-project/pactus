package main

import (
	"path/filepath"
	"testing"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestorePhrase(t *testing.T) {
	phrase := "abandon ability able about above absent absorb abstract absurd abuse access accident"

	t.Run("the phrase of the command line is used as it is", func(t *testing.T) {
		got, err := restorePhrase(phrase, "")
		require.NoError(t, err)
		assert.Equal(t, phrase, got)
	})

	t.Run("the phrase of a file is read without its line breaks and extra spaces", func(t *testing.T) {
		path := filepath.Join(util.TempDirPath(), "seed.txt")
		require.NoError(t, util.WriteFile(path, []byte("  abandon ability able about above absent\n"+
			"absorb abstract  absurd abuse access accident\n\n")))

		got, err := restorePhrase("", path)
		require.NoError(t, err)
		assert.Equal(t, phrase, got)
	})

	t.Run("both options can't be used together", func(t *testing.T) {
		path := filepath.Join(util.TempDirPath(), "seed.txt")
		require.NoError(t, util.WriteFile(path, []byte(phrase)))

		_, err := restorePhrase(phrase, path)
		require.Error(t, err)
	})

	t.Run("a missing file is an error", func(t *testing.T) {
		_, err := restorePhrase("", filepath.Join(util.TempDirPath(), "missing.txt"))
		require.Error(t, err)
	})
}

// The init command restores the wallet from the seed phrase of a file, without any prompt
// when the password and the number of validators are given.
func TestInitRestoresFromFile(t *testing.T) {
	phrase, err := wallet.GenerateMnemonic(128)
	require.NoError(t, err)

	seedPath := filepath.Join(util.TempDirPath(), "seed.txt")
	require.NoError(t, util.WriteFile(seedPath, []byte(phrase+"\n")))
	workingDir := util.TempDirPath()

	rootCmd := &cobra.Command{Use: "pactus-daemon"}
	buildInitCmd(rootCmd)
	rootCmd.SetArgs([]string{
		"init", "--working-dir", workingDir, "--restore-from-file", seedPath,
		"--password", "secret", "--val-num", "2",
	})
	require.NoError(t, rootCmd.Execute())

	// The command keeps the wallet open, as the process exits after it.
	assert.True(t, util.PathExists(cmd.PactusDefaultWalletPath(workingDir)), "the wallet is not created")
	assert.True(t, util.PathExists(cmd.PactusConfigPath(workingDir)), "the config is not created")
	assert.True(t, util.PathExists(cmd.PactusGenesisPath(workingDir)), "the genesis is not created")
}
