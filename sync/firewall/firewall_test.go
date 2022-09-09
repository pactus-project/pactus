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
	"github.com/stretchr/testify/assert"
)

var tFirewall *Firewall
var tBadPeerID peer.ID
var tGoodPeerID peer.ID
var tUnknownPeerID peer.ID
var tNetwork *network.MockNetwork
var tState *state.MockState

func setup(t *testing.T) {
	logger := logger.NewLogger("firewal", nil)
	peerSet := peerset.NewPeerSet(3 * time.Second)
	tState = state.MockingState()
	tNetwork = network.MockingNetwork(network.TestRandomPeerID())
	conf := DefaultConfig()
	conf.Enabled = true
	tFirewall = NewFirewall(conf, tNetwork, peerSet, tState, logger)
	assert.NotNil(t, tFirewall)
	tBadPeerID = network.TestRandomPeerID()
	tGoodPeerID = network.TestRandomPeerID()
	tUnknownPeerID = network.TestRandomPeerID()

	tNetwork.AddAnotherNetwork(network.MockingNetwork(tGoodPeerID))
	tNetwork.AddAnotherNetwork(network.MockingNetwork(tUnknownPeerID))
	tNetwork.AddAnotherNetwork(network.MockingNetwork(tBadPeerID))

	tFirewall.peerSet.UpdateStatus(tGoodPeerID, peerset.StatusCodeKnown)
	tFirewall.peerSet.UpdateStatus(tBadPeerID, peerset.StatusCodeBanned)
}

func TestInvalidBundlesCounter(t *testing.T) {
	setup(t)

	assert.Nil(t, tFirewall.OpenGossipBundle([]byte("bad"), tUnknownPeerID, tUnknownPeerID))
	assert.Nil(t, tFirewall.OpenGossipBundle(nil, tUnknownPeerID, tUnknownPeerID))

	bdl := bundle.NewBundle(tUnknownPeerID, message.NewQueryProposalMessage(0, -1))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	d, _ := bdl.Encode()
	assert.Nil(t, tFirewall.OpenGossipBundle(d, tUnknownPeerID, tUnknownPeerID))

	bdl = bundle.NewBundle(tBadPeerID, message.NewQueryProposalMessage(0, 1))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	d, _ = bdl.Encode()
	assert.Nil(t, tFirewall.OpenGossipBundle(d, tUnknownPeerID, tUnknownPeerID))

	peer := tFirewall.peerSet.GetPeer(tUnknownPeerID)
	assert.Equal(t, peer.InvalidBundles, 4)
}

func TestGossipMesage(t *testing.T) {
	t.Run("Message source: unknown, from: bad => should close the connection", func(t *testing.T) {
		setup(t)

		bdl := bundle.NewBundle(tUnknownPeerID, message.NewQueryProposalMessage(100, 1))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
		d, _ := bdl.Encode()

		assert.False(t, tNetwork.IsClosed(tBadPeerID))
		assert.Nil(t, tFirewall.OpenGossipBundle(d, tUnknownPeerID, tBadPeerID))
		assert.True(t, tNetwork.IsClosed(tBadPeerID))
	})

	t.Run("Message source: bad, from: unknown => should close the connection", func(t *testing.T) {
		setup(t)

		bdl := bundle.NewBundle(tBadPeerID, message.NewQueryProposalMessage(100, 1))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
		d, _ := bdl.Encode()

		assert.False(t, tNetwork.IsClosed(tBadPeerID))
		assert.Nil(t, tFirewall.OpenGossipBundle(d, tBadPeerID, tUnknownPeerID))
		assert.True(t, tNetwork.IsClosed(tBadPeerID))
	})

	t.Run("Message initiator is not the same as source => should close the connection", func(t *testing.T) {
		setup(t)

		bdl := bundle.NewBundle(tBadPeerID, message.NewQueryProposalMessage(100, 1))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
		d, _ := bdl.Encode()

		assert.Nil(t, tFirewall.OpenGossipBundle(d, tUnknownPeerID, tUnknownPeerID))
		assert.True(t, tNetwork.IsClosed(tUnknownPeerID))
	})

	t.Run("Ok => should NOT close the connection", func(t *testing.T) {
		setup(t)

		bdl := bundle.NewBundle(tGoodPeerID, message.NewQueryProposalMessage(100, 1))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
		d, _ := bdl.Encode()

		assert.False(t, tNetwork.IsClosed(tGoodPeerID))
		assert.NotNil(t, tFirewall.OpenGossipBundle(d, tGoodPeerID, tGoodPeerID))
		assert.False(t, tNetwork.IsClosed(tGoodPeerID))
	})
}

func TestStreamMesage(t *testing.T) {
	t.Run("Message source: bad => should close the connection", func(t *testing.T) {
		setup(t)

		bdl := bundle.NewBundle(tBadPeerID, message.NewBlocksRequestMessage(int(util.RandInt32(0)), 1, 100))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
		d, _ := bdl.Encode()

		assert.False(t, tNetwork.IsClosed(tBadPeerID))
		assert.Nil(t, tFirewall.OpenStreamBundle(bytes.NewReader(d), tBadPeerID))
		assert.True(t, tNetwork.IsClosed(tBadPeerID))
	})

	t.Run("Ok => should NOT close the connection", func(t *testing.T) {
		setup(t)

		bdl := bundle.NewBundle(tGoodPeerID, message.NewBlocksRequestMessage(int(util.RandInt32(0)), 1, 100))
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
		d, _ := bdl.Encode()

		assert.False(t, tNetwork.IsClosed(tGoodPeerID))
		assert.NotNil(t, tFirewall.OpenStreamBundle(bytes.NewReader(d), tGoodPeerID))
		assert.False(t, tNetwork.IsClosed(tGoodPeerID))
	})
}

func TestDisabledFirewal(t *testing.T) {
	setup(t)

	bdl := bundle.NewBundle(tGoodPeerID, message.NewQueryProposalMessage(0, -1))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	d, _ := bdl.Encode()

	tFirewall.config.Enabled = false
	assert.Nil(t, tFirewall.OpenGossipBundle(d, tBadPeerID, tBadPeerID))
	assert.False(t, tNetwork.IsClosed(tBadPeerID))
}

func TestUpdateLastSeen(t *testing.T) {
	setup(t)

	bdl := bundle.NewBundle(tGoodPeerID, message.NewQueryProposalMessage(100, 1))
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	d, _ := bdl.Encode()
	now := time.Now().UnixNano()
	assert.Nil(t, tFirewall.OpenGossipBundle(d, tUnknownPeerID, tGoodPeerID))

	assert.GreaterOrEqual(t, tFirewall.peerSet.GetPeer(tUnknownPeerID).LastSeen.UnixNano(), now)
	assert.GreaterOrEqual(t, tFirewall.peerSet.GetPeer(tGoodPeerID).LastSeen.UnixNano(), now)
}

func TestNetworkFlags(t *testing.T) {
	setup(t)

	bdl := bundle.NewBundle(tGoodPeerID, message.NewQueryProposalMessage(100, 1))
	bdl.Flags = util.UnsetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	assert.Error(t, tFirewall.checkBundle(bdl, tGoodPeerID))

	tState.TestParams.BlockVersion = 0x3f
	assert.Error(t, tFirewall.checkBundle(bdl, tGoodPeerID))
}
