//go111:build gtk

package model

import (
	"testing"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigModel(t *testing.T) {
	workingDir := util.TempDirPath()
	gen := genesis.MainnetGenesis()
	defConf := config.DefaultConfigForChain(genesis.Mainnet)

	require.NoError(t, gen.SaveToFile(cmd.PactusGenesisPath(workingDir)))
	require.NoError(t, defConf.Save(cmd.PactusConfigPath(workingDir)))

	model, err := NewConfigModel(workingDir, true)
	require.NoError(t, err)

	t.Run("Check default config", func(t *testing.T) {
		defTOML, err := defConf.ToTOML()
		require.NoError(t, err)

		assert.Equal(t, defTOML, model.DefaultTOML())
		assert.False(t, model.IsDirty(model.SavedContent()))
		assert.True(t, model.IsDirty(model.SavedContent()+"\n# edit"))
	})

	t.Run("Invalid input", func(t *testing.T) {
		badTom := "not valid {{{ toml"
		err = model.Validate(badTom)
		require.Error(t, err)

		err = model.Save(badTom)
		require.Error(t, err)
	})

	t.Run("Save content", func(t *testing.T) {
		edited := model.SavedContent() + "\n"
		require.NoError(t, model.Save(edited))
		assert.False(t, model.IsDirty(edited))
	})
}
