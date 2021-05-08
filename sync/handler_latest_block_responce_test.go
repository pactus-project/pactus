package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

func TestOneBlockShorter(t *testing.T) {
	setup(t)
	disableHeartbeat(t)

	t.Run("Bob commits two blocks. Alice should ask for the lastest block.", func(t *testing.T) {
		joinBobToCommittee(t)
		addMoreBlocksForBobAndAnnounceLastBlock(t, 2)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeBlockAnnounce)

		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeLatestBlocksRequest)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeLatestBlocksResponse) // 22
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeLatestBlocksResponse) // No more block

		assert.Equal(t, tAliceState.LastBlockHash(), tBobState.LastBlockHash())
		assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
	})
}
