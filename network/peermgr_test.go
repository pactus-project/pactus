package network

import (
	"testing"

	lp2pnet "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestGetMultiAddr(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conf := testConfig()
	net := makeTestNetwork(t, conf, nil)

	pid := ts.RandPeerID()
	addr, _ := IPToMultiAddr("1.2.3.4", 1234)
	assert.Nil(t, net.peerMgr.GetMultiAddr(pid))

	net.peerMgr.AddPeer(pid, addr, lp2pnet.DirOutbound)
	assert.Equal(t, addr, net.peerMgr.GetMultiAddr(pid))

	net.peerMgr.RemovePeer(pid)
	assert.Nil(t, net.peerMgr.GetMultiAddr(pid))
}

func TestNumInboundOutbound(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conf := testConfig()
	net := makeTestNetwork(t, conf, nil)

	addr, _ := IPToMultiAddr("1.2.3.4", 1234)

	pid1 := ts.RandPeerID()
	pid2 := ts.RandPeerID()
	pid3 := ts.RandPeerID()

	net.peerMgr.AddPeer(pid1, addr, lp2pnet.DirInbound)
	net.peerMgr.AddPeer(pid2, addr, lp2pnet.DirOutbound)
	net.peerMgr.AddPeer(pid3, addr, lp2pnet.DirOutbound)
	// Adding pid1 again
	net.peerMgr.AddPeer(pid1, addr, lp2pnet.DirInbound)

	assert.Equal(t, 1, net.peerMgr.NumInbound())
	assert.Equal(t, 2, net.peerMgr.NumOutbound())

	net.peerMgr.RemovePeer(pid1)
	net.peerMgr.RemovePeer(pid2)
	net.peerMgr.RemovePeer(ts.RandPeerID())

	assert.Equal(t, 0, net.peerMgr.NumInbound())
	assert.Equal(t, 1, net.peerMgr.NumOutbound())

	assert.Equal(t, net.NumInbound(), net.peerMgr.NumInbound())
	assert.Equal(t, net.NumOutbound(), net.peerMgr.NumOutbound())
}
