package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingAleykMessages(t *testing.T) {
	setup(t)

	t.Run("Alice receives Aleyk message a Peer.", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		pid := util.RandomPeerID()
		pld := payload.NewAleykPayload(tAlicePeerID, payload.ResponseCodeOK, "Welcome", "kitty", pub, 1, 0)
		tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeOK)
	})

	t.Run("Alice receives not welcoming Aleyk message from a peer", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		pid := util.RandomPeerID()
		claimedHeight := tAliceState.LastBlockHeight() + 5
		pld := payload.NewAleykPayload(tAlicePeerID, payload.ResponseCodeRejected, "Not Welcome!", "kitty", pub, claimedHeight, 0)
		tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeBanned)
	})

	t.Run("Alice receives Aleyk message from a peer but not targeted Alice", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		pid := util.RandomPeerID()
		claimedHeight := tAliceState.LastBlockHeight() + 5
		pld := payload.NewAleykPayload(util.RandomPeerID(), payload.ResponseCodeOK, "Welcome", "kitty", pub, claimedHeight, 0)
		tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeUnknown)
	})

	t.Run("Alice eavesdrops Aleyk messages", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		pid := util.RandomPeerID()
		claimedHeight := tAliceState.LastBlockHeight() + 5
		pld := payload.NewAleykPayload(util.RandomPeerID(), payload.ResponseCodeRejected, "Not Welcome!", "kitty", pub, claimedHeight, 0)
		tAliceNet.ReceivingMessageFromOtherPeer(pid, pld)

		peer := tAliceSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Status(), peerset.StatusCodeUnknown)
		assert.Equal(t, peer.PublicKey().String(), pub.String())
	})

}
