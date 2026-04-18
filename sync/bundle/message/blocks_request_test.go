package message

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLatestBlocksRequestType(t *testing.T) {
	msg := &BlocksRequestMessage{}
	assert.Equal(t, TypeBlocksRequest, msg.Type())
}

func TestBlocksRequestMessage(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		msg := NewBlocksRequestMessage(1, 0, 0)

		err := msg.BasicCheck()
		require.ErrorIs(t, err, BasicCheckError{Reason: "invalid height"})
	})
	t.Run("Invalid count", func(t *testing.T) {
		msg := NewBlocksRequestMessage(1, 200, 0)

		err := msg.BasicCheck()
		require.ErrorIs(t, err, BasicCheckError{Reason: "count is zero"})
	})

	t.Run("OK", func(t *testing.T) {
		msg := NewBlocksRequestMessage(1, 100, 7)

		require.NoError(t, msg.BasicCheck())
		assert.Equal(t, types.Height(106), msg.To())
		assert.Contains(t, msg.LogString(), "100")
	})
}
