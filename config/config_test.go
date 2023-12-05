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

	defConf := DefaultConfigMainnet()
	conf, err := LoadFromFile(path, true, defConf)
	assert.NoError(t, err)

	assert.NoError(t, conf.BasicCheck())
}

func TestSaveTestnetConfig(t *testing.T) {
	path := util.TempFilePath()
	assert.NoError(t, SaveTestnetConfig(path))

	defConf := DefaultConfigTestnet()
	conf, err := LoadFromFile(path, true, defConf)
	assert.NoError(t, err)

	assert.NoError(t, conf.BasicCheck())
	assert.Equal(t, conf.Network.NetworkName, "pactus-testnet-v2")
	assert.Equal(t, conf.Network.DefaultPort, 21777)
}

func TestSaveLocalnetConfig(t *testing.T) {
	path := util.TempFilePath()
	assert.NoError(t, SaveLocalnetConfig(path))

	defConf := DefaultConfigLocalnet()
	conf, err := LoadFromFile(path, true, defConf)
	assert.NoError(t, err)

	assert.NoError(t, conf.BasicCheck())
	assert.Empty(t, conf.Network.ListenAddrStrings)
	assert.Empty(t, conf.Network.RelayAddrStrings)
	assert.Equal(t, conf.Network.NetworkName, "pactus-localnet")
	assert.Equal(t, conf.Network.DefaultPort, 21666)
}

func TestLoadFromFile(t *testing.T) {
	path := util.TempFilePath()
	defConf := DefaultConfigTestnet()

	_, err := LoadFromFile(path, true, defConf)
	assert.Error(t, err, "not exists")

	assert.NoError(t, util.WriteFile(path, []byte(`foo = "bar"`)))
	_, err = LoadFromFile(path, true, defConf)
	assert.Error(t, err, "unknown field")

	conf, err := LoadFromFile(path, false, defConf)
	assert.NoError(t, err)
	assert.Equal(t, conf, defConf)
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

	defaultConf := DefaultConfigMainnet()
	defaultToml := string(defaultConf.toTOML())

	exampleToml = strings.ReplaceAll(exampleToml, "##", "")
	exampleToml = strings.ReplaceAll(exampleToml, "\r\n", "\n") // For Windows
	exampleToml = strings.ReplaceAll(exampleToml, "\n\n", "\n")
	defaultToml = strings.ReplaceAll(defaultToml, "\n\n", "\n")

	assert.Equal(t, defaultToml, exampleToml)
}

func TestNodeConfigBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("invalid reward addresses", func(t *testing.T) {
		conf := DefaultNodeConfig()
		conf.RewardAddresses = []string{
			ts.RandAccAddress().String(),
			"abcd",
		}

		assert.Error(t, conf.BasicCheck())
	})

	t.Run("validator address as reward address", func(t *testing.T) {
		conf := DefaultNodeConfig()
		conf.RewardAddresses = []string{
			ts.RandValAddress().String(),
		}

		assert.Error(t, conf.BasicCheck())
	})

	t.Run("ok", func(t *testing.T) {
		conf := DefaultNodeConfig()
		conf.RewardAddresses = []string{
			ts.RandAccAddress().String(),
			ts.RandAccAddress().String(),
		}

		assert.NoError(t, conf.BasicCheck())
	})

	t.Run("no reward addresses inside config, Ok", func(t *testing.T) {
		conf := DefaultNodeConfig()
		conf.RewardAddresses = []string{}

		assert.NoError(t, conf.BasicCheck())
	})
}
