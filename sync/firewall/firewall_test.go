package firewall

import (
	"bytes"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/bundle"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

var tFirewall *Firewall
var tBadPeerID peer.ID
var tGoodPeerID peer.ID
var tUnknownPeerID peer.ID
var tNetwork *network.MockNetwork

func setup(t *testing.T) {
	logger.InitLogger(logger.TestConfig())
	logger := logger.NewLogger("firewal", nil)
	peerSet := peerset.NewPeerSet(3 * time.Second)
	committee, _ := committee.GenerateTestCommittee()
	state := state.MockingState(committee)
	tNetwork = network.MockingNetwork(util.RandomPeerID())
	conf := TestConfig()
	conf.Enabled = true
	tFirewall = NewFirewall(conf, tNetwork, peerSet, state, logger)
	assert.NotNil(t, tFirewall)
	tBadPeerID = util.RandomPeerID()
	tGoodPeerID = util.RandomPeerID()
	tUnknownPeerID = util.RandomPeerID()

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

	msg := bundle.NewBundle(tUnknownPeerID, message.NewQueryProposalMessage(-1, 1))
	d, _ := msg.Encode()
	assert.Nil(t, tFirewall.OpenGossipBundle(d, tUnknownPeerID, tUnknownPeerID))

	msg = bundle.NewBundle(tBadPeerID, message.NewQueryProposalMessage(0, 1))
	d, _ = msg.Encode()
	assert.Nil(t, tFirewall.OpenGossipBundle(d, tUnknownPeerID, tUnknownPeerID))

	peer := tFirewall.peerSet.GetPeer(tUnknownPeerID)
	assert.Equal(t, peer.InvalidBundles, 4)
}

func TestGossipMesage(t *testing.T) {
	t.Run("Message source: unknown, from: bad => should close the connection", func(t *testing.T) {
		setup(t)

		msg := bundle.NewBundle(tUnknownPeerID, message.NewQueryProposalMessage(100, 1))
		d, _ := msg.Encode()

		assert.False(t, tNetwork.IsClosed(tBadPeerID))
		assert.Nil(t, tFirewall.OpenGossipBundle(d, tUnknownPeerID, tBadPeerID))
		assert.True(t, tNetwork.IsClosed(tBadPeerID))
	})

	t.Run("Message source: bad, from: unknown => should close the connection", func(t *testing.T) {
		setup(t)

		msg := bundle.NewBundle(tBadPeerID, message.NewQueryProposalMessage(100, 1))
		d, _ := msg.Encode()

		assert.False(t, tNetwork.IsClosed(tBadPeerID))
		assert.Nil(t, tFirewall.OpenGossipBundle(d, tBadPeerID, tUnknownPeerID))
		assert.True(t, tNetwork.IsClosed(tBadPeerID))
	})

	t.Run("Message initiator is not the same as source => should close the connection", func(t *testing.T) {
		setup(t)

		msg := bundle.NewBundle(tBadPeerID, message.NewQueryProposalMessage(100, 1))
		d, _ := msg.Encode()

		assert.Nil(t, tFirewall.OpenGossipBundle(d, tUnknownPeerID, tUnknownPeerID))
		assert.True(t, tNetwork.IsClosed(tUnknownPeerID))
	})

	t.Run("Ok => should NOT close the connection", func(t *testing.T) {
		setup(t)

		msg := bundle.NewBundle(tGoodPeerID, message.NewQueryProposalMessage(100, 1))
		d, _ := msg.Encode()

		assert.False(t, tNetwork.IsClosed(tGoodPeerID))
		assert.NotNil(t, tFirewall.OpenGossipBundle(d, tGoodPeerID, tGoodPeerID))
		assert.False(t, tNetwork.IsClosed(tGoodPeerID))
	})
}

func TestStreamMesage(t *testing.T) {

	t.Run("Message source: bad => should close the connection", func(t *testing.T) {
		setup(t)

		msg := bundle.NewBundle(tBadPeerID, message.NewBlocksRequestMessage(util.RandInt(0), 1, 100))
		d, _ := msg.Encode()

		assert.False(t, tNetwork.IsClosed(tBadPeerID))
		assert.Nil(t, tFirewall.OpenStreamBundle(bytes.NewReader(d), tBadPeerID))
		assert.True(t, tNetwork.IsClosed(tBadPeerID))
	})

	t.Run("Ok => should NOT close the connection", func(t *testing.T) {
		setup(t)

		msg := bundle.NewBundle(tGoodPeerID, message.NewBlocksRequestMessage(util.RandInt(0), 1, 100))
		d, _ := msg.Encode()

		assert.False(t, tNetwork.IsClosed(tGoodPeerID))
		assert.NotNil(t, tFirewall.OpenStreamBundle(bytes.NewReader(d), tGoodPeerID))
		assert.False(t, tNetwork.IsClosed(tGoodPeerID))
	})
}

func TestDisabledFirewal(t *testing.T) {
	setup(t)

	msg := bundle.NewBundle(tGoodPeerID, message.NewQueryProposalMessage(-1, -1))
	d, _ := msg.Encode()

	tFirewall.config.Enabled = false
	assert.Nil(t, tFirewall.OpenGossipBundle(d, tBadPeerID, tBadPeerID))
	assert.False(t, tNetwork.IsClosed(tBadPeerID))
}

func TestUpdateLastSeen(t *testing.T) {
	setup(t)

	msg := bundle.NewBundle(tGoodPeerID, message.NewQueryProposalMessage(100, 1))
	d, _ := msg.Encode()
	now := time.Now().UnixNano()
	assert.Nil(t, tFirewall.OpenGossipBundle(d, tUnknownPeerID, tGoodPeerID))

	assert.GreaterOrEqual(t, tFirewall.peerSet.GetPeer(tUnknownPeerID).LastSeen.UnixNano(), now)
	assert.GreaterOrEqual(t, tFirewall.peerSet.GetPeer(tGoodPeerID).LastSeen.UnixNano(), now)
}
