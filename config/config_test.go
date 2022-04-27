package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)


func TestSaveMainnetConfig(t *testing.T) {
	path := util.TempFilePath()
	rewardAddr := crypto.GenerateTestAddress()
	assert.NoError(t, SaveMainnetConfig(path, rewardAddr.String()))

	conf, err := LoadFromFile(path)
	assert.NoError(t, err)

	assert.NoError(t, conf.SanityCheck())
	assert.Equal(t, conf.State.RewardAddress, rewardAddr.String())
	assert.Equal(t, conf.Network.Name, "zarb")
}

func TestSaveTestnetConfig(t *testing.T) {
	path := util.TempFilePath()
	rewardAddr := crypto.GenerateTestAddress()
	assert.NoError(t, SaveTestnetConfig(path, rewardAddr.String()))

	conf, err := LoadFromFile(path)
	assert.NoError(t, err)

	assert.NoError(t, conf.SanityCheck())
	assert.Equal(t, conf.State.RewardAddress, rewardAddr.String())
	assert.Equal(t, conf.Network.Name, "zarb-testnet")
}

func TestSaveLocalnetConfig(t *testing.T) {
	path := util.TempFilePath()
	rewardAddr := crypto.GenerateTestAddress()
	assert.NoError(t, SaveLocalnetConfig(path, rewardAddr.String()))

	conf, err := LoadFromFile(path)
	assert.NoError(t, err)

	assert.NoError(t, conf.SanityCheck())
	assert.Equal(t, conf.State.RewardAddress, rewardAddr.String())
	assert.Equal(t, conf.Network.Name, "zarb-localnet")
}

func TestLoadFromFile(t *testing.T) {
	path := util.TempFilePath()
	_, err := LoadFromFile(path)
	assert.Error(t, err, "not exists")

	util.WriteFile(path, []byte(`foo = "bar"`))
	_, err = LoadFromFile(path)
	assert.Error(t, err, "unknown field")
}

func TestExampleConfig(t *testing.T) {
	lines := strings.Split(string(exampleConfigBytes), "\n")
	exampleToml := ""
	for _, line := range lines {
		if !(strings.HasPrefix(line, "# ") ||
			strings.HasPrefix(line, "  # ") ||
			strings.HasPrefix(line, "    # ")) {
			exampleToml += line
			exampleToml += "\n"
		}
	}

	defaultConf := DefaultConfig()
	defaultToml := string(defaultConf.toTOML())

	exampleToml = strings.ReplaceAll(exampleToml, "##", "")
	exampleToml = strings.ReplaceAll(exampleToml, "\n\n", "\n")
	defaultToml = strings.ReplaceAll(defaultToml, "\n\n", "\n")

	// fmt.Println(defaultToml)
	// fmt.Println(exampleToml)
	assert.Equal(t, defaultToml, exampleToml)
}
