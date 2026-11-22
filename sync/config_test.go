package sync

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigBasicCheck(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr string
		updateFn    func(c *Config)
	}{
		{
			name:        "Invalid Session Timeout",
			expectedErr: "time: invalid duration",
			updateFn: func(c *Config) {
				c.SessionTimeoutStr = "INVALID-DURATION"
			},
		},
		{
			name:        "Zero Session Timeout",
			expectedErr: "session timeout should be positive",
			updateFn: func(c *Config) {
				c.SessionTimeoutStr = "0s"
			},
		},
		{
			name:        "Negative Session Timeout",
			expectedErr: "session timeout should be positive",
			updateFn: func(c *Config) {
				c.SessionTimeoutStr = "-5s"
			},
		},
		{
			name:        "Moniker is too long",
			expectedErr: "moniker is too long: 129 bytes, the maximum is 128 bytes",
			updateFn: func(c *Config) {
				c.Moniker = strings.Repeat("a", 129)
			},
		},
		{
			name:        "Moniker is too long in bytes",
			expectedErr: "moniker is too long: 130 bytes, the maximum is 128 bytes",
			updateFn: func(c *Config) {
				c.Moniker = strings.Repeat("\u00e9", 65)
			},
		},
		{
			name:        "Moniker is not valid UTF-8",
			expectedErr: "moniker is not a valid UTF-8 string",
			updateFn: func(c *Config) {
				c.Moniker = "alice\xff"
			},
		},
		{
			name: "Longest moniker",
			updateFn: func(c *Config) {
				c.Moniker = strings.Repeat("a", 128)
			},
		},
		{
			name: "Structured moniker",
			updateFn: func(c *Config) {
				c.Moniker = "name=alice;email=alice@example.com;web=https://example.com"
			},
		},
		{
			name:     "DefaultConfig",
			updateFn: func(*Config) {},
		},
	}

	for no, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := DefaultConfig()
			tt.updateFn(conf)
			if tt.expectedErr != "" {
				err := conf.BasicCheck()
				assert.ErrorContains(t, err, tt.expectedErr,
					"Expected error not matched for test %d-%s, expected: %s, got: %s", no, tt.name, tt.expectedErr, err)
			} else {
				err := conf.BasicCheck()
				require.NoError(t, err, "Expected no error for test %d-%s, get: %s", no, tt.name, err)
			}
		})
	}
}

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	require.NoError(t, c.BasicCheck())
	assert.Equal(t, 10*time.Second, c.SessionTimeout())
}
