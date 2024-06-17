package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetIPFromMultiAddress(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		expectedIP  string
		expectError bool
	}{
		{
			name:       "Valid IPv4 with p2p",
			address:    "/ip4/84.247.165.249/tcp/21888/p2p/12D3KooWQmv2FcNQfh1EhA98twt8ePdkQaxEPeYfinEYyVS16juY",
			expectedIP: "84.247.165.249",
		},
		{
			name:       "Valid IPv4 without p2p",
			address:    "/ip4/115.193.157.138/tcp/21888",
			expectedIP: "115.193.157.138",
		},
		{
			name: "Valid IPv6 with p2p",
			address: "/ip6/240e:390:8a1:ae80:7dbc:64b6:e84c:d2bf/tcp/21888/p2p/" +
				"12D3KooWQmv2FcNQfh1EhA98twt8ePdkQaxEPeYfinEYyVS16juY",
			expectedIP: "240e:390:8a1:ae80:7dbc:64b6:e84c:d2bf",
		},
		{
			name:        "Invalid address",
			address:     "/invalid/address",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, err := GetIPFromMultiAddress(tt.address)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIP, ip)
			}
		})
	}
}
