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
			name:        "Invalid RelayAddrStrings - Expect Error",
			expectError: true,
			updateFn: func(c *Config) {
				c.RelayAddrStrings = []string{"/ip4/127.0.0.1/"}
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
			name:        "Valid RelayAddrStrings - No Error",
			expectError: false,
			updateFn: func(c *Config) {
				c.RelayAddrStrings = []string{"/ip4/127.0.0.1/p2p/12D3KooWQBpPV6NtZy1dvN2oF7dJdLoooRZfEmwtHiDUf42ArDjT"}
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

	assert.False(t, conf.IsBootstrapper(pid1))
	assert.True(t, conf.IsBootstrapper(pid2))
	assert.True(t, conf.IsBootstrapper(pid3))
}
