package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatestBlocksRequestType(t *testing.T) {
	msg := &BlocksRequestMessage{}
	assert.Equal(t, TypeBlocksRequest, msg.Type())
}

func TestBlocksRequestMessage(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		msg := NewBlocksRequestMessage(1, 0, 0)

		err := msg.BasicCheck()
		assert.ErrorIs(t, err, BasicCheckError{Reason: "invalid height"})
	})
	t.Run("Invalid count", func(t *testing.T) {
		msg := NewBlocksRequestMessage(1, 200, 0)

		err := msg.BasicCheck()
		assert.ErrorIs(t, err, BasicCheckError{Reason: "count is zero"})
	})

	t.Run("OK", func(t *testing.T) {
		msg := NewBlocksRequestMessage(1, 100, 7)

		assert.NoError(t, msg.BasicCheck())
		assert.Equal(t, uint32(106), msg.To())
		assert.Contains(t, msg.String(), "100")
	})
}
