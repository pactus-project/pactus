package network

import (
	"testing"

	lp2pspb "github.com/libp2p/go-libp2p-pubsub/pb"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
)

func TestMakeMultiAddrs(t *testing.T) {
	tests := []struct {
		name       string
		inputAddrs []string
		expected   []multiaddr.Multiaddr
	}{
		{
			inputAddrs: []string{
				"/ip4/127.0.0.1/tcp/1234",
				"/ip6/::1/tcp/5678/",
				"/dns4/example.com",
			},
			expected: []multiaddr.Multiaddr{
				multiaddr.Cast([]byte{0x04, 0x7f, 0x00, 0x00, 0x01, 0x06, 0x04, 0xd2}),
				multiaddr.Cast([]byte{
					0x29, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x06, 0x16, 0x2e,
				}),
				multiaddr.Cast([]byte{0x36, 0x0b, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d}),
			},
		},
		{
			inputAddrs: []string{
				"invalid_address",
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualPis, actualError := MakeMultiAddrs(tt.inputAddrs)

			if tt.expected != nil {
				assert.Equal(t, tt.expected, actualPis)
				assert.NoError(t, actualError)
			} else {
				assert.Error(t, actualError)
				assert.Nil(t, actualPis)
			}
		})
	}
}

func TestMakeAddrInfos(t *testing.T) {
	pid, _ := lp2ppeer.Decode("12D3KooWCwQZt8UriVXobQHPXPR8m83eceXVoeT6brPNiBHomebc")
	tests := []struct {
		name        string
		inputAddrs  []string
		expectedPis []lp2ppeer.AddrInfo
	}{
		{
			inputAddrs: []string{
				"/ip4/127.0.0.1/tcp/1234/p2p/12D3KooWCwQZt8UriVXobQHPXPR8m83eceXVoeT6brPNiBHomebc",
				"/ip6/::1/tcp/5678/p2p/12D3KooWCwQZt8UriVXobQHPXPR8m83eceXVoeT6brPNiBHomebc",
				"/dns4/example.com/tcp/4001/p2p/12D3KooWCwQZt8UriVXobQHPXPR8m83eceXVoeT6brPNiBHomebc",
			},
			expectedPis: []lp2ppeer.AddrInfo{
				{
					ID: pid,
					Addrs: []multiaddr.Multiaddr{
						multiaddr.StringCast("/ip4/127.0.0.1/tcp/1234"),
					},
				},
				{
					ID: pid,
					Addrs: []multiaddr.Multiaddr{
						multiaddr.StringCast("/ip6/::1/tcp/5678"),
					},
				},
				{
					ID: pid,
					Addrs: []multiaddr.Multiaddr{
						multiaddr.StringCast("/dns4/example.com/tcp/4001"),
					},
				},
			},
		},
		{
			inputAddrs: []string{
				"/ip4/127.0.0.1/tcp/1234", // No peer id
			},
			expectedPis: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualPis, actualError := MakeAddrInfos(tt.inputAddrs)

			if tt.expectedPis != nil {
				assert.Equal(t, tt.expectedPis, actualPis)
				assert.NoError(t, actualError)
			} else {
				assert.Error(t, actualError)
				assert.Nil(t, actualPis)
			}
		})
	}
}

func TestIPToMultiAddr(t *testing.T) {
	tests := []struct {
		ip       string
		port     int
		expected string
	}{
		{"127.0.0.1", 8080, "/ip4/127.0.0.1/tcp/8080"},
		{"192.168.1.1", 1234, "/ip4/192.168.1.1/tcp/1234"},
		{"::1", 80, "/ip6/::1/tcp/80"},
		{"invalid_ip", 80, ""},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			ma, err := IPToMultiAddr(tt.ip, tt.port)
			if tt.expected != "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, ma.String())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestHasPID(t *testing.T) {
	pids := []lp2ppeer.ID{"peer1", "peer2", "peer3"}

	assert.True(t, HasPID(pids, lp2ppeer.ID("peer1")))
	assert.False(t, HasPID(pids, lp2ppeer.ID("peer4")))
}

func TestSubnetsToFilters(t *testing.T) {
	sns := PrivateSubnets()
	filter := SubnetsToFilters(sns, multiaddr.ActionDeny)

	ma1, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0")
	ma2, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1")
	ma3, _ := multiaddr.NewMultiaddr("/ip4/8.8.8.8")

	assert.False(t, filter.AddrBlocked(ma1))
	assert.True(t, filter.AddrBlocked(ma2))
	assert.False(t, filter.AddrBlocked(ma3))
}

func TestMessageIdFunc(t *testing.T) {
	m := &lp2pspb.Message{Data: []byte("pactus")}
	id := MessageIDFunc(m)

	assert.Equal(t, id, string([]byte{
		0xea, 0x02, 0x0a, 0xce, 0x5c, 0x96, 0x8f, 0x75,
		0x5d, 0xfc, 0x1b, 0x59, 0x21, 0xe5, 0x74, 0x19, 0x1c, 0xd9, 0xff, 0x43,
	}))
}
