package sync

import (
	"testing"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func (td *testData) connect(pid peer.ID, direction lp2pnetwork.Direction) {
	td.sync.processConnectEvent(&network.ConnectEvent{
		PeerID: pid, RemoteAddress: "/ip4/2.2.2.2/tcp/21888", Direction: direction,
	})
}

func (td *testData) identify(pid peer.ID) {
	td.sync.processProtocolsEvent(&network.ProtocolsEvents{PeerID: pid, Protocols: []string{"protocol-1"}})
}

func (td *testData) helloFrom(t *testing.T, pid peer.ID) *message.HelloMessage {
	t.Helper()

	msg := message.NewHelloMessage(pid, td.RandString(12), service.New(service.FullNode),
		td.RandHeight(), td.RandHash(), td.state.Genesis().Hash())
	require.NoError(t, msg.Sign([]*bls.ValidatorKey{td.RandValKey()}))

	return msg
}

// Two peers that dial each other at the same time have two connections. The handshake follows
// the first one, whatever the order in which the connections are reported.
func TestHandshakeWhenBothPeersDial(t *testing.T) {
	t.Run("our connection is reported first", func(t *testing.T) {
		td := setup(t, nil)
		td.state.EXPECT().UpdateValidatorProtocolVersion(gomock.Any(), gomock.Any()).AnyTimes()
		pid := td.RandPeerID()

		td.connect(pid, lp2pnetwork.DirOutbound)
		td.identify(pid)
		td.shouldPublishMessageWithThisType(t, message.TypeHello)

		// The connection of the other peer doesn't change the direction of the handshake.
		td.connect(pid, lp2pnetwork.DirInbound)
		assert.Equal(t, lp2pnetwork.DirOutbound, td.sync.peerSet.GetPeer(pid).Direction)

		td.receivingNewMessage(td.sync, td.helloFrom(t, pid), pid)
		td.shouldPublishHelloAck(t, message.ResponseCodeOK)
		assert.Equal(t, status.StatusKnown, td.sync.peerSet.GetPeerStatus(pid))
	})

	t.Run("each peer reports the connection of the other first", func(t *testing.T) {
		td := setup(t, nil)
		pid := td.RandPeerID()

		td.connect(pid, lp2pnetwork.DirInbound)
		td.identify(pid)
		td.shouldNotPublishAnyMessage(t)

		// The other peer waits for our Hello too: our own connection sends it.
		td.connect(pid, lp2pnetwork.DirOutbound)
		td.shouldPublishMessageWithThisType(t, message.TypeHello)
		peer := td.sync.peerSet.GetPeer(pid)
		assert.Equal(t, lp2pnetwork.DirOutbound, peer.Direction)
		assert.True(t, peer.OutboundHelloSent)
	})

	t.Run("the other peer sent its Hello before our connection is reported", func(t *testing.T) {
		td := setup(t, nil)
		td.state.EXPECT().UpdateValidatorProtocolVersion(gomock.Any(), gomock.Any()).AnyTimes()
		pid := td.RandPeerID()

		td.connect(pid, lp2pnetwork.DirInbound)
		td.identify(pid)
		td.receivingNewMessage(td.sync, td.helloFrom(t, pid), pid)
		td.shouldPublishMessageWithThisType(t, message.TypeHello)
		td.shouldPublishHelloAck(t, message.ResponseCodeOK)

		// The handshake is going on: our connection doesn't send another Hello.
		td.connect(pid, lp2pnetwork.DirOutbound)
		td.shouldNotPublishAnyMessage(t)
		assert.Equal(t, lp2pnetwork.DirInbound, td.sync.peerSet.GetPeer(pid).Direction)
	})
}

// A peer that disconnects and connects again starts a new handshake.
func TestHandshakeAfterReconnection(t *testing.T) {
	td := setup(t, nil)
	td.state.EXPECT().UpdateValidatorProtocolVersion(gomock.Any(), gomock.Any()).AnyTimes()
	pid := td.RandPeerID()

	td.connect(pid, lp2pnetwork.DirOutbound)
	td.identify(pid)
	td.shouldPublishMessageWithThisType(t, message.TypeHello)
	td.receivingNewMessage(td.sync, td.helloFrom(t, pid), pid)
	td.shouldPublishHelloAck(t, message.ResponseCodeOK)
	require.Equal(t, status.StatusKnown, td.sync.peerSet.GetPeerStatus(pid))

	td.sync.processDisconnectEvent(&network.DisconnectEvent{PeerID: pid})

	// This time the other peer dials: it sends the first Hello, and it is accepted.
	td.connect(pid, lp2pnetwork.DirInbound)
	td.identify(pid)
	td.receivingNewMessage(td.sync, td.helloFrom(t, pid), pid)
	td.shouldPublishMessageWithThisType(t, message.TypeHello)
	td.shouldPublishHelloAck(t, message.ResponseCodeOK)
}
