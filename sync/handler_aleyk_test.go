package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingAleykMessages(t *testing.T) {
	setup(t)

	pub, prv := bls.GenerateTestKeyPair()
	sig := prv.Sign(pub.RawBytes())

	t.Run("Alice receives Aleyk message from a Peer. Peer has less blocks than Alice", func(t *testing.T) {
		from := util.RandomPeerID()
		pld := payload.NewAleykPayload("kitty", pub, sig, 1, 0, tAlicePeerID, payload.ResponseCodeOK, "Welcome")
		simulatingReceiveingNewMessage(t, tAliceSync, pld, from)

		peer := tAliceSync.peerSet.GetPeer(from)
		assert.Equal(t, peer.Status(), peerset.StatusCodeOK)
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
	})

	t.Run("Alice receives Aleyk message from a Peer. Peer has more blocks than Alice", func(t *testing.T) {
		tAliceSync.peerSet.Clear()
		from := util.RandomPeerID()
		claimedHeight := tAliceState.LastBlockHeight() + 5
		pld := payload.NewAleykPayload("kitty", pub, sig, claimedHeight, 0, tAlicePeerID, payload.ResponseCodeOK, "Welcome")
		simulatingReceiveingNewMessage(t, tAliceSync, pld, from)

		peer := tAliceSync.peerSet.GetPeer(from)
		assert.Equal(t, peer.Status(), peerset.StatusCodeOK)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
	})

	t.Run("Alice receives not welcoming Aleyk message from a peer", func(t *testing.T) {
		from := util.RandomPeerID()
		pld := payload.NewAleykPayload("kitty", pub, sig, 1, 0, tAlicePeerID, payload.ResponseCodeRejected, "Not Welcome!")
		simulatingReceiveingNewMessage(t, tAliceSync, pld, from)

		peer := tAliceSync.peerSet.GetPeer(from)
		assert.Equal(t, peer.Status(), peerset.StatusCodeBanned)
	})

	t.Run("Alice receives Aleyk message from a peer but not targeted Alice", func(t *testing.T) {
		from := util.RandomPeerID()
		pld := payload.NewAleykPayload("kitty", pub, sig, 1, 0, util.RandomPeerID(), payload.ResponseCodeOK, "Welcome")
		simulatingReceiveingNewMessage(t, tAliceSync, pld, from)

		peer := tAliceSync.peerSet.GetPeer(from)
		assert.Equal(t, peer.Status(), peerset.StatusCodeUnknown)
	})

	t.Run("Alice listens to Aleyk messages", func(t *testing.T) {
		from := util.RandomPeerID()
		pld := payload.NewAleykPayload("kitty", pub, sig, 1, 0, util.RandomPeerID(), payload.ResponseCodeRejected, "Not Welcome!")
		simulatingReceiveingNewMessage(t, tAliceSync, pld, from)

		peer := tAliceSync.peerSet.GetPeer(from)
		assert.Equal(t, peer.Status(), peerset.StatusCodeUnknown)
		assert.Equal(t, peer.PublicKey().String(), pub.String())
	})
}
