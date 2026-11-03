package sync

import (
	"testing"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/sync/peerset/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PeerSet.RemovePeer exists but nothing calls it: every peer ID ever seen stays
// in memory forever (with its metrics maps, moniker, agent and keys). Peer IDs
// are free to generate, so an attacker can make the node keep one record per
// connection attempt.
func TestPeerSetKeepsDisconnectedStrangers(t *testing.T) {
	td := setup(t, nil)

	const numPeers = 500
	for i := 0; i < numPeers; i++ {
		pid := td.RandPeerID()
		td.sync.processConnectEvent(&network.ConnectEvent{
			PeerID:        pid,
			RemoteAddress: "/ip4/2.2.2.2/tcp/21888",
			Direction:     lp2pnetwork.DirInbound,
		})
		td.sync.processDisconnectEvent(&network.DisconnectEvent{PeerID: pid})
	}

	assert.Less(t, td.sync.peerSet.Len(), numPeers/10,
		"records of peers that disconnected without completing the handshake should be reclaimed")
}

// The documentation of GetPeer says it returns a copy, but it returns the live
// pointer stored in the map; callers (firewall, handlers) read its fields without
// holding the lock while other goroutines update it.
func TestPeerSetGetPeerReturnsLivePointer(t *testing.T) {
	td := setup(t, nil)

	pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	snapshot := td.sync.peerSet.GetPeer(pid)

	td.sync.peerSet.UpdateStatus(pid, status.StatusBanned)

	assert.Equal(t, status.StatusKnown, snapshot.Status,
		"GetPeer is documented as returning a copy; later updates must not alter it")
}

// Same issue as GetPeer, for GetRandomPeer: it returned the live
// pointer stored in the map instead of a copy.
func TestGetRandomPeerReturnsLivePointer(t *testing.T) {
	td := setup(t, nil)

	pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	snapshot := td.sync.peerSet.GetRandomPeer()
	assert.Equal(t, pid, snapshot.PeerID)

	td.sync.peerSet.UpdateStatus(pid, status.StatusBanned)

	assert.Equal(t, status.StatusKnown, snapshot.Status,
		"GetRandomPeer must return a copy; later updates must not alter it")
}

// Same issue as GetPeer, for Sessions: it returned the live pointers
// stored in the session manager's map instead of copies.
func TestSessionsReturnsLivePointers(t *testing.T) {
	td := setup(t, nil)

	pid := td.addPeer(t, status.StatusKnown, service.New(service.FullNode))
	sid := td.sync.peerSet.OpenSession(pid, 1, 10)

	sessions := td.sync.peerSet.Sessions()
	require.Len(t, sessions, 1)
	snapshot := sessions[0]

	td.sync.peerSet.SetSessionCompleted(sid)

	assert.Equal(t, session.Open, snapshot.Status,
		"Sessions must return copies; later updates must not alter them")
}
