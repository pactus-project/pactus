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

	t.Run("Alice receives Aleyk message from a Peer. Peer has less blocks than Alice", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()
		from := util.RandomPeerID()
		pld := payload.NewAleykPayload(tAlicePeerID, payload.ResponseCodeOK, "Welcome", "kitty", pub, 1, 0)
		simulatingReceiveingNewMessage(t, tAliceSync, pld, from)

		peer := tAliceSync.peerSet.GetPeer(from)
		assert.Equal(t, peer.Status(), peerset.StatusCodeOK)
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
	})

	t.Run("Alice receives Aleyk message from a Peer. Peer has more blocks than Alice", func(t *testing.T) {
		tAliceSync.peerSet.Clear()
		pub, _ := bls.GenerateTestKeyPair()
		from := util.RandomPeerID()
		claimedHeight := tAliceState.LastBlockHeight() + 5
		pld := payload.NewAleykPayload(tAlicePeerID, payload.ResponseCodeOK, "Welcome", "kitty", pub, claimedHeight, 0)
		simulatingReceiveingNewMessage(t, tAliceSync, pld, from)

		peer := tAliceSync.peerSet.GetPeer(from)
		assert.Equal(t, peer.Status(), peerset.StatusCodeOK)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
	})

	t.Run("Alice receives not welcoming Aleyk message from a peer", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()
		from := util.RandomPeerID()
		pld := payload.NewAleykPayload(tAlicePeerID, payload.ResponseCodeRejected, "Not Welcome!", "kitty", pub, 1, 0)
		simulatingReceiveingNewMessage(t, tAliceSync, pld, from)

		peer := tAliceSync.peerSet.GetPeer(from)
		assert.Equal(t, peer.Status(), peerset.StatusCodeBanned)
	})

	t.Run("Alice receives Aleyk message from a peer but not targeted Alice", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()
		from := util.RandomPeerID()
		pld := payload.NewAleykPayload(util.RandomPeerID(), payload.ResponseCodeOK, "Welcome", "kitty", pub, 1, 0)
		simulatingReceiveingNewMessage(t, tAliceSync, pld, from)

		peer := tAliceSync.peerSet.GetPeer(from)
		assert.Equal(t, peer.Status(), peerset.StatusCodeUnknown)
	})

	t.Run("Alice eavesdrops Aleyk messages", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()
		from := util.RandomPeerID()
		pld := payload.NewAleykPayload(util.RandomPeerID(), payload.ResponseCodeRejected, "Not Welcome!", "kitty", pub, 1, 0)
		simulatingReceiveingNewMessage(t, tAliceSync, pld, from)

		peer := tAliceSync.peerSet.GetPeer(from)
		assert.Equal(t, peer.Status(), peerset.StatusCodeUnknown)
		assert.Equal(t, peer.PublicKey().String(), pub.String())
	})
}
