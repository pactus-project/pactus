package sync

import (
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
