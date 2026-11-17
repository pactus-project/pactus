package firewall

import (
	"bytes"
	"testing"
)

// The firewall is the first code that reads the bytes of a peer, from the gossip
// and from a stream: opening a bundle from any bytes never panics.
func FuzzOpenBundle(f *testing.F) {
	// The seeds are in testdata/fuzz/FuzzOpenBundle: a vote, a proposal, a block announcement
	// and a transactions message.
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		td := setup(t, nil)

		_, _ = td.firewall.OpenGossipBundle(data, td.goodPeerID)
		_, _ = td.firewall.OpenStreamBundle(bytes.NewReader(data), td.goodPeerID)
		_, _ = td.firewall.OpenGossipBundle(data, td.unknownPeerID)
	})
}
