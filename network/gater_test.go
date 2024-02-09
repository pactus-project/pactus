package network

import (
	"testing"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

type mockConnMultiaddrs struct {
	remote multiaddr.Multiaddr
}

func (cma *mockConnMultiaddrs) LocalMultiaddr() multiaddr.Multiaddr {
	return nil
}

func (cma *mockConnMultiaddrs) RemoteMultiaddr() multiaddr.Multiaddr {
	return cma.remote
}

func TestAllowedConnections(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	conf := testConfig()
	net := makeTestNetwork(t, conf, nil)

	maPrivate := multiaddr.StringCast("/ip4/127.0.0.1/tcp/1234")
	maPublic := multiaddr.StringCast("/ip4/8.8.8.8/tcp/1234")
	cmaPrivate := &mockConnMultiaddrs{remote: maPrivate}
	cmaPublic := &mockConnMultiaddrs{remote: maPublic}
	pid := ts.RandPeerID()

	assert.True(t, net.connGater.InterceptPeerDial(pid))
	assert.True(t, net.connGater.InterceptAddrDial(pid, maPrivate))
	assert.True(t, net.connGater.InterceptAddrDial(pid, maPublic))
	assert.True(t, net.connGater.InterceptAccept(cmaPrivate))
	assert.True(t, net.connGater.InterceptAccept(cmaPublic))
}

func TestDenyPrivate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	conf := testConfig()
	conf.ForcePrivateNetwork = false
	net := makeTestNetwork(t, conf, nil)

	maPrivate := multiaddr.StringCast("/ip4/127.0.0.1/tcp/1234")
	maPublic := multiaddr.StringCast("/ip4/8.8.8.8/tcp/1234")
	cmaPrivate := &mockConnMultiaddrs{remote: maPrivate}
	cmaPublic := &mockConnMultiaddrs{remote: maPublic}
	pid := ts.RandPeerID()

	assert.True(t, net.connGater.InterceptPeerDial(pid))
	assert.False(t, net.connGater.InterceptAddrDial(pid, maPrivate))
	assert.True(t, net.connGater.InterceptAddrDial(pid, maPublic))
	assert.False(t, net.connGater.InterceptAccept(cmaPrivate))
	assert.True(t, net.connGater.InterceptAccept(cmaPublic))
}

func TestMaxConnection(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	conf := testConfig()
	conf.MaxConns = 4
	assert.Equal(t, conf.ScaledMaxConns(), 4)
	net := makeTestNetwork(t, conf, nil)

	maPrivate := multiaddr.StringCast("/ip4/127.0.0.1/tcp/1234")
	maPublic := multiaddr.StringCast("/ip4/8.8.8.8/tcp/1234")
	aMultiAddr := multiaddr.StringCast("/ip4/1.1.1.1/tcp/1234")
	cmaPrivate := &mockConnMultiaddrs{remote: maPrivate}
	cmaPublic := &mockConnMultiaddrs{remote: maPublic}
	pid := ts.RandPeerID()

	net.peerMgr.AddPeer(ts.RandPeerID(), aMultiAddr, lp2pnetwork.DirOutbound)
	net.peerMgr.AddPeer(ts.RandPeerID(), aMultiAddr, lp2pnetwork.DirInbound)
	net.peerMgr.AddPeer(ts.RandPeerID(), aMultiAddr, lp2pnetwork.DirInbound)
	net.peerMgr.AddPeer(ts.RandPeerID(), aMultiAddr, lp2pnetwork.DirInbound)
	net.peerMgr.AddPeer(ts.RandPeerID(), aMultiAddr, lp2pnetwork.DirInbound)

	assert.True(t, net.connGater.InterceptPeerDial(pid))
	assert.True(t, net.connGater.InterceptAddrDial(pid, maPrivate))
	assert.True(t, net.connGater.InterceptAddrDial(pid, maPublic))
	assert.True(t, net.connGater.InterceptAccept(cmaPrivate))
	assert.True(t, net.connGater.InterceptAccept(cmaPublic))

	net.peerMgr.AddPeer(ts.RandPeerID(), aMultiAddr, lp2pnetwork.DirOutbound)

	assert.False(t, net.connGater.InterceptPeerDial(pid))
	assert.False(t, net.connGater.InterceptAddrDial(pid, maPrivate))
	assert.False(t, net.connGater.InterceptAddrDial(pid, maPublic))
	assert.True(t, net.connGater.InterceptAccept(cmaPrivate))
	assert.True(t, net.connGater.InterceptAccept(cmaPublic))

	net.peerMgr.AddPeer(ts.RandPeerID(), aMultiAddr, lp2pnetwork.DirInbound)

	assert.False(t, net.connGater.InterceptPeerDial(pid))
	assert.False(t, net.connGater.InterceptAddrDial(pid, maPrivate))
	assert.False(t, net.connGater.InterceptAddrDial(pid, maPublic))
	assert.False(t, net.connGater.InterceptAccept(cmaPrivate))
	assert.False(t, net.connGater.InterceptAccept(cmaPublic))
}
