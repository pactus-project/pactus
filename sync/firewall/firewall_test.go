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
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
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

	peerGood := tFirewall.peerSet.MustGetPeer(tGoodPeerID)
	peerGood.UpdateStatus(peerset.StatusCodeGood)

	badGood := tFirewall.peerSet.MustGetPeer(tBadPeerID)
	badGood.UpdateStatus(peerset.StatusCodeBanned)
}

func TestInvalidMessagesCounter(t *testing.T) {
	setup(t)

	assert.Nil(t, tFirewall.OpenMessage(bytes.NewReader([]byte("bad")), tUnknownPeerID, tUnknownPeerID))
	assert.Nil(t, tFirewall.OpenMessage(bytes.NewReader(nil), tUnknownPeerID, tUnknownPeerID))

	msg := message.NewMessage(tUnknownPeerID, payload.NewQueryProposalPayload(-1, 1))
	d, _ := msg.Encode()
	assert.Nil(t, tFirewall.OpenMessage(bytes.NewReader(d), tUnknownPeerID, tUnknownPeerID))

	msg = message.NewMessage(tBadPeerID, payload.NewQueryProposalPayload(0, 1))
	d, _ = msg.Encode()
	assert.Nil(t, tFirewall.OpenMessage(bytes.NewReader(d), tUnknownPeerID, tUnknownPeerID))

	peer := tFirewall.peerSet.GetPeer(tUnknownPeerID)
	assert.Equal(t, peer.InvalidMessages(), 4)
}

func TestBanPeer(t *testing.T) {
	t.Run("Message source: unknown, from: bad => should close the connection", func(t *testing.T) {
		setup(t)

		msg := message.NewMessage(tUnknownPeerID, payload.NewQueryProposalPayload(100, 1))
		d, _ := msg.Encode()

		assert.False(t, tNetwork.IsClosed(tBadPeerID))
		assert.Nil(t, tFirewall.OpenMessage(bytes.NewReader(d), tUnknownPeerID, tBadPeerID))
		assert.True(t, tNetwork.IsClosed(tBadPeerID))
	})

	t.Run("Message source: bad, from: unknown => should close the connection", func(t *testing.T) {
		setup(t)

		msg := message.NewMessage(tBadPeerID, payload.NewQueryProposalPayload(100, 1))
		d, _ := msg.Encode()

		assert.False(t, tNetwork.IsClosed(tBadPeerID))
		assert.Nil(t, tFirewall.OpenMessage(bytes.NewReader(d), tBadPeerID, tUnknownPeerID))
		assert.True(t, tNetwork.IsClosed(tBadPeerID))
	})

	t.Run("Message initiator is not the same as source => should close the connection", func(t *testing.T) {
		setup(t)

		msg := message.NewMessage(tBadPeerID, payload.NewQueryProposalPayload(100, 1))
		d, _ := msg.Encode()

		assert.Nil(t, tFirewall.OpenMessage(bytes.NewReader(d), tUnknownPeerID, tUnknownPeerID))
		assert.True(t, tNetwork.IsClosed(tUnknownPeerID))
	})

	t.Run("Ok => should NOT close the connection", func(t *testing.T) {
		setup(t)

		msg := message.NewMessage(tGoodPeerID, payload.NewQueryProposalPayload(100, 1))
		d, _ := msg.Encode()

		assert.False(t, tNetwork.IsClosed(tGoodPeerID))
		assert.NotNil(t, tFirewall.OpenMessage(bytes.NewReader(d), tGoodPeerID, tGoodPeerID))
		assert.False(t, tNetwork.IsClosed(tGoodPeerID))
	})
}

func TestDisabledFirewal(t *testing.T) {
	setup(t)

	msg := message.NewMessage(tGoodPeerID, payload.NewQueryProposalPayload(-1, -1))
	d, _ := msg.Encode()

	tFirewall.config.Enabled = false
	assert.Nil(t, tFirewall.OpenMessage(bytes.NewReader(d), tBadPeerID, tBadPeerID))
	assert.False(t, tNetwork.IsClosed(tBadPeerID))
}
