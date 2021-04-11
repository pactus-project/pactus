package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

func TestOneBlockBehind(t *testing.T) {
	setup(t)

	t.Run("Bob is not in the committee. Bob commits one block. Bob should broadcasts heartbeat. Alice should ask for the last block.", func(t *testing.T) {
		addMoreBlocksForBob(t, 1)

		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeHeartBeat)
		shouldPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeLatestBlocksRequest)
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeLatestBlocksResponse) // 22
		shouldPublishPayloadWithThisType(t, tBobNet, payload.PayloadTypeLatestBlocksResponse) // No more block

		assert.Equal(t, tAliceState.LastBlockHash(), tBobState.LastBlockHash())
		assert.Equal(t, tAliceState.LastBlockHeight(), tBobState.LastBlockHeight())
	})
}
