package firewall

import (
	"bytes"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	*testsuite.TestSuite

	firewall      *Firewall
	badPeerID     peer.ID
	goodPeerID    peer.ID
	unknownPeerID peer.ID
	network       *network.MockNetwork
	state         *state.MockState
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	subLogger := logger.NewSubLogger("firewall", nil)
	peerSet := peerset.NewPeerSet(1 * time.Minute)
	st := state.MockingState(ts)
	net := network.MockingNetwork(ts, ts.RandPeerID())
	conf := DefaultConfig()
	conf.Enabled = true
	firewall := NewFirewall(conf, net, peerSet, st, subLogger)
	assert.NotNil(t, firewall)
	badPeerID := ts.RandPeerID()
	goodPeerID := ts.RandPeerID()
	unknownPeerID := ts.RandPeerID()

	net.AddAnotherNetwork(network.MockingNetwork(ts, goodPeerID))
	net.AddAnotherNetwork(network.MockingNetwork(ts, unknownPeerID))
	net.AddAnotherNetwork(network.MockingNetwork(ts, badPeerID))

	firewall.peerSet.UpdateStatus(goodPeerID, peerset.StatusCodeKnown)
	firewall.peerSet.UpdateStatus(badPeerID, peerset.StatusCodeBanned)

	return &testData{
		TestSuite:     ts,
		firewall:      firewall,
		network:       net,
		state:         st,
		badPeerID:     badPeerID,
		goodPeerID:    goodPeerID,
		unknownPeerID: unknownPeerID,
	}
}

func TestInvalidBundlesCounter(t *testing.T) {
	td := setup(t)

	assert.Nil(t, td.firewall.OpenGossipBundle([]byte("bad"), td.unknownPeerID))
	assert.Nil(t, td.firewall.OpenGossipBundle(nil, td.unknownPeerID))

	bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), -1, td.RandValAddress()))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
	d, _ := bdl.Encode()
	assert.Nil(t, td.firewall.OpenGossipBundle(d, td.unknownPeerID))

	p := td.firewall.peerSet.GetPeer(td.unknownPeerID)
	assert.Equal(t, p.InvalidBundles, 3)
}

func TestGossipMessage(t *testing.T) {
	t.Run("Message  from: unknown => should NOT close the connection", func(t *testing.T) {
		td := setup(t)

		bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
		d, _ := bdl.Encode()

		assert.False(t, td.network.IsClosed(td.unknownPeerID))
		assert.NotNil(t, td.firewall.OpenGossipBundle(d, td.unknownPeerID))
		assert.False(t, td.network.IsClosed(td.unknownPeerID))
	})

	t.Run("Message  from: bad => should close the connection", func(t *testing.T) {
		td := setup(t)

		bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
		d, _ := bdl.Encode()

		assert.False(t, td.network.IsClosed(td.badPeerID))
		assert.Nil(t, td.firewall.OpenGossipBundle(d, td.badPeerID))
		assert.True(t, td.network.IsClosed(td.badPeerID))
	})

	t.Run("Message is nil => should close the connection", func(t *testing.T) {
		td := setup(t)

		assert.Nil(t, td.firewall.OpenGossipBundle(nil, td.unknownPeerID))
	})

	t.Run("Ok => should NOT close the connection", func(t *testing.T) {
		td := setup(t)

		bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
		d, _ := bdl.Encode()

		assert.False(t, td.network.IsClosed(td.goodPeerID))
		assert.NotNil(t, td.firewall.OpenGossipBundle(d, td.goodPeerID))
		assert.False(t, td.network.IsClosed(td.goodPeerID))
	})
}

func TestStreamMessage(t *testing.T) {
	t.Run("Message is nil => should close the connection", func(t *testing.T) {
		td := setup(t)

		assert.False(t, td.network.IsClosed(td.badPeerID))
		assert.Nil(t, td.firewall.OpenStreamBundle(bytes.NewReader(nil), td.badPeerID))
		assert.True(t, td.network.IsClosed(td.badPeerID))
	})

	t.Run("Message from: bad => should close the connection", func(t *testing.T) {
		td := setup(t)

		bdl := bundle.NewBundle(message.NewBlocksRequestMessage(td.RandInt(100), 1, 100))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
		d, _ := bdl.Encode()

		assert.False(t, td.network.IsClosed(td.badPeerID))
		assert.Nil(t, td.firewall.OpenStreamBundle(bytes.NewReader(d), td.badPeerID))
		assert.True(t, td.network.IsClosed(td.badPeerID))
	})

	t.Run("Ok => should NOT close the connection", func(t *testing.T) {
		td := setup(t)

		bdl := bundle.NewBundle(message.NewBlocksRequestMessage(td.RandInt(100), 1, 100))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
		d, _ := bdl.Encode()

		assert.False(t, td.network.IsClosed(td.goodPeerID))
		assert.NotNil(t, td.firewall.OpenStreamBundle(bytes.NewReader(d), td.goodPeerID))
		assert.False(t, td.network.IsClosed(td.goodPeerID))
	})
}

func TestDisabledFirewall(t *testing.T) {
	td := setup(t)

	bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), -1, td.RandValAddress()))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
	d, _ := bdl.Encode()

	td.firewall.config.Enabled = false
	assert.Nil(t, td.firewall.OpenGossipBundle(d, td.badPeerID))
	assert.False(t, td.network.IsClosed(td.badPeerID))
}

func TestUpdateLastReceived(t *testing.T) {
	td := setup(t)

	bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
	d, _ := bdl.Encode()
	now := time.Now().UnixNano()
	assert.NotNil(t, td.firewall.OpenGossipBundle(d, td.goodPeerID))

	peerGood := td.firewall.peerSet.GetPeer(td.goodPeerID)
	assert.GreaterOrEqual(t, peerGood.LastReceived.UnixNano(), now)
}

func TestNetworkFlags(t *testing.T) {
	td := setup(t)

	// TODO: add tests for Mainnet and Testnet flags
	bdl := bundle.NewBundle(message.NewQueryVotesMessage(td.RandHeight(), td.RandRound(), td.RandValAddress()))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
	assert.NoError(t, td.firewall.checkBundle(bdl))

	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	assert.Error(t, td.firewall.checkBundle(bdl))

	bdl.Flags = 0
	assert.Error(t, td.firewall.checkBundle(bdl))

	td.state.TestParams.BlockVersion = 0x3f // changing genesis hash
	bdl.Flags = 1
	assert.Error(t, td.firewall.checkBundle(bdl))
}
