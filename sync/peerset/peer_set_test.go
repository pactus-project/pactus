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

	ps := NewPeerSet(time.Minute)

	pk1, _ := ts.RandBLSKeyPair()
	pk2, _ := ts.RandBLSKeyPair()
	pk3, _ := ts.RandBLSKeyPair()
	pk4, _ := ts.RandBLSKeyPair()
	pk5, _ := ts.RandBLSKeyPair()
	pid1 := ts.RandPeerID()
	pid2 := ts.RandPeerID()
	pid3 := ts.RandPeerID()
	ps.UpdateInfo(pid1, "Moniker1", "Agent1",
		[]*bls.PublicKey{pk1, pk2}, service.New(service.FullNode))
	ps.UpdateInfo(pid2, "Moniker2", "Agent2",
		[]*bls.PublicKey{pk3}, service.New(service.None))
	ps.UpdateInfo(pid3, "Moniker3", "Agent3",
		[]*bls.PublicKey{pk4, pk5}, service.New(service.FullNode))

	t.Run("Testing Len", func(t *testing.T) {
		assert.Equal(t, 3, ps.Len())
	})

	t.Run("Testing Iterate peers", func(t *testing.T) {
		// Verify that the peer list contains the expected peers
		found := false
		ps.IteratePeers(func(p *peer.Peer) bool {
			if p.PeerID == pid2 {
				found = true

				return true
			}

			return false
		})

		assert.True(t, found, "Peer with ID %s not found in the peer list", pid2)
	})

	t.Run("Testing GetPeer", func(t *testing.T) {
		p := ps.GetPeer(pid2)
		assert.Equal(t, pid2, p.PeerID)
		assert.True(t, p.Status.IsUnknown())

		p = ps.GetPeer(ts.RandPeerID())
		assert.Nil(t, p)
	})

	t.Run("Testing ConsensusKeys", func(t *testing.T) {
		p := ps.GetPeer(pid3)

		assert.Contains(t, p.ConsensusKeys, pk4)
		assert.Contains(t, p.ConsensusKeys, pk5)
	})

	t.Run("Testing counters", func(t *testing.T) {
		ps.IncreaseInvalidBundlesCounter(pid1)
		ps.IncreaseReceivedBundlesCounter(pid1)
		ps.IncreaseReceivedBytesCounter(pid1, message.TypeBlocksResponse, 100)
		ps.IncreaseReceivedBytesCounter(pid1, message.TypeTransaction, 150)
		ps.IncreaseSentCounters(message.TypeBlocksRequest, 200, nil)
		ps.IncreaseSentCounters(message.TypeBlocksRequest, 250, &pid1)

		peer1 := ps.findPeer(pid1)

		receivedBytes := make(map[message.Type]int64)
		receivedBytes[message.TypeBlocksResponse] = 100
		receivedBytes[message.TypeTransaction] = 150

		sentBytes := make(map[message.Type]int64)
		sentBytes[message.TypeBlocksRequest] = 450

		assert.Equal(t, 1, peer1.InvalidBundles)
		assert.Equal(t, 1, peer1.ReceivedBundles)
		assert.Equal(t, int64(100), peer1.ReceivedBytes[message.TypeBlocksResponse])
		assert.Equal(t, int64(150), peer1.ReceivedBytes[message.TypeTransaction])
		assert.Equal(t, int64(250), peer1.SentBytes[message.TypeBlocksRequest])

		assert.Equal(t, int64(250), ps.TotalReceivedBytes())
		assert.Equal(t, int64(100), ps.ReceivedBytesMessageType(message.TypeBlocksResponse))
		assert.Equal(t, int64(150), ps.ReceivedBytesMessageType(message.TypeTransaction))
		assert.Equal(t, receivedBytes, ps.ReceivedBytes())
		assert.Equal(t, int64(450), ps.TotalSentBytes())
		assert.Equal(t, int64(450), ps.SentBytesMessageType(message.TypeBlocksRequest))
		assert.Equal(t, sentBytes, ps.SentBytes())
		assert.Equal(t, 2, ps.TotalSentBundles())
	})

	t.Run("Testing UpdateHeight", func(t *testing.T) {
		height := ts.RandHeight()
		h := ts.RandHash()
		ps.UpdateHeight(pid1, height, h)

		peer1 := ps.findPeer(pid1)
		assert.Equal(t, height, peer1.Height)
		assert.Equal(t, h, peer1.LastBlockHash)
	})

	t.Run("Testing UpdateStatus", func(t *testing.T) {
		ps.UpdateStatus(pid1, status.StatusBanned)

		peer1 := ps.findPeer(pid1)
		assert.Equal(t, status.StatusBanned, peer1.Status)
	})

	t.Run("Testing UpdateLastSent", func(t *testing.T) {
		now := time.Now()
		ps.UpdateLastSent(pid1)

		peer1 := ps.findPeer(pid1)
		assert.GreaterOrEqual(t, peer1.LastSent, now)
	})

	t.Run("Testing UpdateLastReceived", func(t *testing.T) {
		now := time.Now()
		ps.UpdateLastReceived(pid1)

		peer1 := ps.findPeer(pid1)
		assert.GreaterOrEqual(t, peer1.LastReceived, now)
	})

	t.Run("Testing StartedAt", func(t *testing.T) {
		assert.LessOrEqual(t, ps.StartedAt(), time.Now())
	})

	t.Run("Testing RemovePeer", func(t *testing.T) {
		ps.RemovePeer(ts.RandPeerID())
		assert.Equal(t, 3, ps.Len())

		ps.RemovePeer(pid2)
		assert.Equal(t, 2, ps.Len())
	})
}

