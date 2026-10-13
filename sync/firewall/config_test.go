package firewall

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigBasicCheck(t *testing.T) {
	tests := []struct {
		name       string
		bannedNets []string
		valid      bool
	}{
		{"no banned networks", nil, true},
		{"valid banned networks", []string{"10.0.0.0/8", "192.168.1.0/24", "2001:db8::/32"}, true},
		{"IP address without a mask", []string{"10.0.0.1"}, false},
		{"not an IP address", []string{"not-a-network"}, false},
		{"one invalid network among valid ones", []string{"10.0.0.0/8", "bad"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := DefaultConfig()
			conf.BannedNets = tt.bannedNets

			err := conf.BasicCheck()
			if tt.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}

	t.Run("default config is valid", func(t *testing.T) {
		require.NoError(t, DefaultConfig().BasicCheck())
	})
}

func TestNewFirewallWithInvalidConfig(t *testing.T) {
	conf := DefaultConfig()
	conf.BannedNets = []string{"not-a-network"}

	fw, err := NewFirewall(t.Context(), conf, nil, nil, nil)

	require.Error(t, err)
	assert.Nil(t, fw)
}

func TestConfigErrorMessage(t *testing.T) {
	assert.Equal(t, "invalid config", ConfigError{Reason: "invalid config"}.Error())
}
