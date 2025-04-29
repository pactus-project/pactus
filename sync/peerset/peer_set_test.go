package peerset

import (
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/sync/peerset/session"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func getSessionByID(ps *PeerSet, sid int) *session.Session {
	ssns := ps.Sessions()
	for _, ssn := range ssns {
		if ssn.SessionID == sid {
			return ssn
		}
	}

	return nil
}

func TestPeerSet(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	peerSet := NewPeerSet(time.Minute)

	pk11, _ := ts.RandBLSKeyPair()
	pk12, _ := ts.RandBLSKeyPair()
	pk21, _ := ts.RandBLSKeyPair()
	pk31, _ := ts.RandBLSKeyPair()
	pk32, _ := ts.RandBLSKeyPair()
	pid1 := ts.RandPeerID()
	pid2 := ts.RandPeerID()
	pid3 := ts.RandPeerID()
	peerSet.UpdateInfo(pid1, "Moniker1", "Agent1",
		[]*bls.PublicKey{pk11, pk12}, service.New(service.FullNode))
	peerSet.UpdateInfo(pid2, "Moniker2", "Agent2",
		[]*bls.PublicKey{pk21}, service.New(service.None))
	peerSet.UpdateInfo(pid3, "Moniker3", "Agent3",
		[]*bls.PublicKey{pk31, pk32}, service.New(service.FullNode))

	t.Run("Testing Len", func(t *testing.T) {
		assert.Equal(t, 3, peerSet.Len())
	})

	t.Run("Testing Iterate peers", func(t *testing.T) {
		// Verify that the peer list contains the expected peers
		found := false
		peerSet.IteratePeers(func(p *peer.Peer) bool {
			if p.PeerID == pid2 {
				found = true

				return true
			}

			return false
		})

		assert.True(t, found, "Peer with ID %s not found in the peer list", pid2)
	})

	t.Run("Testing hasPeer", func(t *testing.T) {
		assert.True(t, peerSet.HasPeer(pid1))
		assert.False(t, peerSet.HasPeer(ts.RandPeerID()))
	})

	t.Run("Testing GetPeer", func(t *testing.T) {
		p := peerSet.GetPeer(pid2)
		assert.Equal(t, pid2, p.PeerID)
		assert.True(t, p.Status.IsUnknown())

		p = peerSet.GetPeer(ts.RandPeerID())
		assert.Nil(t, p)
	})

	t.Run("Testing ConsensusKeys", func(t *testing.T) {
		p := peerSet.GetPeer(pid3)

		assert.Contains(t, p.ConsensusKeys, pk31)
		assert.Contains(t, p.ConsensusKeys, pk32)
	})

	t.Run("Testing counters", func(t *testing.T) {
		peerSet.UpdateInvalidMetric(pid1, 123)
		peerSet.UpdateReceivedMetric(pid1, message.TypeBlocksResponse, 100)
		peerSet.UpdateReceivedMetric(pid1, message.TypeTransaction, 150)
		peerSet.UpdateSentMetric(nil, message.TypeBlocksRequest, 200)
		peerSet.UpdateSentMetric(&pid1, message.TypeBlocksRequest, 250)

		peer1 := peerSet.findPeer(pid1)
		assert.Equal(t, int64(1), peer1.Metric.TotalInvalid.Bundles)
		assert.Equal(t, int64(2), peer1.Metric.TotalReceived.Bundles)
		assert.Equal(t, int64(100), peer1.Metric.MessageReceived[message.TypeBlocksResponse].Bytes)
		assert.Equal(t, int64(150), peer1.Metric.MessageReceived[message.TypeTransaction].Bytes)
		assert.Equal(t, int64(250), peer1.Metric.MessageSent[message.TypeBlocksRequest].Bytes)

		peerSetMetric := peerSet.Metric()
		assert.Equal(t, int64(250), peerSetMetric.TotalReceived.Bytes)
		assert.Equal(t, int64(100), peerSetMetric.MessageReceived[message.TypeBlocksResponse].Bytes)
		assert.Equal(t, int64(150), peerSetMetric.MessageReceived[message.TypeTransaction].Bytes)
		assert.Equal(t, int64(450), peerSetMetric.TotalSent.Bytes)
		assert.Equal(t, int64(450), peerSetMetric.MessageSent[message.TypeBlocksRequest].Bytes)
	})

	t.Run("Testing UpdateHeight", func(t *testing.T) {
		height := ts.RandHeight()
		h := ts.RandHash()
		peerSet.UpdateHeight(pid1, height, h)

		peer1 := peerSet.findPeer(pid1)
		assert.Equal(t, height, peer1.Height)
		assert.Equal(t, h, peer1.LastBlockHash)
	})

	t.Run("Testing UpdateStatus", func(t *testing.T) {
		peerSet.UpdateStatus(pid1, status.StatusBanned)
		peerStatus := peerSet.GetPeerStatus(pid1)
		assert.Equal(t, status.StatusBanned, peerStatus)

		peerStatus = peerSet.GetPeerStatus(ts.RandPeerID())
		assert.Equal(t, status.StatusUnknown, peerStatus)
	})

	t.Run("Testing UpdateLastSent", func(t *testing.T) {
		now := time.Now()
		peerSet.UpdateLastSent(pid1)

		peer1 := peerSet.findPeer(pid1)
		assert.GreaterOrEqual(t, peer1.LastSent, now)
	})

	t.Run("Testing UpdateLastReceived", func(t *testing.T) {
		now := time.Now()
		peerSet.UpdateLastReceived(pid1)

		peer1 := peerSet.findPeer(pid1)
		assert.GreaterOrEqual(t, peer1.LastReceived, now)
	})

	t.Run("Testing StartedAt", func(t *testing.T) {
		assert.LessOrEqual(t, peerSet.StartedAt(), time.Now())
	})

	t.Run("Testing RemovePeer", func(t *testing.T) {
		peerSet.RemovePeer(ts.RandPeerID())
		assert.Equal(t, 3, peerSet.Len())

		peerSet.RemovePeer(pid2)
		assert.Equal(t, 2, peerSet.Len())
	})
}

