package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/version"
)

func TestParsingHelloMessages(t *testing.T) {
	setup(t)

	t.Run("Receiving Hello message from a peer. Genesis hash is wrong.", func(t *testing.T) {
		invGenHash := hash.GenerateTestHash()
		signer := bls.GenerateTestSigner()
		pid := util.RandomPeerID()
		pld := payload.NewHelloPayload(pid, "bad-genesis", 0, payload.FlagNeedResponse, invGenHash)
		signer.SignMsg(pld)
		assert.True(t, pld.PublicKey.EqualsTo(signer.PublicKey()))

		assert.Error(t, testReceiveingNewMessage(tSync, pld, pid))
		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeHello)
		checkPeerStatus(t, pid, peerset.StatusCodeBanned)
	})

	t.Run("Receiving Hello message from a peer. It should be acknowledged and updates the peer info", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		height := util.RandInt(0)
		pid := util.RandomPeerID()
		pld := payload.NewHelloPayload(pid, "kitty", height, payload.FlagNeedResponse|payload.FlagInitialBlockDownload, tState.GenHash)
		signer.SignMsg(pld)

		assert.NoError(t, testReceiveingNewMessage(tSync, pld, pid))

		msg := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeHello)
		assert.False(t, util.IsFlagSet(msg.Payload.(*payload.HelloPayload).Flags, payload.FlagNeedResponse))

		// Check if the peer info is updated
		p := tSync.peerSet.GetPeer(pid)
		assert.Equal(t, p.Status(), peerset.StatusCodeKnown)
		assert.Equal(t, p.Agent(), version.Agent())
		assert.Equal(t, p.Moniker(), "kitty")
		assert.True(t, signer.PublicKey().EqualsTo(p.PublicKey()))
		assert.Equal(t, p.PeerID(), pid)
		assert.Equal(t, p.Height(), height)
		assert.Equal(t, p.InitialBlockDownload(), true)
	})

	t.Run("Receiving Hello-ack message from a peer. It should not be acknowledged, but update the peer info", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		pid := util.RandomPeerID()
		pld := payload.NewHelloPayload(pid, "kitty", 0, payload.FlagInitialBlockDownload, tState.GenHash)
		signer.SignMsg(pld)

		assert.NoError(t, testReceiveingNewMessage(tSync, pld, pid))
		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeHello)
		checkPeerStatus(t, pid, peerset.StatusCodeKnown)
	})

	t.Run("Receiving Hello-ack message from a peer. Peer is ahead. It should request for blocks", func(t *testing.T) {
		tSync.peerSet.Clear()
		signer := bls.GenerateTestSigner()
		claimedHeight := tState.LastBlockHeight() + 5
		pid := util.RandomPeerID()
		pld := payload.NewHelloPayload(pid, "kitty", claimedHeight, 0, tState.GenHash)
		signer.SignMsg(pld)

		assert.NoError(t, testReceiveingNewMessage(tSync, pld, pid))
		shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlocksRequest)
		checkPeerStatus(t, pid, peerset.StatusCodeKnown)
		assert.Equal(t, tSync.peerSet.MaxClaimedHeight(), claimedHeight)
	})
}

func TestBroadcastingHelloMessages(t *testing.T) {
	setup(t)

	tSync.sayHello(true)

	msg := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeHello)
	assert.True(t, util.IsFlagSet(msg.Flags, message.MsgFlagHelloMessage))
	assert.True(t, util.IsFlagSet(msg.Payload.(*payload.HelloPayload).Flags, payload.FlagNeedResponse))
}
