package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

func TestParsingBlockAnnounceMessages(t *testing.T) {
	setup(t)

	t.Run("Bob should not broadcast block announce message because he is not in the committee", func(t *testing.T) {
		addMoreBlocksForBobAndAnnounceLastBlock(t, 1)

		shouldNotPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlockAnnounce)
	})

	joinBobToTheSet(t)

	t.Run("Bob should broadcast block announce message because he is in the committee", func(t *testing.T) {
		addMoreBlocksForBobAndAnnounceLastBlock(t, 1)

		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlockAnnounce)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeLatestBlocksRequest)

		assert.Equal(t, tAliceSync.peerSet.FindHighestPeer().PeerID(), tBobPeerID)
	})
}
