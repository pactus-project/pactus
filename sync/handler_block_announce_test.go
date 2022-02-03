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

	joinBobToCommittee(t)

	t.Run("Bob should broadcast block announce message because he is in the committee", func(t *testing.T) {
		addMoreBlocksForBobAndAnnounceLastBlock(t, 1)

		msg1 := shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlockAnnounce)
		assert.Equal(t, msg1.Payload.(*payload.BlockAnnouncePayload).Height, tBobState.LastBlockHeight())
		assert.Equal(t, msg1.Payload.(*payload.BlockAnnouncePayload).Block.Hash(), tBobState.LastBlockHash())
		assert.Equal(t, msg1.Payload.(*payload.BlockAnnouncePayload).Certificate.BlockHash(), tBobState.LastBlockCertificate.BlockHash())

		msg2 := shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeBlocksRequest)
		assert.Equal(t, msg2.Payload.(*payload.BlocksRequestPayload).From, tBobState.LastBlockHeight()-1)
		assert.Equal(t, msg2.Payload.(*payload.BlocksRequestPayload).To, tBobState.LastBlockHeight())
	})
}
