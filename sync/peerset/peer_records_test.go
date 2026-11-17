package peerset

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestForgetDisconnectedPeersOrder(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	peerSet := NewPeerSet(time.Minute)

	start := time.Now()
	disconnected := make([]peer.ID, 0, maxPeerRecords)
	for i := range maxPeerRecords {
		pid := ts.RandPeerID()
		peerSet.UpdateStatus(pid, status.StatusConnected)
		peerSet.UpdateStatus(pid, status.StatusDisconnected)
		peerSet.peers[pid].LastReceived = start.Add(time.Duration(i) * time.Second)
		disconnected = append(disconnected, pid)
	}

	// The oldest disconnected peer has an open session: it is kept.
	withSession := disconnected[0]
	peerSet.OpenSession(withSession, 1, 10)

	// The bound is reached: the records of the disconnected peers that were heard from
	// least recently are forgotten, until a quarter of the records is free.
	peerSet.UpdateStatus(ts.RandPeerID(), status.StatusConnected)

	assert.Equal(t, maxPeerRecords*3/4+1, peerSet.Len())
	assert.True(t, peerSet.HasPeer(withSession), "a peer with an open session is kept")
	assert.False(t, peerSet.HasPeer(disconnected[1]), "the oldest peers are forgotten first")
	assert.True(t, peerSet.HasPeer(disconnected[len(disconnected)-1]), "the newest peers are kept")
}
