package txpool

import (
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	conf := DefaultConfig()
	assert.NoError(t, conf.BasicCheck())

	assert.Equal(t, 600, conf.transferPoolSize())
	assert.Equal(t, 100, conf.bondPoolSize())
	assert.Equal(t, 100, conf.unbondPoolSize())
	assert.Equal(t, 100, conf.withdrawPoolSize())
	assert.Equal(t, 100, conf.sortitionPoolSize())
	assert.Equal(t, amount.Amount(0.1e8), conf.fixedFee())

	assert.Equal(t,
		conf.transferPoolSize()+
			conf.bondPoolSize()+
			conf.unbondPoolSize()+
			conf.withdrawPoolSize()+
			conf.sortitionPoolSize(), conf.MaxSize)
}

func TestConfigBasicCheck(t *testing.T) {
	tests := []struct {
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
				Reason: "dailyLimit should be positive",
			},
			updateFn: func(c *Config) {
				c.Fee.DailyLimit = 0
			},
		},
		{
			name: "Negative DailyLimit",
			expectedErr: ConfigError{
				Reason: "dailyLimit should be positive",
			},
			updateFn: func(c *Config) {
				c.Fee.DailyLimit = -1
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
