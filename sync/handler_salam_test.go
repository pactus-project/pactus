package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/version"
)

func TestParsingSalamMessages(t *testing.T) {
	setup(t)

	t.Run("Alice receives Salam message from a peer. Genesis hash is wrong. Alice should not handshake", func(t *testing.T) {
		invGenHash := hash.GenerateTestHash()
		pub, _ := bls.GenerateTestKeyPair()
		pld := payload.NewSalamPayload("bad-genesis", pub, invGenHash, 0, 0)
		pid := util.RandomPeerID()
		tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeBanned)
		shouldPublishPayloadWithThisTypeAndResponseCode(t, tAliceNet, payload.PayloadTypeAleyk, payload.ResponseCodeRejected)
	})

	t.Run("Alice receives Salam message from a peer. Genesis hash is Ok. Alice should update the peer info", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()

		pld := payload.NewSalamPayload("kitty", pub, tAliceState.GenHash, 3, 0x1)
		pid := util.RandomPeerID()
		tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeOK)
		shouldPublishPayloadWithThisTypeAndResponseCode(t, tAliceNet, payload.PayloadTypeAleyk, payload.ResponseCodeOK)
		assert.Equal(t, tBobSync.peerSet.MaxClaimedHeight(), tAliceState.LastBlockHeight())

		p := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, p.NodeVersion(), version.Version())
		assert.Equal(t, p.Moniker(), "kitty")
		assert.True(t, pub.EqualsTo(p.PublicKey()))
		assert.Equal(t, p.PeerID(), pid)
		assert.Equal(t, p.Height(), 3)
		assert.Equal(t, p.InitialBlockDownload(), true)
	})

	t.Run("Alice receives Salam message from a peer. Peer is ahead. Alice should request for blocks", func(t *testing.T) {
		tAliceSync.peerSet.Clear()
		pub, _ := bls.GenerateTestKeyPair()
		claimedHeight := tAliceState.LastBlockHeight() + 5
		pld := payload.NewSalamPayload("kitty", pub, tAliceState.GenHash, claimedHeight, 0)
		tAliceNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tAliceNet, payload.PayloadTypeAleyk, payload.ResponseCodeOK)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeLatestBlocksRequest)
		assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), claimedHeight)
	})
}
