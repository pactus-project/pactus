package sync

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
)

func TestParsingHelloMessages(t *testing.T) {
	td := setup(t, nil)

	t.Run("Receiving Hello message from a peer. Peer ID is not same as initiator.",
		func(t *testing.T) {
			signer := td.RandomSigner()
			pid := td.RandomPeerID()
			initiator := td.RandomPeerID()
			msg := message.NewHelloMessage(pid, "bad-genesis", 0, 0,
				td.state.LastBlockHash(), td.state.Genesis().Hash())
			msg.Sign(signer)

			assert.Error(t, td.receivingNewMessage(td.sync, msg, initiator))
			assert.Equal(t, td.sync.peerSet.GetPeer(initiator).Status, peerset.StatusCodeBanned)
		})

	t.Run("Receiving Hello message from a peer. Genesis hash is wrong.",
		func(t *testing.T) {
			invGenHash := td.RandomHash()
			signer := td.RandomSigner()
			pid := td.RandomPeerID()
			msg := message.NewHelloMessage(pid, "bad-genesis", 0, 0,
				td.state.LastBlockHash(), invGenHash)
			msg.Sign(signer)

			assert.Error(t, td.receivingNewMessage(td.sync, msg, pid))
			td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeHello)
			td.checkPeerStatus(t, pid, peerset.StatusCodeBanned)
		})

	t.Run("Receiving Hello message from a peer. It should be acknowledged and updates the peer info",
		func(t *testing.T) {
			signer := td.RandomSigner()
			height := td.RandUint32(td.state.LastBlockHeight())
			pid := td.RandomPeerID()
			msg := message.NewHelloMessage(pid, "kitty", height, message.FlagNodeNetwork,
				td.state.LastBlockHash(), td.state.Genesis().Hash())
			msg.Sign(signer)

			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))

			td.shouldPublishMessageWithThisType(t, td.network, message.TypeHello)

			// Check if the peer info is updated
			p := td.sync.peerSet.GetPeer(pid)

			pub := signer.PublicKey().(*bls.PublicKey)
			assert.Equal(t, p.Status, peerset.StatusCodeKnown)
			assert.Equal(t, p.Agent, version.Agent())
			assert.Equal(t, p.Moniker, "kitty")
			assert.Contains(t, p.ConsensusKeys, pub)
			assert.Equal(t, p.PeerID, pid)
			assert.Equal(t, p.Height, height)
			assert.True(t, util.IsFlagSet(p.Flags, peerset.PeerFlagNodeNetwork))
		})

	t.Run("Receiving Hello-ack message from a peer. It should not be acknowledged, but update the peer info",
		func(t *testing.T) {
			signer := td.RandomSigner()
			height := td.RandUint32(td.state.LastBlockHeight())
			pid := td.RandomPeerID()
			msg := message.NewHelloMessage(pid, "kitty", height, message.FlagHelloAck,
				td.state.LastBlockHash(), td.state.Genesis().Hash())
			msg.Sign(signer)

			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))
			td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeHello)
			td.checkPeerStatus(t, pid, peerset.StatusCodeKnown)

			// Check if the peer info is updated
			p := td.sync.peerSet.GetPeer(pid)
			assert.Equal(t, p.Height, height)
		})

	t.Run("Receiving Hello-ack message from a peer. Peer is ahead. It should request for blocks",
		func(t *testing.T) {
			td.sync.peerSet.Clear()
			signer := td.RandomSigner()
			claimedHeight := td.state.LastBlockHeight() + 5
			pid := td.RandomPeerID()
			msg := message.NewHelloMessage(pid, "kitty", claimedHeight, message.FlagHelloAck,
				td.state.LastBlockHash(), td.state.Genesis().Hash())
			msg.Sign(signer)

			assert.NoError(t, td.receivingNewMessage(td.sync, msg, pid))
			td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)
			td.checkPeerStatus(t, pid, peerset.StatusCodeKnown)
			assert.Equal(t, td.sync.peerSet.MaxClaimedHeight(), claimedHeight)
		})
}

func TestBroadcastingHelloMessages(t *testing.T) {
	td := setup(t, nil)

	td.sync.sayHello(true)

	bdl := td.shouldPublishMessageWithThisType(t, td.network, message.TypeHello)
	assert.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagHelloMessage))
	assert.True(t, util.IsFlagSet(bdl.Message.(*message.HelloMessage).Flags, message.FlagHelloAck))
}
