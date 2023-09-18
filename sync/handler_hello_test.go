package sync

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/services"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
)

func TestParsingHelloMessages(t *testing.T) {
	td := setup(t, nil)

	td.state.CommitTestBlocks(21)

	t.Run("Receiving Hello message from a peer. Peer ID is not same as initiator.",
		func(t *testing.T) {
			signer := td.RandSigner()
			pid := td.RandPeerID()
			initiator := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "bad-genesis", 0, 0,
				td.state.LastBlockHash(), td.state.Genesis().Hash())
			msg.Sign(signer)

			assert.NoError(t, td.receivingNewMessage(td.sync, msg, initiator))
			assert.Equal(t, td.sync.peerSet.GetPeer(initiator).Status, peerset.StatusCodeBanned)
			bundle := td.shouldPublishMessageWithThisType(t, td.network, message.TypeHelloAck)
			assert.Equal(t, bundle.Message.(*message.HelloAckMessage).ResponseCode, message.ResponseCodeRejected)
		})

	t.Run("Receiving Hello message from a peer. Genesis hash is wrong.",
		func(t *testing.T) {
			invGenHash := td.RandHash()
			signer := td.RandSigner()
			pid := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "bad-genesis", 0, 0,
				td.state.LastBlockHash(), invGenHash)
			msg.Sign(signer)

			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))
			td.checkPeerStatus(t, pid, peerset.StatusCodeBanned)
			bundle := td.shouldPublishMessageWithThisType(t, td.network, message.TypeHelloAck)
			assert.Equal(t, bundle.Message.(*message.HelloAckMessage).ResponseCode, message.ResponseCodeRejected)
		})

	t.Run("Receiving Hello message from a peer. It should be acknowledged and updates the peer info",
		func(t *testing.T) {
			signer := td.RandSigner()
			height := td.RandUint32NonZero(td.state.LastBlockHeight())
			pid := td.RandPeerID()
			msg := message.NewHelloMessage(pid, "kitty", height, services.New(services.Network),
				td.state.LastBlockHash(), td.state.Genesis().Hash())
			msg.Sign(signer)

			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			bundle := td.shouldPublishMessageWithThisType(t, td.network, message.TypeHelloAck)
			assert.Equal(t, bundle.Message.(*message.HelloAckMessage).ResponseCode, message.ResponseCodeOK)

			// Check if the peer info is updated
			p := td.sync.peerSet.GetPeer(pid)

			pub := signer.PublicKey().(*bls.PublicKey)
			assert.Equal(t, p.Status, peerset.StatusCodeConnected)
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
	assert.True(t, util.IsFlagSet(bdl.Message.(*message.HelloMessage).Services, services.New(services.Network)))
}
