package config

import (
	"strings"
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestSaveMainnetConfig(t *testing.T) {
	path := util.TempFilePath()
	assert.NoError(t, SaveMainnetConfig(path, 7))

	conf, err := LoadFromFile(path)
	assert.NoError(t, err)

	assert.NoError(t, conf.SanityCheck())
	assert.Equal(t, conf.Network.Name, "pactus")
}

func TestSaveTestnetConfig(t *testing.T) {
	path := util.TempFilePath()
	assert.NoError(t, SaveTestnetConfig(path, 7))

	conf, err := LoadFromFile(path)
	assert.NoError(t, err)

	assert.NoError(t, conf.SanityCheck())
	assert.Equal(t, conf.Network.Name, "pactus-testnet")
}

func TestSaveLocalnetConfig(t *testing.T) {
	path := util.TempFilePath()
	assert.NoError(t, SaveLocalnetConfig(path))

	conf, err := LoadFromFile(path)
	assert.NoError(t, err)

	assert.NoError(t, conf.SanityCheck())
	assert.Equal(t, conf.Network.Name, "pactus-localnet")
}

func TestLoadFromFile(t *testing.T) {
	path := util.TempFilePath()
	_, err := LoadFromFile(path)
	assert.Error(t, err, "not exists")

	assert.NoError(t, util.WriteFile(path, []byte(`foo = "bar"`)))
	_, err = LoadFromFile(path)
	assert.Error(t, err, "unknown field")
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
	exampleToml = strings.Replace(exampleToml, "\n", "", 1) // Remove the first \n
	defaultToml = strings.ReplaceAll(defaultToml, "\n\n", "\n")

	// fmt.Println(defaultToml)
	// fmt.Println(exampleToml)
	assert.Equal(t, defaultToml, exampleToml)
}
