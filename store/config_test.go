package store

import (
	"runtime"
	"testing"

	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestConfigBasicCheck(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr error
		updateFn    func(c *Config)
	}{
		{
			name: "Invalid Path",
			expectedErr: ConfigError{
				Reason: "path is not valid",
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
				Reason: "cache size set to zero",
			},
			updateFn: func(c *Config) {
				c.AccountCacheSize = 0
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

func TestConfigStorePath(t *testing.T) {
	conf := DefaultConfig()
	conf.Path = util.TempDirPath()
	assert.NoError(t, conf.BasicCheck())

	if runtime.GOOS != "windows" {
		assert.Equal(t, conf.StorePath(), conf.Path+"/store.db")
	} else {
		assert.Equal(t, conf.StorePath(), conf.Path+"\\store.db")
	}
}
