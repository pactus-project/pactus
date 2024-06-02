package ipblocker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		blackListCidr  []string
		expectedError  bool
		expectedLength int
	}{
		{
			"Valid CIDRs",
			[]string{"240e:390:8a1:ae80::/64", "192.168.1.0/24"},
			false, 2,
		},
		{"Invalid CIDR", []string{"invalid-cidr"}, true, 0},
		{"Empty CIDRs", []string{}, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipBlocker, err := New(tt.blackListCidr)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedLength, len(ipBlocker.cidrs))
			}
		})
	}
}

func TestIsBlocked(t *testing.T) {
	tests := []struct {
		name          string
		blackListCidr []string
		ip            string
		expected      bool
	}{
		{
			"Blocked IPv6",
			[]string{"240e:390:8a1:ae80::/64"},
			"240e:390:8a1:ae80:7dbc:64b6:e84c:d2bf", true,
		},
		{
			"Not Blocked IPv6",
			[]string{"240e:390:8a1:ae80::/64"},
			"240e:391:8a1:ae80:7dbc:64b6:e84c:d2bf", false,
		},
		{"Blocked IPv4", []string{"192.168.1.0/24"}, "192.168.1.1", true},
		{"Not Blocked IPv4", []string{"192.168.1.0/24"}, "10.0.0.1", false},
		{"Empty CIDR List", []string{}, "192.168.1.1", false},
		{"Invalid IP", []string{"192.168.1.0/24"}, "invalid-ip", false},
		{
			"Blocked IPv4 in multiple CIDRs",
			[]string{"10.0.0.0/8", "192.168.1.0/24"},
			"192.168.1.1", true,
		},
		{
			"Blocked IPv6 in multiple CIDRs",
			[]string{"240e:390:8a1:ae80::/64", "2001:db8::/32"},
			"2001:db8::1", true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipBlocker, err := New(tt.blackListCidr)
			assert.NoError(t, err)
			result := ipBlocker.IsBlocked(tt.ip)
			assert.Equal(t, tt.expected, result)
		})
	}
}