func TestOpenSession(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	ps := NewPeerSet(time.Minute)

	pid1 := ts.RandPeerID()
	pid2 := ts.RandPeerID()
	sid1 := ps.OpenSession(pid1, 100, 10)
	sid2 := ps.OpenSession(pid2, 110, 10)

	ssn1 := getSessionByID(ps, sid1)
	ssn2 := getSessionByID(ps, sid1)
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
	assert.True(t, ps.HasOpenSession(pid1))
	assert.True(t, ps.HasOpenSession(pid2))
	assert.False(t, ps.HasOpenSession(ts.RandPeerID()))
	assert.Equal(t, 2, ps.NumberOfSessions())
}

func TestNumberOfSessions(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	// Test when there are no open sessions
	assert.Equal(t, 0, ps.NumberOfSessions())

	// Test when there are multiple open sessions
	ps.OpenSession("peer1", 100, 101)
	ps.OpenSession("peer2", 200, 201)
	ps.OpenSession("peer3", 300, 301)

	assert.Equal(t, 3, ps.NumberOfSessions())
}

func TestHasAnyOpenSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	// Test when there are no open sessions
	assert.False(t, ps.HasAnyOpenSession())

	sid := ps.OpenSession("peer1", 100, 101)
	assert.True(t, ps.HasAnyOpenSession())

	ps.SetSessionCompleted(sid)
	assert.False(t, ps.HasAnyOpenSession())
}

func TestRemoveAllSessions(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	_ = ps.OpenSession("peer1", 100, 101)
	_ = ps.OpenSession("peer2", 100, 101)
	_ = ps.OpenSession("peer3", 100, 101)

	ps.RemoveAllSessions()
	assert.Zero(t, ps.NumberOfSessions())
	assert.False(t, ps.HasAnyOpenSession())
}

func TestCompletedSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	sid := ps.OpenSession("peer1", 100, 101)
	ssn := getSessionByID(ps, sid)
	assert.Equal(t, session.Open, ssn.Status)

	ps.SetSessionCompleted(sid)
	assert.Equal(t, 1, ps.NumberOfSessions())
	assert.False(t, ps.HasAnyOpenSession())
	assert.Equal(t, session.Completed, ssn.Status)
}

func TestUncompletedSession(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	sid := ps.OpenSession("peer1", 100, 101)
	ssn := getSessionByID(ps, sid)
	assert.Equal(t, session.Open, ssn.Status)

	ps.SetSessionUncompleted(sid)
	assert.Equal(t, 1, ps.NumberOfSessions())
	assert.False(t, ps.HasAnyOpenSession())
	assert.Equal(t, session.Uncompleted, ssn.Status)
}

func TestExpireSessions(t *testing.T) {
	timeout := 100 * time.Millisecond
	ps := NewPeerSet(timeout)

	sid := ps.OpenSession("peer1", 100, 101)
	ssn := getSessionByID(ps, sid)
	time.Sleep(timeout)

	ps.SetExpiredSessionsAsUncompleted()
	assert.Equal(t, 1, ps.NumberOfSessions())
	assert.False(t, ps.HasAnyOpenSession())
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
	ps := NewPeerSet(time.Minute)
	for i := 0; i < 6; i++ {
		pid := peer.ID(fmt.Sprintf("peer_%v", i+1))
		ps.UpdateInfo(pid, fmt.Sprintf("Moniker_%v", i+1), "Agent1", nil, service.New())
		ps.UpdateStatus(pid, status.StatusKnown)

		for r := 0; r < 5; r++ {
			sid := ps.OpenSession(pid, 0, 0)

			if r < 5-i {
				ps.SetSessionCompleted(sid)
			}
		}
	}

	// Now let's run TestGetRandomPeer for 1000 times
	hits := make(map[peer.ID]int)
	for i := 0; i < 1000; i++ {
		randomPeer := ps.GetRandomPeer()
		hits[randomPeer.PeerID]++
	}

	assert.Greater(t, hits[peer.ID("peer_1")], hits[peer.ID("peer_3")])
	assert.Greater(t, hits[peer.ID("peer_2")], hits[peer.ID("peer_4")])
	assert.Greater(t, hits[peer.ID("peer_3")], hits[peer.ID("peer_5")])
	assert.Greater(t, hits[peer.ID("peer_4")], hits[peer.ID("peer_6")])
}

