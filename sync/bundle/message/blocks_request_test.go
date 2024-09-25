package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatestBlocksRequestType(t *testing.T) {
	m := &BlocksRequestMessage{}
	assert.Equal(t, TypeBlocksRequest, m.Type())
}

func TestBlocksRequestMessage(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		m := NewBlocksRequestMessage(1, 0, 0)

		err := m.BasicCheck()
		assert.ErrorIs(t, err, BasicCheckError{Reason: "invalid height"})
	})
	t.Run("Invalid count", func(t *testing.T) {
		m := NewBlocksRequestMessage(1, 200, 0)

		err := m.BasicCheck()
		assert.ErrorIs(t, err, BasicCheckError{Reason: "count is zero"})
	})

	t.Run("OK", func(t *testing.T) {
		m := NewBlocksRequestMessage(1, 100, 7)

		assert.NoError(t, m.BasicCheck())
		assert.Equal(t, uint32(106), m.To())
		assert.Contains(t, m.String(), "100")
	})
}
