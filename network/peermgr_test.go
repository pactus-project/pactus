package network

import (
	"testing"

	lp2pnet "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestGetMultiAddr(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conf := testConfig()
	net := makeTestNetwork(t, conf, nil)

	pid := ts.RandPeerID()
	addr, _ := IPToMultiAddr("1.2.3.4", 1234)
	addrInfo, err := lp2ppeer.AddrInfoFromP2pAddr(addr)
	assert.NoError(t, err)
	assert.Nil(t, net.peerMgr.GetMultiAddr(pid))

	net.peerMgr.AddPeer(pid, *addrInfo, lp2pnet.DirOutbound)
	assert.Equal(t, addr, net.peerMgr.GetMultiAddr(pid))

	net.peerMgr.RemovePeer(pid)
	assert.Nil(t, net.peerMgr.GetMultiAddr(pid))
}
