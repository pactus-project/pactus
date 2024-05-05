package network

import (
	"testing"

	lp2pnet "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestNumInboundOutbound(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conf := testConfig()
	net := makeTestNetwork(t, conf, nil)

	addr, _ := IPToMultiAddr("1.2.3.4", 1234)

	pid1 := ts.RandPeerID()
	pid2 := ts.RandPeerID()
	pid3 := ts.RandPeerID()

	net.peerMgr.peerConnected(pid1, addr, lp2pnet.DirInbound)
	net.peerMgr.peerConnected(pid1, addr, lp2pnet.DirInbound) // Duplicated event
	net.peerMgr.peerConnected(pid2, addr, lp2pnet.DirOutbound)
	net.peerMgr.peerConnected(pid3, addr, lp2pnet.DirOutbound)
	net.peerMgr.PeerDisconnected(pid1)
	net.peerMgr.PeerDisconnected(pid1) // Duplicated event
	net.peerMgr.PeerDisconnected(pid2)
	net.peerMgr.PeerDisconnected(ts.RandPeerID())
	net.peerMgr.peerConnected(pid1, addr, lp2pnet.DirInbound) // Connect again

	assert.Equal(t, 1, net.NumInbound())
	assert.Equal(t, 1, net.NumOutbound())
}
