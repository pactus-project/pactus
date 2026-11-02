package store

import (
	"runtime"
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigBasicCheck(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
		updateFn    func(c *Config)
	}{
		{
			name: "Invalid Path",
			expectedErr: ConfigError{
				Reason: "path is invalid or inaccessible: /invalid:path/\x00*folder?\\CON",
			},
			updateFn: func(c *Config) {
				c.Path = "/invalid:path/\x00*folder?\\CON"
			},
		},
		{
			name: "Invalid TxCacheWindow",
			expectedErr: ConfigError{
				Reason: "cache window set to zero",
			},
			updateFn: func(c *Config) {
				c.TxCacheWindow = 0
			},
		},
		{
			name: "Invalid AccountCacheSize",
			expectedErr: ConfigError{
				Reason: "cache size must be positive",
			},
			updateFn: func(c *Config) {
				c.AccountCacheSize = 0
			},
		},
		{
			name: "Negative AccountCacheSize",
			expectedErr: ConfigError{
				Reason: "cache size must be positive",
			},
			updateFn: func(c *Config) {
				c.AccountCacheSize = -1
			},
		},
		{
			name: "Negative PublicKeyCacheSize",
			expectedErr: ConfigError{
				Reason: "cache size must be positive",
			},
			updateFn: func(c *Config) {
				c.PublicKeyCacheSize = -1
			},
		},
		{
			name: "Invalid RetentionDays",
			expectedErr: ConfigError{
				Reason: "retention days can't be less than 10 days",
			},
			updateFn: func(c *Config) {
				c.RetentionDays = 1
			},
		},
		{
			name: "Too large RetentionDays",
			expectedErr: ConfigError{
				Reason: "retention days can't be more than 497102 days",
			},
			updateFn: func(c *Config) {
				c.RetentionDays = 497103
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
				require.ErrorIs(t, err, tt.expectedErr,
					"Expected error not matched for test %d-%s, expected: %s, got: %s", no, tt.name, tt.expectedErr, err)
			} else {
				err := conf.BasicCheck()
				require.NoError(t, err, "Expected no error for test %d-%s, get: %s", no, tt.name, err)
			}
		})
	}
}

func TestRetentionBlocksDoesNotOverflow(t *testing.T) {
	conf := DefaultConfig()
	conf.RetentionDays = maxRetentionDays
	require.NoError(t, conf.BasicCheck())
	assert.Equal(t, uint32(maxRetentionDays*blockPerDay), conf.RetentionBlocks())
	assert.Greater(t, conf.RetentionBlocks(), uint32(maxRetentionDays-1)*blockPerDay)
}

func TestConfigStorePath(t *testing.T) {
	conf := DefaultConfig()
	conf.Path = util.TempDirPath()
	require.NoError(t, conf.BasicCheck())

	if runtime.GOOS != "windows" {
		assert.Equal(t, conf.Path+"/store.db", conf.StorePath())
	} else {
		assert.Equal(t, conf.Path+"\\store.db", conf.StorePath())
	}
}
