package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingAleykMessages(t *testing.T) {
	setup(t)

	t.Run("Alice receives Aleyk message from a peer. The peer is behind alice. Alice should not request for blocks", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		pld := payload.NewAleykPayload(payload.ResponseCodeOK, "Welcome!", "kitty", pub, 1, 0)
		tAliceNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeLatestBlocksRequest)
	})

	t.Run("Alice receives Aleyk message from a peer. The peer is ahead of alice. Alice should request for new blocks", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		claimedHeight := tAliceState.LastBlockHeight() + 5
		pld := payload.NewAleykPayload(payload.ResponseCodeOK, "Welcome!", "kitty", pub, claimedHeight, 0)
		tAliceNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeLatestBlocksRequest)
		assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), claimedHeight)
	})

	t.Run("Alice receives Aleyk message from a peer. The peer is at the same height with alice. Alice should not request for blocks", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		claimedHeight := tAliceState.LastBlockHeight()
		pld := payload.NewAleykPayload(payload.ResponseCodeOK, "Welcome!", "kitty", pub, claimedHeight, 0)
		tAliceNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeLatestBlocksRequest)
	})
}
