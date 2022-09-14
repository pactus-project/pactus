package sync

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
)

func TestParsingHelloMessages(t *testing.T) {
	setup(t)

	t.Run("Receiving Hello message from a peer. Peer ID is not same as initiator.", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		pid := network.TestRandomPeerID()
		initiator := network.TestRandomPeerID()
		msg := message.NewHelloMessage(pid, "bad-genesis", 0, 0, tState.GenesisHash())
		signer.SignMsg(msg)
		assert.True(t, msg.PublicKey.EqualsTo(signer.PublicKey()))

		assert.Error(t, testReceivingNewMessage(tSync, msg, initiator))
		assert.Equal(t, tSync.peerSet.GetPeer(initiator).Status, peerset.StatusCodeBanned)
	})

	t.Run("Receiving Hello message from a peer. Genesis hash is wrong.", func(t *testing.T) {
		invGenHash := hash.GenerateTestHash()
		signer := bls.GenerateTestSigner()
		pid := network.TestRandomPeerID()
		msg := message.NewHelloMessage(pid, "bad-genesis", 0, 0, invGenHash)
		signer.SignMsg(msg)
		assert.True(t, msg.PublicKey.EqualsTo(signer.PublicKey()))

		assert.Error(t, testReceivingNewMessage(tSync, msg, pid))
		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeHello)
		checkPeerStatus(t, pid, peerset.StatusCodeBanned)
	})

	t.Run("Receiving Hello message from a peer. It should be acknowledged and updates the peer info", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		height := util.RandUint32(0)
		pid := network.TestRandomPeerID()
		msg := message.NewHelloMessage(pid, "kitty", height, message.FlagNodeNetwork, tState.GenesisHash())
		signer.SignMsg(msg)

		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))

		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeHello)

		// Check if the peer info is updated
		p := tSync.peerSet.GetPeer(pid)
		assert.Equal(t, p.Status, peerset.StatusCodeKnown)
		assert.Equal(t, p.Agent, version.Agent())
		assert.Equal(t, p.Moniker, "kitty")
		assert.True(t, p.PublicKey.EqualsTo(signer.PublicKey()))
		assert.Equal(t, p.PeerID, pid)
		assert.Equal(t, p.Height, height)
		assert.True(t, util.IsFlagSet(p.Flags, peerset.PeerFlagNodeNetwork))
	})

	t.Run("Receiving Hello-ack message from a peer. It should not be acknowledged, but update the peer info", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		pid := network.TestRandomPeerID()
		msg := message.NewHelloMessage(pid, "kitty", 0, message.FlagHelloAck, tState.GenesisHash())
		signer.SignMsg(msg)

		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))
		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeHello)
		checkPeerStatus(t, pid, peerset.StatusCodeKnown)
	})

	t.Run("Receiving Hello-ack message from a peer. Peer is ahead. It should request for blocks", func(t *testing.T) {
		tSync.peerSet.Clear()
		signer := bls.GenerateTestSigner()
		claimedHeight := tState.LastBlockHeight() + 5
		pid := network.TestRandomPeerID()
		msg := message.NewHelloMessage(pid, "kitty", claimedHeight, message.FlagHelloAck, tState.GenesisHash())
		signer.SignMsg(msg)

		assert.NoError(t, testReceivingNewMessage(tSync, msg, pid))
		shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksRequest)
		checkPeerStatus(t, pid, peerset.StatusCodeKnown)
		assert.Equal(t, tSync.peerSet.MaxClaimedHeight(), claimedHeight)
	})
}

func TestBroadcastingHelloMessages(t *testing.T) {
	setup(t)

	tSync.sayHello(true)

	bdl := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeHello)
	assert.True(t, util.IsFlagSet(bdl.Flags, bundle.BundleFlagHelloMessage))
	assert.True(t, util.IsFlagSet(bdl.Message.(*message.HelloMessage).Flags, message.FlagHelloAck))
}
