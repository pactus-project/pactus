package network

import (
	"testing"

	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"
)

func TestConfigBasicCheck(t *testing.T) {
	testCases := []struct {
		name        string
		expectError bool
		updateFn    func(c *Config)
	}{
		{
			name:        "Empty ListenAddrStrings - Expect Error",
			expectError: true,
			updateFn: func(c *Config) {
				c.ListenAddrStrings = []string{""}
			},
		},
		{
			name:        "Both Relay and Relay Service be true - Expect Error",
			expectError: true,
			updateFn: func(c *Config) {
				c.EnableRelay = true
				c.EnableRelayService = true
			},
		},
		{
			name:        "Invalid ListenAddrStrings - Expect Error",
			expectError: true,
			updateFn: func(c *Config) {
				c.ListenAddrStrings = []string{"127.0.0.1"}
			},
		},
		{
			name:        "Invalid ListenAddrStrings (No port) - Expect Error",
			expectError: true,
			updateFn: func(c *Config) {
				c.ListenAddrStrings = []string{"/ip4"}
			},
		},
		{
			name:        "Invalid Public Address - Expect Error",
			expectError: true,
			updateFn: func(c *Config) {
				c.PublicAddrString = "/ip4"
			},
		},
		{
			name:        "Invalid DefaultRelayAddrStrings - Expect Error",
			expectError: true,
			updateFn: func(c *Config) {
				c.DefaultRelayAddrStrings = []string{"/ip4/127.0.0.1/"}
			},
		},
		{
			name:        "Invalid DefaultBootstrapAddrStrings - Expect Error",
			expectError: true,
			updateFn: func(c *Config) {
				c.DefaultBootstrapAddrStrings = []string{"/ip4/127.0.0.1/"}
			},
		},
		{
			name:        "Invalid BootstrapAddrStrings - Expect Error",
			expectError: true,
			updateFn: func(c *Config) {
				c.BootstrapAddrStrings = []string{"/ip4/127.0.0.1/"}
			},
		},
		{
			name:        "Valid Public Address - No Error",
			expectError: false,
			updateFn: func(c *Config) {
				c.PublicAddrString = "/ip4/127.0.0.1/"
			},
		},
		{
			name:        "Valid ListenAddrStrings - No Error",
			expectError: false,
			updateFn: func(c *Config) {
				c.ListenAddrStrings = []string{"/ip4/127.0.0.1"}
			},
		},
		{
			name:        "Valid BootstrapAddrStrings - No Error",
			expectError: false,
			updateFn: func(c *Config) {
				c.BootstrapAddrStrings = []string{"/ip4/127.0.0.1/p2p/12D3KooWQBpPV6NtZy1dvN2oF7dJdLoooRZfEmwtHiDUf42ArDjT"}
			},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := DefaultConfig()
			tc.updateFn(conf)
			if tc.expectError {
				assert.Error(t, conf.BasicCheck(), "Expected error for Test %d: %s", i, tc.name)
			} else {
				assert.NoError(t, conf.BasicCheck(), "Expected no error for Test %d: %s", i, tc.name)
			}
		})
	}
}

func TestIsBootstrapper(t *testing.T) {
	conf := DefaultConfig()
	conf.BootstrapAddrStrings = []string{"/ip4/127.0.0.1/p2p/12D3KooWQBpPV6NtZy1dvN2oF7dJdLoooRZfEmwtHiDUf42ArDjT"}
	conf.DefaultBootstrapAddrStrings = []string{"/ip4/127.0.0.2/p2p/12D3KooWBqutgDboACf1i1c9uN9BQg9xdREoeXYb2rvFHQU1QcAp"}

	pid1, _ := lp2ppeer.Decode("12D3KooWQQKidG8Nn6fLgxjryHFhRCfG9fUWU88yGSZNd59Kbqka")
	pid2, _ := lp2ppeer.Decode("12D3KooWQBpPV6NtZy1dvN2oF7dJdLoooRZfEmwtHiDUf42ArDjT")
	pid3, _ := lp2ppeer.Decode("12D3KooWBqutgDboACf1i1c9uN9BQg9xdREoeXYb2rvFHQU1QcAp")

	conf.CheckIsBootstrapper(pid1)
	assert.False(t, conf.IsBootstrapper)

	conf.CheckIsBootstrapper(pid2)
	assert.True(t, conf.IsBootstrapper)

	conf.CheckIsBootstrapper(pid3)
	assert.True(t, conf.IsBootstrapper)
}

func TestScaledConns(t *testing.T) {
	tests := []struct {
		config            Config
		expectedMax       int
		expectedMin       int
		expectedThreshold int
	}{
		{Config{MaxConns: 1}, 1, 0, 0},
		{Config{MaxConns: 8}, 8, 2, 1},
		{Config{MaxConns: 30}, 32, 8, 4},
		{Config{MaxConns: 1000}, 1024, 256, 128},
	}

	for _, test := range tests {
		resultMax := test.config.ScaledMaxConns()
		resultMin := test.config.ScaledMinConns()
		resultThreshold := test.config.ConnsThreshold()
		if resultMax != test.expectedMax ||
			resultMin != test.expectedMin ||
			resultThreshold != test.expectedThreshold {
			t.Errorf("For MaxConns %d, "+
				"NormedMaxConns() returned %d (expected %d), "+
				"NormedMinConns() returned %d (expected %d), "+
				"ConnsThreshold() returned %d (expected %d)",
				test.config.MaxConns,
				resultMax, test.expectedMax,
				resultMin, test.expectedMin,
				resultThreshold, test.expectedThreshold)
		}
	}
}
