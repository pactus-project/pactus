package peerset

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/services"
	"github.com/stretchr/testify/assert"
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

	p1.Services = services.New(services.Network)

	t.Run("NewPeer", func(t *testing.T) {
		assert.NotNil(t, p1)
		assert.Equal(t, p1.Status, StatusCodeUnknown)
	})

	t.Run("status check", func(t *testing.T) {
		unknown := p1.IsKnownOrTrusty()
		known := p2.IsKnownOrTrusty()
		banned := p3.IsBanned()

		assert.False(t, unknown)
		assert.True(t, known)
		assert.True(t, banned)
	})

	t.Run("has network service", func(t *testing.T) {
		assert.True(t, p1.HasNetworkService())
		assert.False(t, p2.HasNetworkService())
	})
}
