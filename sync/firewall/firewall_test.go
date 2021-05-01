package firewall

import (
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
	tFirewall = NewFirewall(true, tNetwork, peerSet, state, logger)
	tBadPeerID = util.RandomPeerID()
	tGoodPeerID = util.RandomPeerID()
	tUnknownPeerID = util.RandomPeerID()

	peerGood := tFirewall.peerSet.MustGetPeer(tGoodPeerID)
	peerGood.UpdateStatus(peerset.StatusCodeOK)

	badGood := tFirewall.peerSet.MustGetPeer(tBadPeerID)
	badGood.UpdateStatus(peerset.StatusCodeBanned)
}

func TestInvalidMessagesCounter(t *testing.T) {
	setup(t)

	assert.Nil(t, tFirewall.OpenMessage([]byte("bad"), tUnknownPeerID))
	assert.Nil(t, tFirewall.OpenMessage(nil, tUnknownPeerID))

	msg := message.NewMessage(tUnknownPeerID, payload.NewQueryProposalPayload(-1, 1))
	d, _ := msg.Encode()
	assert.Nil(t, tFirewall.OpenMessage(d, tUnknownPeerID))

	peer := tFirewall.peerSet.GetPeer(tUnknownPeerID)
	assert.Equal(t, peer.InvalidMessages(), 3)
}

func TestBanPeer(t *testing.T) {
	t.Run("Message from bad peer, initiated from bad peer => should close the connection", func(t *testing.T) {
		setup(t)

		msg := message.NewMessage(tBadPeerID, payload.NewQueryProposalPayload(100, 1))
		d, _ := msg.Encode()

		assert.Nil(t, tFirewall.OpenMessage(d, tBadPeerID))
		assert.True(t, tNetwork.Closed)
	})

	t.Run("Message from bad peer, initiated from good peer => should close the connection", func(t *testing.T) {
		setup(t)

		msg := message.NewMessage(tGoodPeerID, payload.NewQueryProposalPayload(100, 1))
		d, _ := msg.Encode()

		assert.Nil(t, tFirewall.OpenMessage(d, tBadPeerID))
		assert.True(t, tNetwork.Closed)
	})

	t.Run("Message from unknown peer, initiated from bad peer => should NOT close the connection", func(t *testing.T) {
		setup(t)

		msg := message.NewMessage(tBadPeerID, payload.NewQueryProposalPayload(100, 1))
		d, _ := msg.Encode()

		assert.Nil(t, tFirewall.OpenMessage(d, tUnknownPeerID))
		assert.False(t, tNetwork.Closed)
	})

	t.Run("Receive many bad messages from a good peer => should close the connection", func(t *testing.T) {
		setup(t)

		peer := tFirewall.peerSet.GetPeer(tGoodPeerID)
		peer.UpdateInvalidMessage(101)
		peer.UpdateReceivedMessage(1001)

		msg := message.NewMessage(tBadPeerID, payload.NewQueryProposalPayload(100, 1))
		d, _ := msg.Encode()

		assert.Nil(t, tFirewall.OpenMessage(d, tGoodPeerID))
		assert.True(t, tNetwork.Closed)
	})
}

func TestDropMessage(t *testing.T) {
	t.Run("Message initiated from bad peer => should drop them", func(t *testing.T) {
		setup(t)

		msg := message.NewMessage(tBadPeerID, payload.NewQueryProposalPayload(1, 0))
		d, _ := msg.Encode()

		assert.Nil(t, tFirewall.OpenMessage(d, tGoodPeerID))
	})

	t.Run("Message initiated from unknown peer", func(t *testing.T) {
		setup(t)

		msg := message.NewMessage(tUnknownPeerID, payload.NewQueryProposalPayload(1, 0))
		d, _ := msg.Encode()

		assert.NotNil(t, tFirewall.OpenMessage(d, tGoodPeerID))
	})

	t.Run("Message initiated from good peer", func(t *testing.T) {
		msg := message.NewMessage(tGoodPeerID, payload.NewQueryProposalPayload(1, 0))
		d, _ := msg.Encode()

		assert.NotNil(t, tFirewall.OpenMessage(d, tUnknownPeerID))
	})
}

func TestDisabledFirewal(t *testing.T) {
	setup(t)

	msg := message.NewMessage(tBadPeerID, payload.NewQueryProposalPayload(1, 0))
	d, _ := msg.Encode()

	tFirewall.Enabled = false
	assert.NotNil(t, tFirewall.OpenMessage(d, tBadPeerID))
}
