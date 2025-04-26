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
	assert.Equal(t, conf, DefaultConfigMainnet())

	confData, _ := util.ReadFile(path)
	exampleData, _ := util.ReadFile("example_config.toml")
	assert.Equal(t, exampleData, confData)
}

func TestSaveTestnetConfig(t *testing.T) {
	path := util.TempFilePath()
	defConf := DefaultConfigTestnet()
	assert.NoError(t, defConf.Save(path))

	conf, err := LoadFromFile(path, true, defConf)
	assert.NoError(t, err)
	assert.Equal(t, conf, DefaultConfigTestnet())

	assert.NoError(t, conf.BasicCheck())
}

func TestDefaultConfig(t *testing.T) {
	conf := defaultConfig()

	assert.NoError(t, conf.BasicCheck())
	assert.Empty(t, conf.Network.ListenAddrStrings)
	assert.Zero(t, conf.Network.NetworkName)
	assert.Zero(t, conf.Network.DefaultPort)

	assert.False(t, conf.GRPC.Enable)
	assert.False(t, conf.Rest.Enable)
	assert.False(t, conf.HTTP.Enable)

	assert.Zero(t, conf.GRPC.Listen)
	assert.Zero(t, conf.Rest.Listen)
	assert.Zero(t, conf.HTTP.Listen)
}

func TestMainnetConfig(t *testing.T) {
	conf := DefaultConfigMainnet()

	assert.NoError(t, conf.BasicCheck())
	assert.Empty(t, conf.Network.ListenAddrStrings)
	assert.Equal(t, "pactus", conf.Network.NetworkName)
	assert.Equal(t, 21888, conf.Network.DefaultPort)

	assert.True(t, conf.GRPC.Enable)
	assert.False(t, conf.Rest.Enable)
	assert.False(t, conf.HTTP.Enable)

	assert.Equal(t, "127.0.0.1:50051", conf.GRPC.Listen)
	assert.Equal(t, "127.0.0.1:8080", conf.Rest.Listen)
	assert.Equal(t, "127.0.0.1:80", conf.HTTP.Listen)
}

func TestTestnetConfig(t *testing.T) {
	conf := DefaultConfigTestnet()

	assert.NoError(t, conf.BasicCheck())
	assert.Empty(t, conf.Network.ListenAddrStrings)
	assert.Equal(t, "pactus-testnet", conf.Network.NetworkName)
	assert.Equal(t, 21777, conf.Network.DefaultPort)

	assert.True(t, conf.GRPC.Enable)
	assert.True(t, conf.Rest.Enable)
	assert.False(t, conf.HTTP.Enable)

	assert.Equal(t, "[::]:50052", conf.GRPC.Listen)
	assert.Equal(t, "[::]:8080", conf.Rest.Listen)
	assert.Equal(t, "[::]:80", conf.HTTP.Listen)
}

func TestLocalnetConfig(t *testing.T) {
	conf := DefaultConfigLocalnet()

	assert.NoError(t, conf.BasicCheck())
	assert.Empty(t, conf.Network.ListenAddrStrings)
	assert.Equal(t, "pactus-localnet", conf.Network.NetworkName)
	assert.Equal(t, 0, conf.Network.DefaultPort)

	assert.True(t, conf.GRPC.Enable)
	assert.True(t, conf.Rest.Enable)
	assert.True(t, conf.HTTP.Enable)

	assert.Equal(t, "[::]:50052", conf.GRPC.Listen)
	assert.Equal(t, "[::]:8080", conf.Rest.Listen)
	assert.Equal(t, "[::]:0", conf.HTTP.Listen)
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
	assert.Equal(t, defConf, conf)
}

func TestExampleConfig(t *testing.T) {
	lines := strings.Split(string(exampleConfigBytes), "\n")
	exampleToml := ""
	for _, line := range lines {
		if !(strings.HasPrefix(line, "# ") ||
			strings.HasPrefix(line, "###") ||
			strings.HasPrefix(line, "  # ") ||
			strings.HasPrefix(line, "    # ") ||
			strings.HasPrefix(line, "      # ")) {
			exampleToml += line
			exampleToml += "\n"
		}
	}

	defaultConf := DefaultConfigMainnet()
	defaultToml := string(defaultConf.toTOML())

	exampleToml = strings.ReplaceAll(exampleToml, "\r\n", "\n") // For Windows
	exampleToml = strings.ReplaceAll(exampleToml, "\n\n", "\n")
	defaultToml = strings.ReplaceAll(defaultToml, "\n\n", "\n")

	defaultToml = strings.TrimSpace(defaultToml)
	exampleToml = strings.TrimSpace(exampleToml)

	assert.Equal(t, defaultToml, exampleToml)
}

func TestNodeConfigBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	randValAddr := ts.RandValAddress()

	tests := []struct {
		name        string
		expectedErr error
		updateFn    func(c *NodeConfig)
	}{
		{
			name: "Invalid reward addresses",
			expectedErr: NodeConfigError{
				Reason: "invalid reward address: invalid bech32 string length 4",
			},
			updateFn: func(c *NodeConfig) {
				c.RewardAddresses = []string{
					"abcd",
				}
			},
		},
		{
			name: "Validator address as reward address",
			expectedErr: NodeConfigError{
				Reason: "reward address is not an account address: " + randValAddr.String(),
			},
			updateFn: func(c *NodeConfig) {
				c.RewardAddresses = []string{
					randValAddr.String(),
				}
			},
		},
		{
			name: "Two rewards addresses",
			updateFn: func(c *NodeConfig) {
				c.RewardAddresses = []string{
					ts.RandAccAddress().String(),
					ts.RandAccAddress().String(),
				}
			},
		},
		{
			name: "No reward address",
			updateFn: func(c *NodeConfig) {
				c.RewardAddresses = []string{}
			},
		},
		{
			name:     "DefaultConfig",
			updateFn: func(*NodeConfig) {},
		},
	}

	for no, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := DefaultNodeConfig()
			tt.updateFn(conf)
			if tt.expectedErr != nil {
				err := conf.BasicCheck()
				assert.ErrorIs(t, tt.expectedErr, err,
					"Expected error not matched for test %d-%s, expected: %s, got: %s",
					no, tt.name, tt.expectedErr, err)
			} else {
				err := conf.BasicCheck()
				assert.NoError(t, err, "Expected no error for test %d-%s, get: %s", no, tt.name, err)
			}
		})
	}
}
