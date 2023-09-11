package config

import (
	"strings"
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestSaveMainnetConfig(t *testing.T) {
	path := util.TempFilePath()
	assert.NoError(t, SaveMainnetConfig(path, 7))

	conf, err := LoadFromFile(path, true)
	assert.NoError(t, err)

	assert.NoError(t, conf.BasicCheck())
}

func TestSaveTestnetConfig(t *testing.T) {
	path := util.TempFilePath()
	assert.NoError(t, SaveTestnetConfig(path, 7))

	conf, err := LoadFromFile(path, true)
	assert.NoError(t, err)

	assert.NoError(t, conf.BasicCheck())
}

func TestSaveLocalnetConfig(t *testing.T) {
	path := util.TempFilePath()
	assert.NoError(t, SaveLocalnetConfig(path, 4))

	conf, err := LoadFromFile(path, true)
	assert.NoError(t, err)

	assert.NoError(t, conf.BasicCheck())
	assert.Empty(t, conf.Network.Listens)
	assert.Empty(t, conf.Network.RelayAddrs)
}

func TestLoadFromFile(t *testing.T) {
	path := util.TempFilePath()
	_, err := LoadFromFile(path, true)
	assert.Error(t, err, "not exists")

	assert.NoError(t, util.WriteFile(path, []byte(`foo = "bar"`)))
	_, err = LoadFromFile(path, true)
	assert.Error(t, err, "unknown field")

	_, err = LoadFromFile(path, false)
	assert.NoError(t, err)
}

func TestExampleConfig(t *testing.T) {
	lines := strings.Split(string(exampleConfigBytes), "\n")
	exampleToml := ""
	for _, line := range lines {
		if !(strings.HasPrefix(line, "# ") ||
			strings.HasPrefix(line, "###") ||
			strings.HasPrefix(line, "  # ") ||
			strings.HasPrefix(line, "    # ")) {
			exampleToml += line
			exampleToml += "\n"
		}
	}

	defaultConf := DefaultConfig()
	defaultToml := string(defaultConf.toTOML())

	exampleToml = strings.ReplaceAll(exampleToml, "%num_validators%", "7")
	exampleToml = strings.ReplaceAll(exampleToml, "##", "")
	exampleToml = strings.ReplaceAll(exampleToml, "\r\n", "\n") // For Windows
	exampleToml = strings.ReplaceAll(exampleToml, "\n\n", "\n")
	defaultToml = strings.ReplaceAll(defaultToml, "\n\n", "\n")

	// fmt.Println(defaultToml)
	// fmt.Println(exampleToml)
	assert.Equal(t, defaultToml, exampleToml)
}

func TestNodeConfigBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("invalid number of validators", func(t *testing.T) {
		conf := DefaultNodeConfig()
		conf.NumValidators = 0

		assert.Error(t, conf.BasicCheck())
	})

	t.Run("invalid number of reward addresses", func(t *testing.T) {
		conf := DefaultNodeConfig()
		conf.RewardAddresses = []string{
			ts.RandAddress().String(),
		}

		assert.Error(t, conf.BasicCheck())
	})

	t.Run("invalid reward addresses", func(t *testing.T) {
		conf := DefaultNodeConfig()
		conf.NumValidators = 2
		conf.RewardAddresses = []string{
			ts.RandAddress().String(),
			"abcd",
		}

		assert.Error(t, conf.BasicCheck())
	})

	t.Run("ok", func(t *testing.T) {
		conf := DefaultNodeConfig()
		conf.NumValidators = 2
		conf.RewardAddresses = []string{
			ts.RandAddress().String(),
			ts.RandAddress().String(),
		}

		assert.NoError(t, conf.BasicCheck())
	})
}
