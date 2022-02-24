package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatestBlocksRequestType(t *testing.T) {
	m := &BlocksRequestMessage{}
	assert.Equal(t, m.Type(), MessageTypeBlocksRequest)
}

func TestBlocksRequestMessage(t *testing.T) {
	t.Run("Invalid from", func(t *testing.T) {
		m := NewBlocksRequestMessage(1, -100, 200)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("Invalid range", func(t *testing.T) {
		m := NewBlocksRequestMessage(1, 200, 100)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		m := NewBlocksRequestMessage(1, 100, 200)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "100")
	})
}