func TestGetRandomPeerConnected(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	pidBanned := peer.ID("known")
	pidConnected := peer.ID("connected")
	pidDisconnected := peer.ID("disconnected")
	ps.UpdateInfo(pidBanned, "moniker", "agent", nil, service.New())
	ps.UpdateInfo(pidConnected, "moniker", "agent", nil, service.New())
	ps.UpdateInfo(pidDisconnected, "moniker", "agent", nil, service.New())

	ps.UpdateStatus(pidBanned, status.StatusBanned)
	ps.UpdateStatus(pidConnected, status.StatusConnected)
	ps.UpdateStatus(pidDisconnected, status.StatusDisconnected)

	p := ps.GetRandomPeer()

	assert.NotEqual(t, p.PeerID, pidBanned)
	assert.NotEqual(t, p.PeerID, pidDisconnected)
	assert.Equal(t, p.PeerID, pidConnected)
}

func TestGetRandomPeerNoPeer(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	randomPeer := ps.GetRandomPeer()
	assert.Nil(t, randomPeer)
}

func TestGetRandomPeerOnePeer(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	pidAlice := peer.ID("alice")
	ps.UpdateInfo(pidAlice, "alice", "agent", nil, service.New())
	ps.UpdateStatus(pidAlice, status.StatusKnown)

	randomPeer := ps.GetRandomPeer()
	assert.Equal(t, randomPeer.PeerID, pidAlice)
}

func TestUpdateAddress(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	pid := peer.ID("peer1")
	addr := "pid-1-address"
	dir := "Inbound"
	ps.UpdateAddress(pid, addr, dir)

	p := ps.GetPeer(pid)
	assert.Equal(t, addr, p.Address)
	assert.Equal(t, dir, p.Direction)
}

func TestUpdateSessionLastActivity(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	sid := ps.OpenSession("peer1", 100, 101)
	ssn := getSessionByID(ps, sid)
	activity1 := ssn.LastActivity
	time.Sleep(10 * time.Millisecond)
	ps.UpdateSessionLastActivity(sid)
	assert.Greater(t, ssn.LastActivity, activity1)
}

func TestUpdateProtocols(t *testing.T) {
	ps := NewPeerSet(time.Minute)

	pid := peer.ID("peer-1")
	protocols := []string{"protocol-1"}
	ps.UpdateProtocols(pid, protocols)

	p := ps.GetPeer(pid)
	assert.Equal(t, protocols, p.Protocols)
}

func TestUpdateStatus(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	tests := []struct {
		name           string
		initialStatus  status.Status
		newStatus      status.Status
		expectedStatus status.Status
	}{
		{
			name:           "Peer is banned, status should not change (attempting to connect)",
			initialStatus:  status.StatusBanned,
			newStatus:      status.StatusConnected,
			expectedStatus: status.StatusBanned,
		},
		{
			name:           "Peer is banned, status should not change (attempting to disconnect)",
			initialStatus:  status.StatusBanned,
			newStatus:      status.StatusDisconnected,
			expectedStatus: status.StatusBanned,
		},
		{
			name:           "Peer is banned, status should not change (attempting to set known)",
			initialStatus:  status.StatusBanned,
			newStatus:      status.StatusKnown,
			expectedStatus: status.StatusBanned,
		},
		{
			name:           "Peer is known, trying to change status to connected, should not change",
			initialStatus:  status.StatusKnown,
			newStatus:      status.StatusConnected,
			expectedStatus: status.StatusKnown,
		},
		{
			name:           "Peer is known, changing status to disconnected",
			initialStatus:  status.StatusKnown,
			newStatus:      status.StatusDisconnected,
			expectedStatus: status.StatusDisconnected,
		},
		{
			name:           "Updating unknown status to connected",
			initialStatus:  status.StatusUnknown,
			newStatus:      status.StatusConnected,
			expectedStatus: status.StatusConnected,
		},
		{
			name:           "Updating connected status to disconnected",
			initialStatus:  status.StatusConnected,
			newStatus:      status.StatusDisconnected,
			expectedStatus: status.StatusDisconnected,
		},
	}

	ps := NewPeerSet(time.Minute)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid := ts.RandPeerID()

			ps.UpdateStatus(pid, tt.initialStatus)
			ps.UpdateStatus(pid, tt.newStatus)

			actualStatus := ps.GetPeerStatus(pid)
			assert.Equal(t, tt.expectedStatus, actualStatus)
		})
	}
}
