package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func TestLatestBlocksResponseType(t *testing.T) {
	p := &BlocksResponsePayload{}
	assert.Equal(t, p.Type(), PayloadTypeBlocksResponse)
}

func TestBlocksResponsePayload(t *testing.T) {
	t.Run("Invalid from", func(t *testing.T) {
		b, trxs := block.GenerateTestBlock(nil, nil)
		p := NewBlocksResponsePayload(ResponseCodeMoreBlocks, 1, -1, []*block.Block{b}, trxs, nil)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid certificate", func(t *testing.T) {
		b, trxs := block.GenerateTestBlock(nil, nil)
		cert := block.GenerateTestCertificate(hash.UndefHash)
		p := NewBlocksResponsePayload(ResponseCodeMoreBlocks, 1, 100, []*block.Block{b}, trxs, cert)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		b1, trxs1 := block.GenerateTestBlock(nil, nil)
		b2, trxs2 := block.GenerateTestBlock(nil, nil)
		trxs := append(trxs1, trxs2...)
		p := NewBlocksResponsePayload(ResponseCodeMoreBlocks, 1, 100, []*block.Block{b1, b2}, trxs, nil)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "100")
	})
}

func TestLatestBlocksResponseCode(t *testing.T) {
	t.Run("busy", func(t *testing.T) {
		p := NewBlocksResponsePayload(ResponseCodeBusy, 1, 0, nil, nil, nil)

		assert.NoError(t, p.SanityCheck())
		assert.Equal(t, p.(*BlocksResponsePayload).To(), 0)
		assert.True(t, p.(*BlocksResponsePayload).IsRequestRejected())
	})

	t.Run("rejected", func(t *testing.T) {
		p := NewBlocksResponsePayload(ResponseCodeRejected, 1, 0, nil, nil, nil)

		assert.NoError(t, p.SanityCheck())
		assert.Equal(t, p.(*BlocksResponsePayload).To(), 0)
		assert.True(t, p.(*BlocksResponsePayload).IsRequestRejected())
	})

	t.Run("OK", func(t *testing.T) {
		b1, trxs1 := block.GenerateTestBlock(nil, nil)
		b2, trxs2 := block.GenerateTestBlock(nil, nil)
		trxs := append(trxs1, trxs2...)

		p := NewBlocksResponsePayload(ResponseCodeMoreBlocks, 1, 100, []*block.Block{b1, b2}, trxs, nil)
		assert.NoError(t, p.SanityCheck())
		assert.Equal(t, p.(*BlocksResponsePayload).To(), 101)
		assert.False(t, p.(*BlocksResponsePayload).IsRequestRejected())
	})
}
