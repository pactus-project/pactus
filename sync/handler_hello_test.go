package sync

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/peerset/service"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
)

func TestParsingHelloMessages(t *testing.T) {
	td := setup(t, nil)

	td.state.CommitTestBlocks(21)

	t.Run("Receiving Hello message from a peer. Peer ID is not same as initiator.",
		func(t *testing.T) {
			valKey := td.RandValKey()
			pid := td.RandPeerID()
			initiator := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "bad-genesis", 0, 0,
				td.state.LastBlockHash(), td.state.Genesis().Hash())
			msg.Sign([]*bls.ValidatorKey{valKey})

			assert.NoError(t, td.receivingNewMessage(td.sync, msg, initiator))
			assert.Equal(t, td.sync.peerSet.GetPeer(initiator).Status, peerset.StatusCodeBanned)
			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeHelloAck)
			assert.Equal(t, bdl.Message.(*message.HelloAckMessage).ResponseCode, message.ResponseCodeRejected)
		})

	t.Run("Receiving Hello message from a peer. Genesis hash is wrong.",
		func(t *testing.T) {
			invGenHash := td.RandHash()
			valKey := td.RandValKey()
			pid := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "bad-genesis", 0, 0,
				td.state.LastBlockHash(), invGenHash)
			msg.Sign([]*bls.ValidatorKey{valKey})

			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))
			td.checkPeerStatus(t, pid, peerset.StatusCodeBanned)
			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeHelloAck)
			assert.Equal(t, bdl.Message.(*message.HelloAckMessage).ResponseCode, message.ResponseCodeRejected)
		})

	t.Run("Receiving Hello message from a peer. Difference is greater or equal than -10 seconds.",
		func(t *testing.T) {
			valKey := td.RandValKey()
			height := td.RandUint32NonZero(td.state.LastBlockHeight())
			pid := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "kitty", height, service.New(service.Network),
				td.state.LastBlockHash(), td.state.Genesis().Hash())
			msg.Sign([]*bls.ValidatorKey{valKey})

			msg.MyTimeUnixMilli = msg.MyTime().Add(-10 * time.Second).UnixMilli()
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))
			td.checkPeerStatus(t, pid, peerset.StatusCodeBanned)
			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeHelloAck)
			assert.Equal(t, bdl.Message.(*message.HelloAckMessage).ResponseCode, message.ResponseCodeRejected)
		})

	t.Run("Receiving Hello message from a peer. Difference is less or equal than 20 seconds.",
		func(t *testing.T) {
			valKey := td.RandValKey()
			height := td.RandUint32NonZero(td.state.LastBlockHeight())
			pid := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "kitty", height, service.New(service.Network),
				td.state.LastBlockHash(), td.state.Genesis().Hash())
			msg.Sign([]*bls.ValidatorKey{valKey})

			msg.MyTimeUnixMilli = msg.MyTime().Add(20 * time.Second).UnixMilli()
			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))
			td.checkPeerStatus(t, pid, peerset.StatusCodeBanned)
			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeHelloAck)
			assert.Equal(t, bdl.Message.(*message.HelloAckMessage).ResponseCode, message.ResponseCodeRejected)
		})

	t.Run("Receiving Hello message from a peer. It should be acknowledged and updates the peer info",
		func(t *testing.T) {
			valKey := td.RandValKey()
			height := td.RandUint32NonZero(td.state.LastBlockHeight())
			pid := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "kitty", height, service.New(service.Network),
				td.state.LastBlockHash(), td.state.Genesis().Hash())
			msg.Sign([]*bls.ValidatorKey{valKey})

			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeHelloAck)
			assert.Equal(t, bdl.Message.(*message.HelloAckMessage).ResponseCode, message.ResponseCodeOK)

			// Check if the peer info is updated
			p := td.sync.peerSet.GetPeer(pid)

			pub := valKey.PublicKey()
			assert.Equal(t, p.Status, peerset.StatusCodeKnown)
			assert.Equal(t, p.Agent, version.Agent())
			assert.Equal(t, p.Moniker, "kitty")
			assert.Contains(t, p.ConsensusKeys, pub)
			assert.Equal(t, p.PeerID, pid)
			assert.Equal(t, p.Height, height)
			assert.True(t, p.HasNetworkService())
		})
}

func TestSendingHelloMessage(t *testing.T) {
	td := setup(t, nil)

	to := td.RandPeerID()
	assert.NoError(t, td.sync.sayHello(to))

	bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeHello)
	assert.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagHandshaking))
	assert.True(t, util.IsFlagSet(bdl.Message.(*message.HelloMessage).Services, service.New(service.Network)))
}
