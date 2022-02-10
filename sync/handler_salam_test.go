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
		pub, prv := bls.GenerateTestKeyPair()
		sig := prv.Sign(pub.RawBytes())
		pld := payload.NewSalamPayload("bad-genesis", pub, sig, 0, 0, invGenHash)
		pid := util.RandomPeerID()
		simulatingReceiveingNewMessage(t, tAliceSync, pld, pid)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeBanned)
		msg := shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeAleyk)
		assert.True(t, msg.Payload.(*payload.AleykPayload).PublicKey.EqualsTo(tAliceSync.signer.PublicKey()))
		assert.Equal(t, msg.Payload.(*payload.AleykPayload).ResponseCode, payload.ResponseCodeRejected)
	})

	t.Run("Alice receives Salam message from a peer. Genesis hash is Ok. Alice should update the peer info", func(t *testing.T) {
		pub, prv := bls.GenerateTestKeyPair()
		sig := prv.Sign(pub.RawBytes())

		pld := payload.NewSalamPayload("kitty", pub, sig, 3, 0x1, tAliceState.GenHash)
		pid := util.RandomPeerID()
		simulatingReceiveingNewMessage(t, tAliceSync, pld, pid)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeOK)
		msg := shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeAleyk)
		assert.Equal(t, msg.Payload.(*payload.AleykPayload).ResponseCode, payload.ResponseCodeOK)
		assert.Equal(t, tBobSync.peerSet.MaxClaimedHeight(), tAliceState.LastBlockHeight())

		p := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, p.Agent(), version.Agent())
		assert.Equal(t, p.Moniker(), "kitty")
		assert.True(t, pub.EqualsTo(p.PublicKey()))
		assert.Equal(t, p.PeerID(), pid)
		assert.Equal(t, p.Height(), 3)
		assert.Equal(t, p.InitialBlockDownload(), true)
	})

	t.Run("Alice receives Salam message from a peer. Peer is ahead. Alice should request for blocks", func(t *testing.T) {
		tAliceSync.peerSet.Clear()
		pub, prv := bls.GenerateTestKeyPair()
		sig := prv.Sign(pub.RawBytes())
		claimedHeight := tAliceState.LastBlockHeight() + 5
		pld := payload.NewSalamPayload("kitty", pub, sig, claimedHeight, 0, tAliceState.GenHash)
		simulatingReceiveingNewMessage(t, tAliceSync, pld, util.RandomPeerID())

		msg := shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeAleyk)
		assert.Equal(t, msg.Payload.(*payload.AleykPayload).ResponseCode, payload.ResponseCodeOK)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
		assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), claimedHeight)
	})
}
