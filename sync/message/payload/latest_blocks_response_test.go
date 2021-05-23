package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

func TestLatestBlocksResponseType(t *testing.T) {
	p := &LatestBlocksResponsePayload{}
	assert.Equal(t, p.Type(), PayloadTypeLatestBlocksResponse)
}

func TestLatestBlocksResponsePayload(t *testing.T) {
	t.Run("Invalid target", func(t *testing.T) {
		b, trxs := block.GenerateTestBlock(nil, nil)
		p1 := NewLatestBlocksResponsePayload(ResponseCodeMoreBlocks, 1, "", 100, []*block.Block{b}, trxs, nil)
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("Invalid from", func(t *testing.T) {
		b, trxs := block.GenerateTestBlock(nil, nil)
		p1 := NewLatestBlocksResponsePayload(ResponseCodeMoreBlocks, 1, util.RandomPeerID(), -1, []*block.Block{b}, trxs, nil)
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("Invalid certificate", func(t *testing.T) {
		b, trxs := block.GenerateTestBlock(nil, nil)
		cert := block.GenerateTestCertificate(crypto.UndefHash)
		p1 := NewLatestBlocksResponsePayload(ResponseCodeMoreBlocks, 1, util.RandomPeerID(), 100, []*block.Block{b}, trxs, cert)
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		b1, trxs1 := block.GenerateTestBlock(nil, nil)
		b2, trxs2 := block.GenerateTestBlock(nil, nil)
		trxs := append(trxs1, trxs2...)
		p1 := NewLatestBlocksResponsePayload(ResponseCodeMoreBlocks, 1, util.RandomPeerID(), 100, []*block.Block{b1, b2}, trxs, nil)
		assert.NoError(t, p1.SanityCheck())
	})
}

func TestLatestBlocksResponseCode(t *testing.T) {
	t.Run("busy", func(t *testing.T) {
		p1 := NewLatestBlocksResponsePayload(ResponseCodeBusy, 1, util.RandomPeerID(), 0, nil, nil, nil)
		assert.NoError(t, p1.SanityCheck())
		assert.Equal(t, p1.(*LatestBlocksResponsePayload).To(), 0)
		assert.True(t, p1.(*LatestBlocksResponsePayload).IsRequestNotProcessed())
	})

	t.Run("rejected", func(t *testing.T) {
		p1 := NewLatestBlocksResponsePayload(ResponseCodeRejected, 1, util.RandomPeerID(), 0, nil, nil, nil)
		assert.NoError(t, p1.SanityCheck())
		assert.Equal(t, p1.(*LatestBlocksResponsePayload).To(), 0)
		assert.True(t, p1.(*LatestBlocksResponsePayload).IsRequestNotProcessed())
	})

	t.Run("OK", func(t *testing.T) {
		b1, trxs1 := block.GenerateTestBlock(nil, nil)
		b2, trxs2 := block.GenerateTestBlock(nil, nil)
		trxs := append(trxs1, trxs2...)
		p1 := NewLatestBlocksResponsePayload(ResponseCodeMoreBlocks, 1, util.RandomPeerID(), 100, []*block.Block{b1, b2}, trxs, nil)
		assert.NoError(t, p1.SanityCheck())
		assert.Equal(t, p1.(*LatestBlocksResponsePayload).To(), 101)
		assert.False(t, p1.(*LatestBlocksResponsePayload).IsRequestNotProcessed())
	})
}
