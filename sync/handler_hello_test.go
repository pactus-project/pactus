package sync

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
)

func TestParsingHelloMessages(t *testing.T) {
	td := setup(t, nil)

	td.state.CommitTestBlocks(21)

	t.Run("Receiving Hello message from an unknown peer.",
		func(t *testing.T) {
			valKey := td.RandValKey()
			pid := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "unknown-peer", service.New(service.FullNode),
				td.RandHeight(), td.RandHash(), td.state.Genesis().Hash())
			msg.Sign([]*bls.ValidatorKey{valKey})

			from := td.RandPeerID()
			td.receivingNewMessage(td.sync, msg, from)
			bdl := td.shouldPublishMessageWithThisType(t, message.TypeHelloAck)
			assert.Equal(t, message.ResponseCodeRejected, bdl.Message.(*message.HelloAckMessage).ResponseCode)
		})

	t.Run("Receiving Hello message from a peer. Genesis hash is wrong.",
		func(t *testing.T) {
			invGenHash := td.RandHash()
			valKey := td.RandValKey()
			pid := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "bad-genesis", service.New(service.FullNode),
				td.RandHeight(), td.RandHash(), invGenHash)
			msg.Sign([]*bls.ValidatorKey{valKey})

			td.receivingNewMessage(td.sync, msg, pid)
			td.checkPeerStatus(t, pid, status.StatusBanned)
			bdl := td.shouldPublishMessageWithThisType(t, message.TypeHelloAck)
			assert.Equal(t, message.ResponseCodeRejected, bdl.Message.(*message.HelloAckMessage).ResponseCode)
		})

	t.Run("Receiving a Hello message from a peer. The time difference is greater than or equal to -10",
		func(t *testing.T) {
			valKey := td.RandValKey()
			pid := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "kitty", service.New(service.FullNode),
				td.RandHeight(), td.RandHash(), td.state.Genesis().Hash())
			msg.Sign([]*bls.ValidatorKey{valKey})

			msg.MyTimeUnixMilli = msg.MyTime().Add(-10 * time.Second).UnixMilli()
			td.receivingNewMessage(td.sync, msg, pid)
			td.checkPeerStatus(t, pid, status.StatusBanned)
			bdl := td.shouldPublishMessageWithThisType(t, message.TypeHelloAck)
			assert.Equal(t, message.ResponseCodeRejected, bdl.Message.(*message.HelloAckMessage).ResponseCode)
		})

	t.Run("Receiving Hello message from a peer. Difference is less or equal than 20 seconds.",
		func(t *testing.T) {
			valKey := td.RandValKey()
			pid := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "kitty", service.New(service.FullNode),
				td.RandHeight(), td.RandHash(), td.state.Genesis().Hash())
			msg.Sign([]*bls.ValidatorKey{valKey})

			msg.MyTimeUnixMilli = msg.MyTime().Add(20 * time.Second).UnixMilli()
			td.receivingNewMessage(td.sync, msg, pid)
			td.checkPeerStatus(t, pid, status.StatusBanned)
			bdl := td.shouldPublishMessageWithThisType(t, message.TypeHelloAck)
			assert.Equal(t, message.ResponseCodeRejected, bdl.Message.(*message.HelloAckMessage).ResponseCode)
		})

	t.Run("Non supporting version.",
		func(t *testing.T) {
			valKey := td.RandValKey()
			pid := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "kitty", service.New(service.FullNode),
				td.RandHeight(), td.RandHash(), td.state.Genesis().Hash())
			nodeAgent := version.NodeAgent
			nodeAgent.Version = version.Version{
				Major: 1,
				Minor: 0,
				Patch: 2,
			}
			msg.Agent = nodeAgent.String()
			msg.Sign([]*bls.ValidatorKey{valKey})

			td.receivingNewMessage(td.sync, msg, pid)
			td.checkPeerStatus(t, pid, status.StatusBanned)
			bdl := td.shouldPublishMessageWithThisType(t, message.TypeHelloAck)
			assert.Equal(t, message.ResponseCodeRejected, bdl.Message.(*message.HelloAckMessage).ResponseCode)
		})

	t.Run("Invalid agent.",
		func(t *testing.T) {
			valKey := td.RandValKey()
			pid := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "kitty", service.New(service.FullNode),
				td.RandHeight(), td.RandHash(), td.state.Genesis().Hash())
			msg.Agent = "invalid-agent"
			msg.Sign([]*bls.ValidatorKey{valKey})

			td.receivingNewMessage(td.sync, msg, pid)
			td.checkPeerStatus(t, pid, status.StatusBanned)
			bdl := td.shouldPublishMessageWithThisType(t, message.TypeHelloAck)
			assert.Equal(t, message.ResponseCodeRejected, bdl.Message.(*message.HelloAckMessage).ResponseCode)
		})

	t.Run("Receiving Hello message from a peer. It should be acknowledged and updates the peer info",
		func(t *testing.T) {
			valKey := td.RandValKey()
			pid := td.RandPeerID()
			peerHeight := td.RandHeight()
			msg := message.NewHelloMessage(pid, "kitty", service.New(service.FullNode),
				peerHeight, td.RandHash(), td.state.Genesis().Hash())
			msg.Sign([]*bls.ValidatorKey{valKey})

			td.receivingNewMessage(td.sync, msg, pid)

			bdl := td.shouldPublishMessageWithThisType(t, message.TypeHelloAck)
			assert.Equal(t, message.ResponseCodeOK, bdl.Message.(*message.HelloAckMessage).ResponseCode)

			// Check if the peer info is updated
			peer := td.sync.peerSet.GetPeer(pid)

			pub := valKey.PublicKey()
			assert.Equal(t, status.StatusConnected, peer.Status)
			assert.Equal(t, version.NodeAgent.String(), peer.Agent)
			assert.Equal(t, "kitty", peer.Moniker)
			assert.Contains(t, peer.ConsensusKeys, pub)
			assert.Equal(t, pid, peer.PeerID)
			assert.Equal(t, peerHeight, peer.Height)
			assert.True(t, peer.IsFullNode())
		})
}
