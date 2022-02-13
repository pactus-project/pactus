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

	signer := bls.GenerateTestSigner()

	t.Run("Alice receives Aleyk message from a Peer. Peer has less blocks than Alice", func(t *testing.T) {
		pid := util.RandomPeerID()
		pld := payload.NewAleykPayload(pid, "kitty", 1, 0, tAlicePeerID, payload.ResponseCodeOK, "Welcome")
		signer.SignMsg(pld)
		simulatingReceiveingNewMessage(t, tAliceSync, pld, pid)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeGood)
		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
	})

	t.Run("Alice receives Aleyk message from a Peer. Peer has more blocks than Alice", func(t *testing.T) {
		tAliceSync.peerSet.Clear()
		claimedHeight := tAliceState.LastBlockHeight() + 5
		pid := util.RandomPeerID()
		pld := payload.NewAleykPayload(pid, "kitty", claimedHeight, 0, tAlicePeerID, payload.ResponseCodeOK, "Welcome")
		signer.SignMsg(pld)
		simulatingReceiveingNewMessage(t, tAliceSync, pld, pid)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeGood)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
	})

	t.Run("Alice receives not welcoming Aleyk message from a peer", func(t *testing.T) {
		pid := util.RandomPeerID()
		pld := payload.NewAleykPayload(pid, "kitty", 1, 0, tAlicePeerID, payload.ResponseCodeRejected, "Not Welcome!")
		signer.SignMsg(pld)
		simulatingReceiveingNewMessage(t, tAliceSync, pld, pid)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeUnknown)
	})

	t.Run("Alice receives Aleyk message from a peer but not targeted Alice", func(t *testing.T) {
		pid := util.RandomPeerID()
		pld := payload.NewAleykPayload(pid, "kitty", 1, 0, util.RandomPeerID(), payload.ResponseCodeOK, "Welcome")
		signer.SignMsg(pld)
		simulatingReceiveingNewMessage(t, tAliceSync, pld, pid)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeUnknown)
	})

	t.Run("Alice listens to Aleyk messages", func(t *testing.T) {
		pid := util.RandomPeerID()
		pld := payload.NewAleykPayload(pid, "kitty", 1, 0, util.RandomPeerID(), payload.ResponseCodeRejected, "Not Welcome!")
		signer.SignMsg(pld)
		simulatingReceiveingNewMessage(t, tAliceSync, pld, pid)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeUnknown)
		assert.Equal(t, peer.PublicKey().String(), signer.PublicKey().String())
	})
}