func TestOpenSession(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	peerSet := NewPeerSet(time.Minute)

	pid1 := ts.RandPeerID()
	pid2 := ts.RandPeerID()
	sid1 := peerSet.OpenSession(pid1, 100, 10)
	sid2 := peerSet.OpenSession(pid2, 110, 10)

	ssn1 := getSessionByID(peerSet, sid1)
	ssn2 := getSessionByID(peerSet, sid1)
	assert.NotNil(t, ssn1)
	assert.Equal(t, uint32(100), ssn1.From)
	assert.Equal(t, uint32(100), ssn2.From)
	assert.Equal(t, uint32(10), ssn1.Count)
	assert.Equal(t, uint32(10), ssn2.Count)
	assert.Equal(t, pid1, ssn1.PeerID)
	assert.Equal(t, session.Open, ssn1.Status)
	assert.LessOrEqual(t, ssn1.LastActivity, time.Now())
	assert.Equal(t, 0, sid1)
	assert.Equal(t, 1, sid2)
	assert.True(t, peerSet.HasOpenSession(pid1))
	assert.True(t, peerSet.HasOpenSession(pid2))
	assert.False(t, peerSet.HasOpenSession(ts.RandPeerID()))
	assert.Equal(t, 2, peerSet.NumberOfSessions())
}

func TestNumberOfSessions(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	// Test when there are no open sessions
	assert.Equal(t, 0, peerSet.NumberOfSessions())

	// Test when there are multiple open sessions
	peerSet.OpenSession("peer1", 100, 101)
	peerSet.OpenSession("peer2", 200, 201)
	peerSet.OpenSession("peer3", 300, 301)

	assert.Equal(t, 3, peerSet.NumberOfSessions())
}

func TestHasAnyOpenSession(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	// Test when there are no open sessions
	assert.False(t, peerSet.HasAnyOpenSession())

	sid := peerSet.OpenSession("peer1", 100, 101)
	assert.True(t, peerSet.HasAnyOpenSession())

	peerSet.SetSessionCompleted(sid)
	assert.False(t, peerSet.HasAnyOpenSession())
}

func TestRemoveAllSessions(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	_ = peerSet.OpenSession("peer1", 100, 101)
	_ = peerSet.OpenSession("peer2", 100, 101)
	_ = peerSet.OpenSession("peer3", 100, 101)

	peerSet.RemoveAllSessions()
	assert.Zero(t, peerSet.NumberOfSessions())
	assert.False(t, peerSet.HasAnyOpenSession())
}

func TestCompletedSession(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	sid := peerSet.OpenSession("peer1", 100, 101)
	ssn := getSessionByID(peerSet, sid)
	assert.Equal(t, session.Open, ssn.Status)

	peerSet.SetSessionCompleted(sid)
	assert.Equal(t, 1, peerSet.NumberOfSessions())
	assert.False(t, peerSet.HasAnyOpenSession())
	assert.Equal(t, session.Completed, ssn.Status)
}

func TestUncompletedSession(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	sid := peerSet.OpenSession("peer1", 100, 101)
	ssn := getSessionByID(peerSet, sid)
	assert.Equal(t, session.Open, ssn.Status)

	peerSet.SetSessionUncompleted(sid)
	assert.Equal(t, 1, peerSet.NumberOfSessions())
	assert.False(t, peerSet.HasAnyOpenSession())
	assert.Equal(t, session.Uncompleted, ssn.Status)
}

