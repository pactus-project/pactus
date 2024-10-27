package peer

import (
	"testing"

	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/stretchr/testify/assert"
)

func TestPeerStatus(t *testing.T) {
	tests := []struct {
		status             status.Status
		isDisconnected     bool
		isBanned           bool
		isConnected        bool
		isKnown            bool
		isKnownOrConnected bool
	}{
		{status: status.StatusBanned, isBanned: true},
		{status: status.StatusDisconnected, isDisconnected: true},
		{status: status.StatusConnected, isConnected: true, isKnownOrConnected: true},
		{status: status.StatusKnown, isKnown: true, isKnownOrConnected: true},
	}

	for _, tt := range tests {
		peer := NewPeer("test")
		peer.Status = tt.status

		assert.Equal(t, tt.isDisconnected, peer.Status.IsDisconnected())
		assert.Equal(t, tt.isBanned, peer.Status.IsBanned())
		assert.Equal(t, tt.isConnected, peer.Status.IsConnected())
		assert.Equal(t, tt.isKnown, peer.Status.IsKnown())
		assert.Equal(t, tt.isKnownOrConnected, peer.Status.IsConnectedOrKnown())
	}
}

func TestIsFullNode(t *testing.T) {
	p1 := NewPeer("peer-1")
	p2 := NewPeer("peer-1")
	p1.Services = service.New(service.PrunedNode)
	p2.Services = service.New(service.FullNode)

	assert.False(t, p1.IsFullNode())
	assert.True(t, p2.IsFullNode())
}

func TestDownloadScore(t *testing.T) {
	tests := []struct {
		totalSession     int
		completedSession int
		expectedScore    int
	}{
		{0, 0, 100},
		{1, 1, 100},
		{1, 0, 50},
		{2, 0, 33},
		{3, 0, 25},
		{5, 4, 83},
		{6, 5, 85},
		{7, 6, 87},
	}

	for no, tt := range tests {
		peer := NewPeer("test")
		peer.TotalSessions = tt.totalSession
		peer.CompletedSessions = tt.completedSession

		score := peer.DownloadScore()
		assert.Equal(t, tt.expectedScore, score,
			"Test %v failed. expected score %d, got %d",
			no+1, tt.expectedScore, score)
	}
}
