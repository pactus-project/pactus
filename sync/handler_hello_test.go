package sync

import (
	"testing"
	"time"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func (td *testData) validHelloMessage(t *testing.T) *message.HelloMessage {
	t.Helper()

	valKey1 := td.RandValKey()
	valKey2 := td.RandValKey()
	msg := message.NewHelloMessage(td.RandPeerID(), td.RandString(12),
		service.New(service.FullNode),
		td.RandHeight(), td.RandHash(), td.state.Genesis().Hash())

	require.NoError(t, msg.Sign([]*bls.ValidatorKey{valKey1, valKey2}))

	return msg
}

func (td *testData) shouldPublishHelloAck(t *testing.T, code message.ResponseCode) {
	t.Helper()

	bdl := td.shouldPublishMessageWithThisType(t, message.TypeHelloAck)
	assert.Equal(t, code, bdl.Message.(*message.HelloAckMessage).ResponseCode)
}

func (td *testData) connectPeer(pid peer.ID, direction lp2pnetwork.Direction, outboundHelloSent bool) {
	td.sync.peerSet.UpdateAddress(pid, "some-address", direction)
	td.sync.peerSet.UpdateStatus(pid, status.StatusConnected)
	td.sync.peerSet.UpdateOutboundHelloSent(pid, outboundHelloSent)
}

func TestHandlerHelloParsingMessages(t *testing.T) {
	td := setup(t, nil)

	td.state.EXPECT().UpdateValidatorProtocolVersion(gomock.Any(), gomock.Any()).AnyTimes()

	t.Run("Receiving Hello message from a peer. Genesis hash is wrong.",
		func(t *testing.T) {
			msg := td.validHelloMessage(t)

			td.sync.peerSet.UpdateAddress(msg.PeerID, "some-address", lp2pnetwork.DirInbound)
			td.sync.peerSet.UpdateStatus(msg.PeerID, status.StatusConnected)

			msg.GenesisHash = td.RandHash()

			td.receivingNewMessage(td.sync, msg, msg.PeerID)
			td.checkPeerStatus(t, msg.PeerID, status.StatusBanned)
			td.shouldPublishHelloAck(t, message.ResponseCodeRejected)
		})

	t.Run("Receiving a Hello message from a peer. The time difference is greater than or equal to -10",
		func(t *testing.T) {
			msg := td.validHelloMessage(t)

			td.sync.peerSet.UpdateAddress(msg.PeerID, "some-address", lp2pnetwork.DirInbound)
			td.sync.peerSet.UpdateStatus(msg.PeerID, status.StatusConnected)

			msg.MyTimeUnixMilli = msg.MyTime().Add(-10 * time.Second).UnixMilli()

			td.receivingNewMessage(td.sync, msg, msg.PeerID)
			td.checkPeerStatus(t, msg.PeerID, status.StatusBanned)
			td.shouldPublishHelloAck(t, message.ResponseCodeRejected)
		})

	t.Run("Receiving Hello message from a peer. Difference is less or equal than 20 seconds.",
		func(t *testing.T) {
			msg := td.validHelloMessage(t)

			td.sync.peerSet.UpdateAddress(msg.PeerID, "some-address", lp2pnetwork.DirInbound)
			td.sync.peerSet.UpdateStatus(msg.PeerID, status.StatusConnected)

			msg.MyTimeUnixMilli = msg.MyTime().Add(20 * time.Second).UnixMilli()

			td.receivingNewMessage(td.sync, msg, msg.PeerID)
			td.checkPeerStatus(t, msg.PeerID, status.StatusBanned)
			td.shouldPublishHelloAck(t, message.ResponseCodeRejected)
		})

	t.Run("Non supporting version",
		func(t *testing.T) {
			msg := td.validHelloMessage(t)

			td.sync.peerSet.UpdateAddress(msg.PeerID, "some-address", lp2pnetwork.DirInbound)
			td.sync.peerSet.UpdateStatus(msg.PeerID, status.StatusConnected)

			nodeAgent := version.NodeAgent
			nodeAgent.Version = version.Version{
				Major: 1,
				Minor: 8,
				Patch: 0,
			}
			msg.Agent = nodeAgent.String()

			td.receivingNewMessage(td.sync, msg, msg.PeerID)
			td.checkPeerStatus(t, msg.PeerID, status.StatusBanned)
			td.shouldPublishHelloAck(t, message.ResponseCodeRejected)
		})

	t.Run("Invalid agent",
		func(t *testing.T) {
			msg := td.validHelloMessage(t)

			td.sync.peerSet.UpdateAddress(msg.PeerID, "some-address", lp2pnetwork.DirInbound)
			td.sync.peerSet.UpdateStatus(msg.PeerID, status.StatusConnected)

			msg.Agent = "invalid-agent"

			td.receivingNewMessage(td.sync, msg, msg.PeerID)
			td.checkPeerStatus(t, msg.PeerID, status.StatusBanned)
			td.shouldPublishHelloAck(t, message.ResponseCodeRejected)
		})

	t.Run("Outdated protocol version, network not upgraded",
		func(t *testing.T) {
			msg := td.validHelloMessage(t)

			td.sync.peerSet.UpdateAddress(msg.PeerID, "some-address", lp2pnetwork.DirInbound)
			td.sync.peerSet.UpdateStatus(msg.PeerID, status.StatusConnected)

			nodeAgent := version.NodeAgent
			nodeAgent.ProtocolVersion = protocol.ProtocolVersionLatest - 1
			msg.Agent = nodeAgent.String()
			td.state.StateParams.BlockVersion = protocol.ProtocolVersionLatest - 1

			td.receivingNewMessage(td.sync, msg, msg.PeerID)
			td.checkPeerStatus(t, msg.PeerID, status.StatusConnected)
			td.shouldPublishMessageWithThisType(t, message.TypeHello)
			td.shouldPublishHelloAck(t, message.ResponseCodeOK)
		})

	t.Run("Outdated protocol version, network upgraded",
		func(t *testing.T) {
			msg := td.validHelloMessage(t)

			td.sync.peerSet.UpdateAddress(msg.PeerID, "some-address", lp2pnetwork.DirInbound)
			td.sync.peerSet.UpdateStatus(msg.PeerID, status.StatusConnected)

			nodeAgent := version.NodeAgent
			nodeAgent.ProtocolVersion = protocol.ProtocolVersionLatest - 1
			msg.Agent = nodeAgent.String()
			td.state.StateParams.BlockVersion = protocol.ProtocolVersionLatest

			td.receivingNewMessage(td.sync, msg, msg.PeerID)
			td.checkPeerStatus(t, msg.PeerID, status.StatusBanned)
			td.shouldPublishHelloAck(t, message.ResponseCodeRejected)
		})

	t.Run("Receiving Hello message from a peer. It should be acknowledged and updates the peer info",
		func(t *testing.T) {
			msg := td.validHelloMessage(t)

			td.sync.peerSet.UpdateAddress(msg.PeerID, "some-address", lp2pnetwork.DirInbound)
			td.sync.peerSet.UpdateStatus(msg.PeerID, status.StatusConnected)

			td.receivingNewMessage(td.sync, msg, msg.PeerID)
			td.shouldPublishMessageWithThisType(t, message.TypeHello)
			td.shouldPublishHelloAck(t, message.ResponseCodeOK)

			// Check if the peer info is updated
			peer := td.sync.peerSet.GetPeer(msg.PeerID)

			assert.Equal(t, status.StatusConnected, peer.Status)
			assert.Equal(t, version.NodeAgent.String(), peer.Agent)
			assert.Equal(t, msg.Moniker, peer.Moniker)
			assert.Equal(t, msg.PublicKeys, peer.ConsensusKeys)
			assert.Equal(t, msg.PeerID, peer.PeerID)
			assert.Equal(t, msg.Height, peer.Height)
			assert.True(t, peer.IsFullNode())
		})
}

func TestHandlerHelloHandshaking(t *testing.T) {
	td := setup(t, nil)

	td.state.EXPECT().UpdateValidatorProtocolVersion(gomock.Any(), protocol.ProtocolVersionLatest).AnyTimes()

	t.Run("Unknown Direction", func(t *testing.T) {
		msg := td.validHelloMessage(t)

		td.connectPeer(msg.PeerID, lp2pnetwork.DirUnknown, false)
		td.receivingNewMessage(td.sync, msg, msg.PeerID)
		td.shouldNotPublishAnyMessage(t)
	})

	t.Run("Inbound Direction, Outbound Hello Not Sent", func(t *testing.T) {
		msg := td.validHelloMessage(t)

		td.connectPeer(msg.PeerID, lp2pnetwork.DirInbound, false)
		td.receivingNewMessage(td.sync, msg, msg.PeerID)
		td.shouldPublishMessageWithThisType(t, message.TypeHello)
		td.shouldPublishMessageWithThisType(t, message.TypeHelloAck)
		td.checkPeerStatus(t, msg.PeerID, status.StatusConnected)
	})

	t.Run("Outbound Direction, Outbound Hello Not Sent", func(t *testing.T) {
		msg := td.validHelloMessage(t)

		td.connectPeer(msg.PeerID, lp2pnetwork.DirOutbound, false)
		td.receivingNewMessage(td.sync, msg, msg.PeerID)
		td.shouldNotPublishAnyMessage(t)
		td.checkPeerStatus(t, msg.PeerID, status.StatusConnected)
	})

	t.Run("Inbound Direction, Outbound Hello Sent", func(t *testing.T) {
		msg := td.validHelloMessage(t)

		td.connectPeer(msg.PeerID, lp2pnetwork.DirInbound, true)
		td.receivingNewMessage(td.sync, msg, msg.PeerID)
		td.shouldNotPublishAnyMessage(t)
		td.checkPeerStatus(t, msg.PeerID, status.StatusConnected)
	})

	t.Run("Outbound Direction, Outbound Hello Sent", func(t *testing.T) {
		msg := td.validHelloMessage(t)

		td.connectPeer(msg.PeerID, lp2pnetwork.DirOutbound, true)
		td.receivingNewMessage(td.sync, msg, msg.PeerID)
		td.shouldPublishMessageWithThisType(t, message.TypeHelloAck)
		td.checkPeerStatus(t, msg.PeerID, status.StatusKnown)
	})
}

func TestHandlerHelloIsBoundToTheConnection(t *testing.T) {
	td := setup(t, nil)

	td.state.EXPECT().UpdateValidatorProtocolVersion(gomock.Any(), gomock.Any()).AnyTimes()

	t.Run("Hello replayed by another peer", func(t *testing.T) {
		msg := td.validHelloMessage(t)
		replayer := td.RandPeerID()

		td.connectPeer(replayer, lp2pnetwork.DirInbound, false)
		td.receivingNewMessage(td.sync, msg, replayer)

		td.shouldPublishHelloAck(t, message.ResponseCodeRejected)
		td.checkPeerStatus(t, replayer, status.StatusBanned)
		assert.Empty(t, td.sync.peerSet.GetPeer(replayer).ConsensusKeys)
		assert.Empty(t, td.sync.peerSet.GetPeer(replayer).Moniker)
	})

	t.Run("Rejected Hello doesn't record the peer info", func(t *testing.T) {
		msg := td.validHelloMessage(t)
		msg.GenesisHash = td.RandHash()

		td.connectPeer(msg.PeerID, lp2pnetwork.DirInbound, false)
		td.receivingNewMessage(td.sync, msg, msg.PeerID)

		td.shouldPublishHelloAck(t, message.ResponseCodeRejected)
		peer := td.sync.peerSet.GetPeer(msg.PeerID)
		assert.Empty(t, peer.ConsensusKeys)
		assert.Empty(t, peer.Moniker)
		assert.Zero(t, peer.Height)
	})

	t.Run("Unexpected Hello doesn't overwrite the peer info", func(t *testing.T) {
		first := td.validHelloMessage(t)

		td.connectPeer(first.PeerID, lp2pnetwork.DirInbound, false)
		td.receivingNewMessage(td.sync, first, first.PeerID)
		td.shouldPublishMessageWithThisType(t, message.TypeHello)
		td.shouldPublishHelloAck(t, message.ResponseCodeOK)

		second := td.validHelloMessage(t)
		second.PeerID = first.PeerID
		second.Moniker = "other"
		td.receivingNewMessage(td.sync, second, first.PeerID)

		assert.Equal(t, first.Moniker, td.sync.peerSet.GetPeer(first.PeerID).Moniker)
		assert.Equal(t, first.PublicKeys, td.sync.peerSet.GetPeer(first.PeerID).ConsensusKeys)
	})
}
