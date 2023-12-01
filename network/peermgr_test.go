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
