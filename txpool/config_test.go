package txpool

import (
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.BasicCheck())

	assert.Equal(t, 600, c.transferPoolSize())
	assert.Equal(t, 100, c.bondPoolSize())
	assert.Equal(t, 100, c.unbondPoolSize())
	assert.Equal(t, 100, c.withdrawPoolSize())
	assert.Equal(t, 100, c.sortitionPoolSize())
	assert.Equal(t, amount.Amount(0.1e8), c.minFee())

	assert.Equal(t,
		c.transferPoolSize()+
			c.bondPoolSize()+
			c.unbondPoolSize()+
			c.withdrawPoolSize()+
			c.sortitionPoolSize(), c.MaxSize)
}

func TestConfigBasicCheck(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr error
		updateFn    func(c *Config)
	}{
		{
			name: "Invalid MaxSize",
			expectedErr: ConfigError{
				Reason: "maxSize can't be less than 10",
			},
			updateFn: func(c *Config) {
				c.MaxSize = 0
			},
		},
		{
			name: "Invalid MaxSize",
			expectedErr: ConfigError{
				Reason: "maxSize can't be less than 10",
			},
			updateFn: func(c *Config) {
				c.MaxSize = 9
			},
		},
		{
			name: "Invalid DailyLimit",
			expectedErr: ConfigError{
				Reason: "dailyLimit can't be zero",
			},
			updateFn: func(c *Config) {
				c.Fee.DailyLimit = 0
			},
		},
		{
			name: "Valid Config",
			updateFn: func(c *Config) {
				c.MaxSize = 100
				c.Fee = &FeeConfig{
					FixedFee:   0.01,
					DailyLimit: 280,
					UnitPrice:  0,
				}
			},
		},
		{
			name:     "DefaultConfig",
			updateFn: func(*Config) {},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := DefaultConfig()
			tc.updateFn(conf)
			if tc.expectedErr != nil {
				err := conf.BasicCheck()
				assert.ErrorIs(t, tc.expectedErr, err,
					"Expected error not matched for test %d-%s, expected: %s, got: %s", i, tc.name, tc.expectedErr, err)
			} else {
				err := conf.BasicCheck()
				assert.NoError(t, err, "Expected no error for test %d-%s, get: %s", i, tc.name, err)
			}
		})
	}
}
