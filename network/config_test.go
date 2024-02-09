package network

import (
	"testing"

	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"
)

func TestConfigBasicCheck(t *testing.T) {
	testCases := []struct {
		name        string
		expectError error
		updateFn    func(c *Config)
	}{
		{
			name: "Empty ListenAddrStrings - Expect Error",
			expectError: ConfigError{
				Reason: "address is not valid: failed to parse multiaddr \"\": empty multiaddr",
			},
			updateFn: func(c *Config) {
				c.ListenAddrStrings = []string{""}
			},
		},
		{
			name: "Both Relay and Relay Service be true - Expect Error",
			expectError: ConfigError{
				Reason: "both the relay and relay service cannot be active at the same time",
			},
			updateFn: func(c *Config) {
				c.EnableRelay = true
				c.EnableRelayService = true
			},
		},
		{
			name: "Invalid ListenAddrStrings - Expect Error",
			expectError: ConfigError{
				Reason: "address is not valid: failed to parse multiaddr \"127.0.0.1\": must begin with /",
			},
			updateFn: func(c *Config) {
				c.ListenAddrStrings = []string{"127.0.0.1"}
			},
		},
		{
			name: "Invalid ListenAddrStrings (No port) - Expect Error",
			expectError: ConfigError{
				Reason: "address is not valid: failed to parse multiaddr \"/ip4\": unexpected end of multiaddr",
			},
			updateFn: func(c *Config) {
				c.ListenAddrStrings = []string{"/ip4"}
			},
		},
		{
			name: "Invalid Public Address - Expect Error",
			expectError: ConfigError{
				Reason: "address is not valid: failed to parse multiaddr \"/ip4\": unexpected end of multiaddr",
			},
			updateFn: func(c *Config) {
				c.PublicAddrString = "/ip4"
			},
		},
		{
			name: "Invalid DefaultBootstrapAddrStrings - Expect Error",
			expectError: ConfigError{
				Reason: "address is not valid: invalid p2p multiaddr",
			},
			updateFn: func(c *Config) {
				c.DefaultBootstrapAddrStrings = []string{"/ip4/127.0.0.1/"}
			},
		},
		{
			name: "Invalid BootstrapAddrStrings - Expect Error",
			expectError: ConfigError{
				Reason: "address is not valid: invalid p2p multiaddr",
			},
			updateFn: func(c *Config) {
				c.BootstrapAddrStrings = []string{"/ip4/127.0.0.1/"}
			},
		},
		{
			name:        "Valid Public Address - No Error",
			expectError: nil,
			updateFn: func(c *Config) {
				c.PublicAddrString = "/ip4/127.0.0.1/"
			},
		},
		{
			name:        "Valid ListenAddrStrings - No Error",
			expectError: nil,
			updateFn: func(c *Config) {
				c.ListenAddrStrings = []string{"/ip4/127.0.0.1"}
			},
		},
		{
			name:        "Valid BootstrapAddrStrings - No Error",
			expectError: nil,
			updateFn: func(c *Config) {
				c.BootstrapAddrStrings = []string{"/ip4/127.0.0.1/p2p/12D3KooWQBpPV6NtZy1dvN2oF7dJdLoooRZfEmwtHiDUf42ArDjT"}
			},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := DefaultConfig()
			tc.updateFn(conf)
			if tc.expectError != nil {
				err := conf.BasicCheck()
				assert.ErrorIs(t, tc.expectError, err,
					"Expected error not matched for test %d-%s: expected %s, got: %s", i, tc.name, tc.expectError, err)
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
		config      Config
		expectedMax int
		expectedMin int
	}{
		{Config{MaxConns: 1}, 1, 0},
		{Config{MaxConns: 8}, 8, 2},
		{Config{MaxConns: 30}, 32, 8},
		{Config{MaxConns: 1000}, 1024, 256},
	}

	for _, test := range tests {
		resultMax := test.config.ScaledMaxConns()
		resultMin := test.config.ScaledMinConns()
		if resultMax != test.expectedMax ||
			resultMin != test.expectedMin {
			t.Errorf("For MaxConns %d, "+
				"NormedMaxConns() returned %d (expected %d), "+
				"NormedMinConns() returned %d (expected %d)",
				test.config.MaxConns,
				resultMax, test.expectedMax,
				resultMin, test.expectedMin)
		}
	}
}

func TestPublicAddr(t *testing.T) {
	conf1 := DefaultConfig()
	assert.Nil(t, conf1.PublicAddr())

	conf2 := DefaultConfig()
	conf2.PublicAddrString = "/ip4/127.0.0.1/p2p/12D3KooWQBpPV6NtZy1dvN2oF7dJdLoooRZfEmwtHiDUf42ArDjT"
	assert.NotNil(t, conf2.PublicAddr())
}