func TestExpireSessions(t *testing.T) {
	timeout := 100 * time.Millisecond
	peerSet := NewPeerSet(timeout)

	sid := peerSet.OpenSession("peer1", 100, 101)
	ssn := getSessionByID(peerSet, sid)
	time.Sleep(timeout)

	peerSet.SetExpiredSessionsAsUncompleted()
	assert.Equal(t, 1, peerSet.NumberOfSessions())
	assert.False(t, peerSet.HasAnyOpenSession())
	assert.Equal(t, session.Uncompleted, ssn.Status)
}

func TestGetRandomPeer(t *testing.T) {
	// We create 6 peers for testing:
	//
	// peer_1 has score 100
	// peer_2 has score 83
	// peer_3 has score 66
	// peer_4 has score 50
	// peer_5 has score 33
	// peer_6 has score 16
	peerSet := NewPeerSet(time.Minute)
	for index := 0; index < 6; index++ {
		pid := peer.ID(fmt.Sprintf("peer_%v", index+1))
		peerSet.UpdateInfo(pid, fmt.Sprintf("Moniker_%v", index+1), "Agent1", nil, service.New())
		peerSet.UpdateStatus(pid, status.StatusKnown)

		for r := 0; r < 5; r++ {
			sid := peerSet.OpenSession(pid, 0, 0)

			if r < 5-index {
				peerSet.SetSessionCompleted(sid)
			}
		}
	}

	// Now let's run TestGetRandomPeer for 1000 times
	hits := make(map[peer.ID]int)
	for i := 0; i < 1000; i++ {
		randomPeer := peerSet.GetRandomPeer()
		hits[randomPeer.PeerID]++
	}

	assert.Greater(t, hits[peer.ID("peer_1")], hits[peer.ID("peer_3")])
	assert.Greater(t, hits[peer.ID("peer_2")], hits[peer.ID("peer_4")])
	assert.Greater(t, hits[peer.ID("peer_3")], hits[peer.ID("peer_5")])
	assert.Greater(t, hits[peer.ID("peer_4")], hits[peer.ID("peer_6")])
}

func TestGetRandomPeerConnected(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	pidBanned := peer.ID("known")
	pidConnected := peer.ID("connected")
	pidDisconnected := peer.ID("disconnected")
	peerSet.UpdateInfo(pidBanned, "moniker", "agent", nil, service.New())
	peerSet.UpdateInfo(pidConnected, "moniker", "agent", nil, service.New())
	peerSet.UpdateInfo(pidDisconnected, "moniker", "agent", nil, service.New())

	peerSet.UpdateStatus(pidBanned, status.StatusBanned)
	peerSet.UpdateStatus(pidConnected, status.StatusConnected)
	peerSet.UpdateStatus(pidDisconnected, status.StatusDisconnected)

	peer := peerSet.GetRandomPeer()

	assert.NotEqual(t, peer.PeerID, pidBanned)
	assert.NotEqual(t, peer.PeerID, pidDisconnected)
	assert.Equal(t, peer.PeerID, pidConnected)
}

func TestGetRandomPeerNoPeer(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	randomPeer := peerSet.GetRandomPeer()
	assert.Nil(t, randomPeer)
}

func TestGetRandomPeerOnePeer(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	pidAlice := peer.ID("alice")
	peerSet.UpdateInfo(pidAlice, "alice", "agent", nil, service.New())
	peerSet.UpdateStatus(pidAlice, status.StatusKnown)

	randomPeer := peerSet.GetRandomPeer()
	assert.Equal(t, randomPeer.PeerID, pidAlice)
}

func TestUpdateAddress(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	pid := peer.ID("peer1")
	addr := "pid-1-address"
	dir := "Inbound"
	peerSet.UpdateAddress(pid, addr, dir)

	p := peerSet.GetPeer(pid)
	assert.Equal(t, addr, p.Address)
	assert.Equal(t, dir, p.Direction)
}

func TestUpdateSessionLastActivity(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	sid := peerSet.OpenSession("peer1", 100, 101)
	ssn := getSessionByID(peerSet, sid)
	activity1 := ssn.LastActivity
	time.Sleep(10 * time.Millisecond)
	peerSet.UpdateSessionLastActivity(sid)
	assert.Greater(t, ssn.LastActivity, activity1)
}

func TestUpdateProtocols(t *testing.T) {
	peerSet := NewPeerSet(time.Minute)

	pid := peer.ID("peer-1")
	protocols := []string{"protocol-1"}
	peerSet.UpdateProtocols(pid, protocols)

	p := peerSet.GetPeer(pid)
	assert.Equal(t, protocols, p.Protocols)
}
