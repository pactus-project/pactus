package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatestBlocksRequestType(t *testing.T) {
	p := &BlocksRequestPayload{}
	assert.Equal(t, p.Type(), PayloadTypeBlocksRequest)
}

func TestBlocksRequestPayload(t *testing.T) {
	t.Run("Invalid from", func(t *testing.T) {
		p := NewBlocksRequestPayload(1, -100, 200)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid range", func(t *testing.T) {
		p := NewBlocksRequestPayload(1, 200, 100)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		p := NewBlocksRequestPayload(1, 100, 200)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "100")
	})
}
