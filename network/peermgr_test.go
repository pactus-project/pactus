package network

import (
	"encoding/json"
	"testing"

	lp2pnet "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumInboundOutbound(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conf := testConfig()
	net := makeTestNetwork(t, conf, nil)

	addr, _ := IPToMultiAddr("1.2.3.4", 1234)

	pid1 := ts.RandPeerID()
	pid2 := ts.RandPeerID()
	pid3 := ts.RandPeerID()

	net.peerMgr.SetPeerConnected(pid1, addr, lp2pnet.DirInbound)
	net.peerMgr.SetPeerConnected(pid1, addr, lp2pnet.DirInbound) // Duplicated event
	net.peerMgr.SetPeerConnected(pid2, addr, lp2pnet.DirOutbound)
	net.peerMgr.SetPeerConnected(pid3, addr, lp2pnet.DirOutbound)
	net.peerMgr.SetPeerDisconnected(pid1)
	net.peerMgr.SetPeerDisconnected(pid1) // Duplicated event
	net.peerMgr.SetPeerDisconnected(pid2)
	net.peerMgr.SetPeerDisconnected(ts.RandPeerID())
	net.peerMgr.SetPeerConnected(pid1, addr, lp2pnet.DirInbound) // Connect again

	assert.Equal(t, 1, net.NumInbound())
	assert.Equal(t, 1, net.NumOutbound())
}

func TestNewPeerMgrIgnoresPeerWithoutAddress(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conf := testConfig()
	conf.BootstrapAddrStrings = []string{"/p2p/" + ts.RandPeerID().String()}

	// A bootstrap address without a transport address is skipped instead of
	// stopping the node.
	net := makeTestNetwork(t, conf, nil)

	assert.Empty(t, net.peerMgr.peers)
}

func TestLoadPeerStore(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	validAddr1 := "/ip4/127.0.0.1/tcp/1234/p2p/" + ts.RandPeerID().String()
	validAddr2 := "/ip4/127.0.0.2/tcp/1234/p2p/" + ts.RandPeerID().String()

	t.Run("an invalid entry does not drop the valid ones", func(t *testing.T) {
		path := util.TempFilePath()
		data, err := json.Marshal([]string{validAddr1, "invalid", "/ip4/127.0.0.3/tcp/1234", validAddr2})
		require.NoError(t, err)
		require.NoError(t, util.WriteFile(path, data))

		peerStore, err := loadPeerStore(path)
		require.Error(t, err)
		require.Len(t, peerStore, 2)
		assert.Equal(t, validAddr1, peerStore[0].Addrs[0].String()+"/p2p/"+peerStore[0].ID.String())
		assert.Equal(t, validAddr2, peerStore[1].Addrs[0].String()+"/p2p/"+peerStore[1].ID.String())
	})

	t.Run("valid entries only", func(t *testing.T) {
		path := util.TempFilePath()
		data, err := json.Marshal([]string{validAddr1, validAddr2})
		require.NoError(t, err)
		require.NoError(t, util.WriteFile(path, data))

		peerStore, err := loadPeerStore(path)
		require.NoError(t, err)
		assert.Len(t, peerStore, 2)
	})

	t.Run("missing file", func(t *testing.T) {
		peerStore, err := loadPeerStore(util.TempFilePath())
		require.Error(t, err)
		assert.Empty(t, peerStore)
	})
}
