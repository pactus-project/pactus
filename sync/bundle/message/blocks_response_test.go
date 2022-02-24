package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func TestLatestBlocksResponseType(t *testing.T) {
	m := &BlocksResponseMessage{}
	assert.Equal(t, m.Type(), MessageTypeBlocksResponse)
}

func TestBlocksResponseMessage(t *testing.T) {
	t.Run("Invalid from", func(t *testing.T) {
		b, trxs := block.GenerateTestBlock(nil, nil)
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, 1, -1, []*block.Block{b}, trxs, nil)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("Invalid certificate", func(t *testing.T) {
		b, trxs := block.GenerateTestBlock(nil, nil)
		cert := block.GenerateTestCertificate(hash.UndefHash)
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, 1, 100, []*block.Block{b}, trxs, cert)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		b1, trxs1 := block.GenerateTestBlock(nil, nil)
		b2, trxs2 := block.GenerateTestBlock(nil, nil)
		trxs := append(trxs1, trxs2...)
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, 1, 100, []*block.Block{b1, b2}, trxs, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "100")
	})
}

func TestLatestBlocksResponseCode(t *testing.T) {
	t.Run("busy", func(t *testing.T) {
		m := NewBlocksResponseMessage(ResponseCodeBusy, 1, 0, nil, nil, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Equal(t, m.To(), 0)
		assert.True(t, m.IsRequestRejected())
	})

	t.Run("rejected", func(t *testing.T) {
		m := NewBlocksResponseMessage(ResponseCodeRejected, 1, 0, nil, nil, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Equal(t, m.To(), 0)
		assert.True(t, m.IsRequestRejected())
	})

	t.Run("OK", func(t *testing.T) {
		b1, trxs1 := block.GenerateTestBlock(nil, nil)
		b2, trxs2 := block.GenerateTestBlock(nil, nil)
		trxs := append(trxs1, trxs2...)

		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, 1, 100, []*block.Block{b1, b2}, trxs, nil)
		assert.NoError(t, m.SanityCheck())
		assert.Equal(t, m.To(), 101)
		assert.False(t, m.IsRequestRejected())
	})
}
