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
	t.Run("Invalid range", func(t *testing.T) {
		m := NewBlocksRequestMessage(1, 200, 100)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidHeight)
	})

	t.Run("OK", func(t *testing.T) {
		m := NewBlocksRequestMessage(1, 100, 200)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "100")
	})
}
