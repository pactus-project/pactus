package addr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseP2PAddress(t *testing.T) {
	testCases := []struct {
		input    string
		expected P2PAddr
		err      bool
	}{
		{
			input: "/dns/bootstrap2.pactus.org/tcp/21888/p2p/12D3KooWM39ag7ghta49qybf7McADgT8FLakTYkCsiBvwdnjuG5q",
			expected: P2PAddr{
				transport: "dns",
				protocol:  "tcp",
				addr:      "bootstrap2.pactus.org",
				port:      "21888",
				peerID:    "12D3KooWM39ag7ghta49qybf7McADgT8FLakTYkCsiBvwdnjuG5q",
			},
			err: false,
		},
		{
			input: "/ip4/115.193.157.138/tcp/21888",
			expected: P2PAddr{
				transport: "ip4",
				protocol:  "tcp",
				addr:      "115.193.157.138",
				port:      "21888",
				peerID:    "",
			},
			err: false,
		},
		{
			input: "/ip4/84.247.165.249/tcp/21888/p2p/12D3KooWQmv2FcNQfh1EhA98twt8ePdkQaxEPeYfinEYyVS16juY",
			expected: P2PAddr{
				transport: "ip4",
				protocol:  "tcp",
				addr:      "84.247.165.249",
				port:      "21888",
				peerID:    "12D3KooWQmv2FcNQfh1EhA98twt8ePdkQaxEPeYfinEYyVS16juY",
			},
			err: false,
		},
		{
			input: "/ip6/2001:db8:85a3:0000:0000:8a2e:0370:7334/tcp/21888/" +
				"p2p/12D3KooWQmv2FcNQfh1EhA98twt8ePdkQaxEPeYfinEYyVS16juY",
			expected: P2PAddr{
				transport: "ip6",
				protocol:  "tcp",
				addr:      "2001:db8:85a3:0000:0000:8a2e:0370:7334",
				port:      "21888",
				peerID:    "12D3KooWQmv2FcNQfh1EhA98twt8ePdkQaxEPeYfinEYyVS16juY",
			},
			err: false,
		},
		{
			input: "/ip4/159.148.146.149/udp/21888/p2p/12D3KooWKCokWtpdudxgsLRoFcnrW35vhn6w632iGWCga7E5e68Q",
			expected: P2PAddr{
				transport: "ip4",
				protocol:  "udp",
				addr:      "159.148.146.149",
				port:      "21888",
				peerID:    "12D3KooWKCokWtpdudxgsLRoFcnrW35vhn6w632iGWCga7E5e68Q",
			},
			err: false,
		},
		{
			input: "/ip6/2a01:4f9:4a:1d85::2/tcp/21888",
			expected: P2PAddr{
				transport: "ip6",
				protocol:  "tcp",
				addr:      "2a01:4f9:4a:1d85::2",
				port:      "21888",
				peerID:    "",
			},
			err: false,
		},
		{
			input: "/ip6/2a02:ab88:7601:f700:7302:a9be:3ac6:75d5/tcp/21888",
			expected: P2PAddr{
				transport: "ip6",
				protocol:  "tcp",
				addr:      "2a02:ab88:7601:f700:7302:a9be:3ac6:75d5",
				port:      "21888",
				peerID:    "",
			},
			err: false,
		},
		// Additional test cases
		{
			input: "/ip4/127.0.0.1/tcp/8080",
			expected: P2PAddr{
				transport: "ip4",
				protocol:  "tcp",
				addr:      "127.0.0.1",
				port:      "8080",
				peerID:    "",
			},
			err: false,
		},
		{
			input: "/ip6/::1/tcp/8080",
			expected: P2PAddr{
				transport: "ip6",
				protocol:  "tcp",
				addr:      "::1",
				port:      "8080",
				peerID:    "",
			},
			err: false,
		},
		{
			input: "/dns/example.com/udp/5353",
			expected: P2PAddr{
				transport: "dns",
				protocol:  "udp",
				addr:      "example.com",
				port:      "5353",
				peerID:    "",
			},
			err: false,
		},
		{
			input: "/ip4/192.168.1.1/tcp/9999/p2p/12D3KooWKCokWtpdudxgsLRoFcnrW35vhn6w632iGWCga7E5e68Q",
			expected: P2PAddr{
				transport: "ip4",
				protocol:  "tcp",
				addr:      "192.168.1.1",
				port:      "9999",
				peerID:    "12D3KooWKCokWtpdudxgsLRoFcnrW35vhn6w632iGWCga7E5e68Q",
			},
			err: false,
		},
		{
			input: "/ip6/2001:db8::1/tcp/12345/p2p/12D3KooWKCokWtpdudxgsLRoFcnrW35vhn6w632iGWCga7E5e68Q",
			expected: P2PAddr{
				transport: "ip6",
				protocol:  "tcp",
				addr:      "2001:db8::1",
				port:      "12345",
				peerID:    "12D3KooWKCokWtpdudxgsLRoFcnrW35vhn6w632iGWCga7E5e68Q",
			},
			err: false,
		},
		// Invalid cases
		{
			input:    "/ip4/192.168.1.1/tcp", // Missing port
			expected: P2PAddr{},
			err:      true,
		},
		{
			input:    "/ip4//tcp/12345", // Missing address
			expected: P2PAddr{},
			err:      true,
		},
		{
			input:    "/ip4/192.168.1.1/12345", // Missing protocol
			expected: P2PAddr{},
			err:      true,
		},
		{
			input:    "/unknown/192.168.1.1/tcp/12345", // Unknown transport
			expected: P2PAddr{},
			err:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Parse(tc.input)
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
