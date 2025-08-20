package consensusv2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigBasicCheck(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
		updateFn    func(c *Config)
	}{
		{
			name: "Invalid ChangeProposerDelta",
			expectedErr: ConfigError{
				Reason: "change proposer delta must be greater than zero",
			},
			updateFn: func(c *Config) {
				c.ChangeProposerDelta = 0
			},
		},
		{
			name: "Invalid ChangeProposerTimeout",
			expectedErr: ConfigError{
				Reason: "change proposer timeout must be greater than zero",
			},
			updateFn: func(c *Config) {
				c.ChangeProposerTimeout = -1 * time.Second
			},
		},
		{
			name: "Invalid MinimumAvailabilityScore",
			expectedErr: ConfigError{
				Reason: "minimum availability score can't be negative or more than 1",
			},
			updateFn: func(c *Config) {
				c.MinimumAvailabilityScore = 1.5
			},
		},
		{
			name: "Invalid MinimumAvailabilityScore - Negative",
			expectedErr: ConfigError{
				Reason: "minimum availability score can't be negative or more than 1",
			},
			updateFn: func(c *Config) {
				c.MinimumAvailabilityScore = -0.8
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
			if tt.expectedErr != nil {
				err := conf.BasicCheck()
				assert.ErrorIs(t, tt.expectedErr, err,
					"Expected error not matched for test %d-%s, expected: %s, got: %s", no, tt.name, tt.expectedErr, err)
			} else {
				err := conf.BasicCheck()
				assert.NoError(t, err, "Expected no error for test %d-%s, get: %s", no, tt.name, err)
			}
		})
	}
}

func TestConfigCalculateChangeProposerTimeout(t *testing.T) {
	c := DefaultConfig()

	assert.Equal(t, c.ChangeProposerTimeout, c.CalculateChangeProposerTimeout(0))
	assert.Equal(t, c.ChangeProposerTimeout+c.ChangeProposerDelta, c.CalculateChangeProposerTimeout(1))
	assert.Equal(t, c.ChangeProposerTimeout+(4*c.ChangeProposerDelta), c.CalculateChangeProposerTimeout(4))
}
