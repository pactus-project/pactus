package message

import (
	"testing"

	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestLatestBlocksRequestType(t *testing.T) {
	m := &BlocksRequestMessage{}
	assert.Equal(t, m.Type(), MessageTypeBlocksRequest)
}

func TestBlocksRequestMessage(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		m := NewBlocksRequestMessage(1, 0, 0)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidHeight)
	})
	t.Run("Invalid count", func(t *testing.T) {
		m := NewBlocksRequestMessage(1, 200, 0)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidMessage)
	})

	t.Run("OK", func(t *testing.T) {
		m := NewBlocksRequestMessage(1, 100, 7)

		assert.NoError(t, m.SanityCheck())
		assert.Equal(t, m.To(), uint32(106))
		assert.Contains(t, m.Fingerprint(), "100")
	})
}
