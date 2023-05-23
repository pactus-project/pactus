package peerset

import (
	"fmt"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/stretchr/testify/assert"
)

func TestPeerSet(t *testing.T) {
	peerSet := NewPeerSet(time.Second)

	// Add peers using UpdatePeerInfo
	pk1, _ := bls.GenerateTestKeyPair()
	pk2, _ := bls.GenerateTestKeyPair()
	pk3, _ := bls.GenerateTestKeyPair()
	peerSet.UpdatePeerInfo(peer.ID("peer1"), StatusCodeBanned, "Moniker1", "Agent1", pk1, true)
	peerSet.UpdatePeerInfo(peer.ID("peer2"), StatusCodeKnown, "Moniker2", "Agent2", pk2, false)
	peerSet.UpdatePeerInfo(peer.ID("peer3"), StatusCodeTrusty, "Moniker3", "Agent3", pk3, true)

	t.Run("Testing Len", func(t *testing.T) {
		assert.Equal(t, 3, peerSet.Len())
	})

	t.Run("Testing MaxClaimedHeight", func(t *testing.T) {
		assert.Equal(t, uint32(0), peerSet.MaxClaimedHeight())

		peerSet.UpdateHeight(peer.ID("peer1"), 100)
		peerSet.UpdateHeight(peer.ID("peer2"), 200)
		peerSet.UpdateHeight(peer.ID("peer3"), 150)

		assert.Equal(t, uint32(200), peerSet.MaxClaimedHeight())
	})

	t.Run("Testing GetPeerList", func(t *testing.T) {
		peerList := peerSet.GetPeerList()

		// Verify that the peer list contains the expected peers
		expectedPeerIDs := []peer.ID{"peer1", "peer2", "peer3"}
		for _, expectedID := range expectedPeerIDs {
			found := false
			for _, peer := range peerList {
				if peer.PeerID == expectedID {
					found = true
					break
				}
			}
			assert.True(t, found, "Peer with ID %s not found in the peer list", expectedID)
		}
	})

	t.Run("Testing GetPeer", func(t *testing.T) {
		p := peerSet.GetPeer(peer.ID("peer2"))

		assert.Equal(t, peer.ID("peer2"), p.PeerID)
		assert.Equal(t, StatusCodeKnown, p.Status)

		p = peerSet.GetPeer(peer.ID("unknown"))
		assert.Equal(t, peer.ID(""), p.PeerID)
		assert.Equal(t, StatusCodeUnknown, p.Status)
	})

	t.Run("Testing RemovePeer", func(t *testing.T) {
		peerSet.RemovePeer(peer.ID("unknown"))
		assert.Equal(t, peerSet.Len(), 3)

		peerSet.RemovePeer(peer.ID("peer2"))
		assert.Equal(t, peerSet.Len(), 2)
	})

	t.Run("Clear", func(t *testing.T) {
		peerSet.Clear()

		assert.Equal(t, 0, peerSet.Len())
		assert.Equal(t, uint32(0), peerSet.MaxClaimedHeight())
	})
}

func TestOpenSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	session := ps.OpenSession("peer1")

	assert.NotNil(t, session)

	assert.Equal(t, 1, ps.NumberOfOpenSessions())
}

func TestFindSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)
	session := ps.OpenSession("peer1")

	// Test finding an existing session
	foundSession := ps.FindSession(session.SessionID())

	assert.Equal(t, session, foundSession)

	// Test finding a non-existing session
	nonExistingSession := ps.FindSession(999)

	assert.Nil(t, nonExistingSession)
}

func TestNumberOfOpenSessions(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	// Test when there are no open sessions
	assert.Equal(t, 0, ps.NumberOfOpenSessions())

	// Test when there are multiple open sessions
	ps.OpenSession("peer1")
	ps.OpenSession("peer2")
	ps.OpenSession("peer3")

	assert.Equal(t, 3, ps.NumberOfOpenSessions())
}

func TestHasAnyOpenSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	// Test when there are no open sessions
	assert.False(t, ps.HasAnyOpenSession())

	// Test when there are open sessions
	ps.OpenSession("peer1")

	assert.True(t, ps.HasAnyOpenSession())
}

func TestCloseSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)
	session := ps.OpenSession("peer1")

	// Test closing an existing session
	ps.CloseSession(session.SessionID())

	assert.Equal(t, 0, ps.NumberOfOpenSessions())

	// Test closing a non-existing session
	ps.CloseSession(999)

	assert.Equal(t, 0, ps.NumberOfOpenSessions())
}

func TestGetRandomWeightedPeer(t *testing.T) {
	// We create 6 peers with varying success and failure counts:
	// peer_1 has 4 successful attempts and 0 failed attempts
	// peer_2 has 4 successful attempts and 1 failed attempt
	// ...
	// peer_6 has 4 successful attempts and 4 failed attempts
	// peer_7 has 4 successful attempts and 5 failed attempts
	peerSet := NewPeerSet(time.Second)
	for i := 0; i < 6; i++ {
		pk, _ := bls.GenerateTestKeyPair()
		pid := peer.ID(fmt.Sprintf("peer_%v", i+1))
		peerSet.UpdatePeerInfo(
			pid, StatusCodeKnown,
			fmt.Sprintf("Moniker_%v", i+1), "Agent1", pk, true)

		for s := 0; s < 4; s++ {
			peerSet.IncreaseSendSuccessCounter(pid)
		}
		for f := 0; f < i; f++ {
			peerSet.IncreaseSendFailedCounter(pid)
		}
	}

	// Now let's run TestGetRandomPeer for 1000 times

	hits := make(map[peer.ID]int)
	for i := 0; i < 1000; i++ {
		p := peerSet.GetRandomPeer()
		hits[p.PeerID]++
	}

	assert.Greater(t, hits[peer.ID("peer_1")], hits[peer.ID("peer_3")])
	assert.Greater(t, hits[peer.ID("peer_3")], hits[peer.ID("peer_5")])
	assert.Greater(t, hits[peer.ID("peer_5")], 0)
	assert.Greater(t, hits[peer.ID("peer_6")], 0)
}

func TestGetRandomPeerUnknown(t *testing.T) {
	peerSet := NewPeerSet(time.Second)

	pk, _ := bls.GenerateTestKeyPair()
	pidUnknown := peer.ID("peer_unknown")
	peerSet.UpdatePeerInfo(pidUnknown, StatusCodeUnknown, "Moniker_unknown", "Agent1", pk, true)

	pk, _ = bls.GenerateTestKeyPair()
	pidBanned := peer.ID("peer_banned")
	peerSet.UpdatePeerInfo(pidBanned, StatusCodeBanned, "Moniker_banned", "Agent1", pk, true)

	p := peerSet.GetRandomPeer()

	assert.NotEqual(t, p.PeerID, pidUnknown)
	assert.NotEqual(t, p.PeerID, pidBanned)
}

func TestGetRandomPeerOnePeer(t *testing.T) {
	peerSet := NewPeerSet(time.Second)

	pk, _ := bls.GenerateTestKeyPair()
	pidKnown := peer.ID("peer_known")
	peerSet.UpdatePeerInfo(pidKnown, StatusCodeKnown, "Moniker_known", "Agent1", pk, true)
	peerSet.IncreaseSendSuccessCounter(pidKnown)

	p := peerSet.GetRandomPeer()

	assert.Equal(t, p.PeerID, pidKnown)
}
