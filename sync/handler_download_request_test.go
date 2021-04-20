package sync

import (
	"testing"

	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func TestDownloadBlocksRequestMessages(t *testing.T) {
	setup(t)
	disableHeartbeat(t)

	t.Run("Alice should not pay attention to requests for other peers", func(t *testing.T) {
		pld := payload.NewDownloadRequestPayload(6, util.RandomPeerID(), 0, 10)
		tAliceNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeDownloadResponse)
	})

	t.Run("Alice should not pay attention to requests with invalid ranges", func(t *testing.T) {
		pld := payload.NewDownloadRequestPayload(6, tAlicePeerID, 1000, 2000)
		tAliceNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldNotPublishPayloadWithThisType(t, tAliceNet, payload.PayloadTypeDownloadResponse)
	})

	t.Run("Alice doesn't have request blocks", func(t *testing.T) {
		pld := payload.NewDownloadRequestPayload(6, tAlicePeerID, 100, 105)
		tAliceNet.ReceivingMessageFromOtherPeer(util.RandomPeerID(), pld)

		shouldPublishPayloadWithThisTypeAndResponseCode(t, tAliceNet, payload.PayloadTypeDownloadResponse, payload.ResponseCodeNoMoreBlocks)
	})
}
