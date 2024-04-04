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
	assert.NoError(t, SaveMainnetConfig(path))

	defConf := DefaultConfigMainnet()
	conf, err := LoadFromFile(path, true, defConf)
	assert.NoError(t, err)

	assert.NoError(t, conf.BasicCheck())
	assert.Equal(t, DefaultConfigMainnet(), conf)

	confData, _ := util.ReadFile(path)
	exampleData, _ := util.ReadFile("example_config.toml")
	assert.Equal(t, confData, exampleData)
}

func TestSaveTestnetConfig(t *testing.T) {
	path := util.TempFilePath()
	defConf := DefaultConfigTestnet()
	assert.NoError(t, defConf.Save(path))

	conf, err := LoadFromFile(path, true, defConf)
	assert.NoError(t, err)
	assert.Equal(t, DefaultConfigTestnet(), conf)

	assert.NoError(t, conf.BasicCheck())
}

func TestDefaultConfig(t *testing.T) {
	conf := defaultConfig()

	assert.NoError(t, conf.BasicCheck())
	assert.Empty(t, conf.Network.ListenAddrStrings)
	assert.Equal(t, conf.Network.NetworkName, "")
	assert.Equal(t, conf.Network.DefaultPort, 0)

	assert.False(t, conf.GRPC.Enable)
	assert.False(t, conf.GRPC.Gateway.Enable)
	assert.False(t, conf.HTTP.Enable)
	assert.False(t, conf.Nanomsg.Enable)

	assert.Equal(t, conf.GRPC.Listen, "")
	assert.Equal(t, conf.GRPC.Gateway.Listen, "")
	assert.Equal(t, conf.HTTP.Listen, "")
	assert.Equal(t, conf.Nanomsg.Listen, "")
}

func TestMainnetConfig(t *testing.T) {
	conf := DefaultConfigMainnet()

	assert.NoError(t, conf.BasicCheck())
	assert.Empty(t, conf.Network.ListenAddrStrings)
	assert.Equal(t, conf.Network.NetworkName, "pactus")
	assert.Equal(t, conf.Network.DefaultPort, 21888)

	assert.True(t, conf.GRPC.Enable)
	assert.False(t, conf.GRPC.Gateway.Enable)
	assert.False(t, conf.HTTP.Enable)
	assert.False(t, conf.Nanomsg.Enable)

	assert.Equal(t, conf.GRPC.Listen, "127.0.0.1:50051")
	assert.Equal(t, conf.GRPC.Gateway.Listen, "127.0.0.1:8080")
	assert.Equal(t, conf.HTTP.Listen, "127.0.0.1:80")
	assert.Equal(t, conf.Nanomsg.Listen, "tcp://127.0.0.1:40899")
}

func TestTestnetConfig(t *testing.T) {
	conf := DefaultConfigTestnet()

	assert.NoError(t, conf.BasicCheck())
	assert.Empty(t, conf.Network.ListenAddrStrings)
	assert.Equal(t, conf.Network.NetworkName, "pactus-testnet")
	assert.Equal(t, conf.Network.DefaultPort, 21777)

	assert.True(t, conf.GRPC.Enable)
	assert.True(t, conf.GRPC.Gateway.Enable)
	assert.False(t, conf.HTTP.Enable)
	assert.False(t, conf.Nanomsg.Enable)

	assert.Equal(t, conf.GRPC.Listen, "[::]:50052")
	assert.Equal(t, conf.GRPC.Gateway.Listen, "[::]:8080")
	assert.Equal(t, conf.HTTP.Listen, "[::]:80")
	assert.Equal(t, conf.Nanomsg.Listen, "tcp://[::]:40799")
}

func TestLocalnetConfig(t *testing.T) {
	conf := DefaultConfigLocalnet()

	assert.NoError(t, conf.BasicCheck())
	assert.Empty(t, conf.Network.ListenAddrStrings)
	assert.Equal(t, conf.Network.NetworkName, "pactus-localnet")
	assert.Equal(t, conf.Network.DefaultPort, 0)

	assert.True(t, conf.GRPC.Enable)
	assert.True(t, conf.GRPC.Gateway.Enable)
	assert.True(t, conf.HTTP.Enable)
	assert.True(t, conf.Nanomsg.Enable)

	assert.Equal(t, conf.GRPC.Listen, "[::]:50052")
	assert.Equal(t, conf.GRPC.Gateway.Listen, "[::]:8080")
	assert.Equal(t, conf.HTTP.Listen, "[::]:0")
	assert.Equal(t, conf.Nanomsg.Listen, "tcp://[::]:40799")
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
