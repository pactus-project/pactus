package network

import (
	"testing"

	lp2pspb "github.com/libp2p/go-libp2p-pubsub/pb"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
)

func TestMakeMultiAddrs(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualPis, actualError := MakeMultiAddrs(tc.inputAddrs)

			if tc.expected != nil {
				assert.Equal(t, actualPis, tc.expected)
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
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualPis, actualError := MakeAddrInfos(tc.inputAddrs)

			if tc.expectedPis != nil {
				assert.Equal(t, actualPis, tc.expectedPis)
				assert.NoError(t, actualError)
			} else {
				assert.Error(t, actualError)
				assert.Nil(t, actualPis)
			}
		})
	}
}

func TestIPToMultiAddr(t *testing.T) {
	testCases := []struct {
		ip       string
		port     int
		expected string
	}{
		{"127.0.0.1", 8080, "/ip4/127.0.0.1/tcp/8080"},
		{"192.168.1.1", 1234, "/ip4/192.168.1.1/tcp/1234"},
		{"::1", 80, "/ip6/::1/tcp/80"},
		{"invalid_ip", 80, ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expected, func(t *testing.T) {
			ma, err := IPToMultiAddr(testCase.ip, testCase.port)
			if testCase.expected != "" {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, ma.String())
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
	f := SubnetsToFilters(sns, multiaddr.ActionDeny)

	ma1, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0")
	ma2, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1")
	ma3, _ := multiaddr.NewMultiaddr("/ip4/8.8.8.8")

	assert.False(t, f.AddrBlocked(ma1))
	assert.True(t, f.AddrBlocked(ma2))
	assert.False(t, f.AddrBlocked(ma3))
}

func TestMessageIdFunc(t *testing.T) {
	m := &lp2pspb.Message{Data: []byte("zarb")}
	id := MessageIDFunc(m)

	assert.Equal(t, "\x12\xb3\x89\x77\xf2\xd6\x7f\x06\xf0\xc0\xcd\x54\xaa\xf7\x32\x4c\xf4\xfe\xe1\x84", id)
}
