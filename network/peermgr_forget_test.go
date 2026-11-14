package network

import (
	"testing"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// peerMgr.peers is only ever added to: every inbound stranger (random peer ID,
// ephemeral source port as "address") stays forever, is written to the peer
// store file, and is re-dialled by CheckConnectivity when the node has fewer
// than MinConns connections.
func TestPeerMgrForgetsDisconnectedInboundPeers(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	mgr := &peerMgr{peers: map[lp2ppeer.ID]*peerInfo{}}

	ma, err := multiaddr.NewMultiaddr("/ip4/2.2.2.2/tcp/54321")
	require.NoError(t, err)

	const numPeers = 500
	for i := 0; i < numPeers; i++ {
		pid := ts.RandPeerID()
		mgr.SetPeerConnected(pid, ma, lp2pnetwork.DirInbound)
		mgr.SetPeerDisconnected(pid)
	}

	assert.Less(t, len(mgr.peers), numPeers/10,
		"disconnected inbound peers (not bootstrap/known addresses) should not be kept and re-dialled")
}
