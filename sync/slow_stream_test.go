package sync

import (
	"io"
	"testing"
	"time"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/network"
	"github.com/stretchr/testify/assert"
)

// A peer that opens a stream and never writes anything must not freeze the
// processing of every other network event.
//
// processStreamMessage reads the stream (bdl.Decode) inside the single receiver
// goroutine of the network pipeline. A silent stream therefore blocks the whole
// event loop (gossip, connect/disconnect events, hello...) until the stream
// deadline (network.Config.StreamTimeout, 20s by default) expires.
func TestSlowStreamBlocksEventLoop(t *testing.T) {
	td := setup(t, nil)

	silentReader, silentWriter := io.Pipe()
	// Releases the blocked goroutine at the end of the test.
	t.Cleanup(func() { _ = silentWriter.Close() })

	// A hostile peer opens a stream but never sends any byte.
	td.network.EventPipe.Send(&network.StreamMessage{
		From:   td.RandPeerID(),
		Reader: silentReader,
	})

	// A legitimate connect event arrives just after it.
	honestPeer := td.RandPeerID()
	td.network.EventPipe.Send(&network.ConnectEvent{
		PeerID:        honestPeer,
		RemoteAddress: "/ip4/2.2.2.2/tcp/21888",
		Direction:     lp2pnetwork.DirInbound,
	})

	assert.Eventually(t, func() bool {
		return td.sync.peerSet.HasPeer(honestPeer)
	}, 2*time.Second, 20*time.Millisecond,
		"the connect event of an honest peer must be processed while a silent stream is pending")
}
