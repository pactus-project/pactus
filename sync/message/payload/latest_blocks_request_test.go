package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/util"
)

func TestLatestBlocksRequestType(t *testing.T) {
	p := &LatestBlocksRequestPayload{}
	assert.Equal(t, p.Type(), PayloadTypeLatestBlocksRequest)
}

func TestLatestBlocksRequestPayload(t *testing.T) {
	t.Run("Invalid target", func(t *testing.T) {
		p := NewLatestBlocksRequestPayload(1, "", 100, 200)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid from", func(t *testing.T) {
		p := NewLatestBlocksRequestPayload(1, util.RandomPeerID(), -100, 200)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid range", func(t *testing.T) {
		p := NewLatestBlocksRequestPayload(1, util.RandomPeerID(), 200, 100)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		p := NewLatestBlocksRequestPayload(1, util.RandomPeerID(), 100, 200)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "100")
	})
}
