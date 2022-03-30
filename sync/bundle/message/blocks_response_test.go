package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
)

func TestLatestBlocksResponseType(t *testing.T) {
	m := &BlocksResponseMessage{}
	assert.Equal(t, m.Type(), MessageTypeBlocksResponse)
}

func TestBlocksResponseMessage(t *testing.T) {
	t.Run("Invalid from", func(t *testing.T) {
		b := block.GenerateTestBlock(nil, nil)
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, 1, -1, []*block.Block{b}, nil)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidHeight)
	})

	t.Run("Invalid certificate", func(t *testing.T) {
		b := block.GenerateTestBlock(nil, nil)
		c := block.NewCertificate(-1, nil, nil, nil)
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, 1, 100, []*block.Block{b}, c)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidRound)
	})

	t.Run("OK", func(t *testing.T) {
		b1 := block.GenerateTestBlock(nil, nil)
		b2 := block.GenerateTestBlock(nil, nil)
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, 1, 100, []*block.Block{b1, b2}, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "100")
	})
}

func TestLatestBlocksResponseCode(t *testing.T) {
	t.Run("busy", func(t *testing.T) {
		m := NewBlocksResponseMessage(ResponseCodeBusy, 1, 0, nil, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Equal(t, m.To(), int32(0))
		assert.True(t, m.IsRequestRejected())
	})

	t.Run("rejected", func(t *testing.T) {
		m := NewBlocksResponseMessage(ResponseCodeRejected, 1, 0, nil, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Equal(t, m.To(), int32(0))
		assert.True(t, m.IsRequestRejected())
	})

	t.Run("OK", func(t *testing.T) {
		b1 := block.GenerateTestBlock(nil, nil)
		b2 := block.GenerateTestBlock(nil, nil)

		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, 1, 100, []*block.Block{b1, b2}, nil)
		assert.NoError(t, m.SanityCheck())
		assert.Equal(t, m.To(), int32(101))
		assert.False(t, m.IsRequestRejected())
	})
}
