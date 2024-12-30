package zmq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.NotNil(t, cfg, "DefaultConfig should not return nil")
	assert.Equal(t, "", cfg.ZmqPubBlockInfo, "ZmqPubBlockInfo should be empty")
	assert.Equal(t, "", cfg.ZmqPubTxInfo, "ZmqPubTxInfo should be empty")
	assert.Equal(t, "", cfg.ZmqPubRawBlock, "ZmqPubRawBlock should be empty")
	assert.Equal(t, "", cfg.ZmqPubRawTx, "ZmqPubRawTx should be empty")
	assert.Equal(t, 1000, cfg.ZmqPubHWM, "ZmqPubHWM should default to 1000")
}

func TestBasicCheck(t *testing.T) {
	testCases := []struct {
		name      string
		config    *Config
		expectErr bool
	}{
		{
			name: "Valid configuration",
			config: &Config{
				ZmqPubBlockInfo: "tcp://127.0.0.1:28332",
				ZmqPubTxInfo:    "tcp://127.0.0.1:28333",
				ZmqPubRawBlock:  "tcp://127.0.0.1:28334",
				ZmqPubRawTx:     "tcp://127.0.0.1:28335",
				ZmqPubHWM:       1000,
			},
			expectErr: false,
		},
		{
			name: "Invalid scheme",
			config: &Config{
				ZmqPubBlockInfo: "udp://127.0.0.1:28332",
			},
			expectErr: true,
		},
		{
			name: "Missing port",
			config: &Config{
				ZmqPubBlockInfo: "tcp://127.0.0.1",
			},
			expectErr: true,
		},
		{
			name: "Empty host",
			config: &Config{
				ZmqPubBlockInfo: "tcp://:28332",
			},
			expectErr: true,
		},
		{
			name: "Negative ZmqPubHWM",
			config: &Config{
				ZmqPubHWM: -1,
			},
			expectErr: true,
		},
		{
			name:      "Empty configuration",
			config:    DefaultConfig(),
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.BasicCheck()
			if tc.expectErr {
				assert.Error(t, err, "BasicCheck should return an error")
			} else {
				assert.NoError(t, err, "BasicCheck should not return an error")
			}
		})
	}
}
