package network

import (
	"math"
	"testing"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2prcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceLimitsAreBounded(t *testing.T) {
	conf := DefaultConfig()
	conf.MaxConns = 1000

	limits := resourceLimits(conf)
	system := limits.ToPartialLimitConfig().System

	assert.Less(t, int(system.ConnsInbound), math.MaxInt, "inbound connections should be capped")
	assert.GreaterOrEqual(t, int(system.ConnsInbound), conf.MaxConns+conf.MaxConns/4,
		"the cap should not be below the connection manager high watermark")
	assert.Less(t, int(system.Conns), math.MaxInt)
}

// The resource manager of the host uses the bounded limits.
func TestHostResourceManagerIsBounded(t *testing.T) {
	conf := testConfig()
	net := makeTestNetwork(t, conf, nil)

	type limited interface {
		Limit() lp2prcmgr.Limit
	}

	checked := false
	err := net.host.Network().ResourceManager().ViewSystem(func(scope lp2pnetwork.ResourceScope) error {
		limits, ok := scope.(limited)
		require.True(t, ok, "the system scope has limits")

		assert.Less(t, limits.Limit().GetConnTotalLimit(), math.MaxInt)
		assert.Less(t, limits.Limit().GetConnLimit(lp2pnetwork.DirInbound), math.MaxInt)
		checked = true

		return nil
	})
	require.NoError(t, err)
	assert.True(t, checked)
}
