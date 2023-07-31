package peerset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/libp2p/go-libp2p/core/peer"

)

func TestPeer(t *testing.T) {
	pid1 := peer.ID("peer-1")
	pid2 := peer.ID("peer-2")
	pid3 := peer.ID("peer-3")

	p1 := NewPeer(pid1)
	p2 := NewPeer(pid2)
	p3 := NewPeer(pid3)

	p2.Status = StatusCodeKnown
	p3.Status = StatusCodeBanned
	
	p1.Flags = PeerFlagNodeNetwork

	t.Run("NewPeer", func(t *testing.T) {
		assert.NotNil(t, p1)
		assert.Equal(t, p1.Status, StatusCodeUnknown)
	})

	t.Run("status check", func(t *testing.T) {
		unknown :=  p1.IsKnownOrTrusty()
		known := p2.IsKnownOrTrusty()
		banned := p3.IsBanned()

		assert.False(t, unknown)
		assert.True(t, known)
		assert.True(t, banned)
	})

	t.Run("is node network", func(t *testing.T) {
		nodeNetwork := p1.IsNodeNetwork()
		notNodeNetwork := p2.IsNodeNetwork()

		assert.True(t, nodeNetwork)
		assert.False(t, notNodeNetwork)
	})
}

