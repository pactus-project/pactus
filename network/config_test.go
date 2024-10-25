package network

import (
	"testing"

	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"
)

func TestConfigBasicCheck(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
		updateFn    func(c *Config)
	}{
		{
			name: "Empty ListenAddrStrings",
			expectedErr: ConfigError{
				Reason: "address is not valid: failed to parse multiaddr \"\": empty multiaddr",
			},
			updateFn: func(c *Config) {
				c.ListenAddrStrings = []string{""}
			},
		},
		{
			name: "Both Relay and Relay Service be true",
			expectedErr: ConfigError{
				Reason: "both the relay and relay service cannot be active at the same time",
			},
			updateFn: func(c *Config) {
				c.EnableRelay = true
				c.EnableRelayService = true
			},
		},
		{
			name: "Invalid ListenAddrStrings",
			expectedErr: ConfigError{
				Reason: "address is not valid: failed to parse multiaddr \"127.0.0.1\": must begin with /",
			},
			updateFn: func(c *Config) {
				c.ListenAddrStrings = []string{"127.0.0.1"}
			},
		},
		{
			name: "Invalid ListenAddrStrings (No port)",
			expectedErr: ConfigError{
				Reason: "address is not valid: failed to parse multiaddr \"/ip4\": unexpected end of multiaddr",
			},
			updateFn: func(c *Config) {
				c.ListenAddrStrings = []string{"/ip4"}
			},
		},
		{
			name: "Invalid Public Address",
			expectedErr: ConfigError{
				Reason: "address is not valid: failed to parse multiaddr \"/ip4\": unexpected end of multiaddr",
			},
			updateFn: func(c *Config) {
				c.PublicAddrString = "/ip4"
			},
		},
		{
			name: "Invalid DefaultBootstrapAddrStrings",
			expectedErr: ConfigError{
				Reason: "address is not valid: invalid p2p multiaddr",
			},
			updateFn: func(c *Config) {
				c.DefaultBootstrapAddrStrings = []string{"/ip4/127.0.0.1/"}
			},
		},
		{
			name: "Invalid BootstrapAddrStrings",
			expectedErr: ConfigError{
				Reason: "address is not valid: invalid p2p multiaddr",
			},
			updateFn: func(c *Config) {
				c.BootstrapAddrStrings = []string{"/ip4/127.0.0.1/"}
			},
		},
		{
			name: "Low MaxConns",
			expectedErr: ConfigError{
				Reason: "maximum connection should be greater than 16",
			},
			updateFn: func(c *Config) {
				c.MaxConns = 8
			},
		},
		{
			name: "Valid Public Address",
			updateFn: func(c *Config) {
				c.PublicAddrString = "/ip4/127.0.0.1/"
			},
		},
		{
			name: "Valid ListenAddrStrings",
			updateFn: func(c *Config) {
				c.ListenAddrStrings = []string{"/ip4/127.0.0.1"}
			},
		},
		{
			name: "Valid BootstrapAddrStrings",
			updateFn: func(c *Config) {
				c.BootstrapAddrStrings = []string{
					"/ip4/127.0.0.1/p2p/12D3KooWQBpPV6NtZy1dvN2oF7dJdLoooRZfEmwtHiDUf42ArDjT",
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
				assert.ErrorIs(t, err, tt.expectedErr,
					"Expected error not matched for test %d-%s, expected: %s, got: %s",
					no, tt.name, tt.expectedErr, err)
			} else {
				err := conf.BasicCheck()
				assert.NoError(t, err,
					"Expected no error for test %d-%s, get: %s",
					no, tt.name, err)
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

func TestMinConns(t *testing.T) {
	tests := []struct {
		config      Config
		expectedMin int
	}{
		{Config{MaxConns: 16}, 2},
		{Config{MaxConns: 30}, 5},
		{Config{MaxConns: 128}, 30},
	}

	for _, tt := range tests {
		resultMin := tt.config.MinConns()
		if resultMin != tt.expectedMin {
			t.Errorf("For MaxConns %d, "+
				"MinConns() returned %d (expected %d)",
				tt.config.MaxConns,
				resultMin, tt.expectedMin)
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
