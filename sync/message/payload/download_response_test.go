package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/util"
)

func TestDownloadResponseType(t *testing.T) {
	p := &DownloadResponsePayload{}
	assert.Equal(t, p.Type(), PayloadTypeDownloadResponse)
}

func TestDownloadResponsePayload(t *testing.T) {
	t.Run("Invalid target", func(t *testing.T) {
		b, trxs := block.GenerateTestBlock(nil, nil)
		p1 := NewDownloadResponsePayload(ResponseCodeMoreBlocks, 1, "", 100, []*block.Block{b}, trxs)
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("Invalid from", func(t *testing.T) {
		b, trxs := block.GenerateTestBlock(nil, nil)
		p1 := NewDownloadResponsePayload(ResponseCodeMoreBlocks, 1, util.RandomPeerID(), -1, []*block.Block{b}, trxs)
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		b1, trxs1 := block.GenerateTestBlock(nil, nil)
		b2, trxs2 := block.GenerateTestBlock(nil, nil)
		trxs := append(trxs1, trxs2...)
		p1 := NewDownloadResponsePayload(ResponseCodeMoreBlocks, 1, util.RandomPeerID(), 100, []*block.Block{b1, b2}, trxs)
		assert.NoError(t, p1.SanityCheck())
	})
}

func TestDownloadResponseCode(t *testing.T) {
	t.Run("busy", func(t *testing.T) {
		p1 := NewDownloadResponsePayload(ResponseCodeBusy, 1, util.RandomPeerID(), 0, nil, nil)
		assert.NoError(t, p1.SanityCheck())
		assert.Equal(t, p1.(*DownloadResponsePayload).To(), 0)
		assert.True(t, p1.(*DownloadResponsePayload).IsRequestNotProcessed())
	})

	t.Run("rejected", func(t *testing.T) {
		p1 := NewDownloadResponsePayload(ResponseCodeRejected, 1, util.RandomPeerID(), 0, nil, nil)
		assert.NoError(t, p1.SanityCheck())
		assert.Equal(t, p1.(*DownloadResponsePayload).To(), 0)
		assert.True(t, p1.(*DownloadResponsePayload).IsRequestNotProcessed())
	})

	t.Run("OK", func(t *testing.T) {
		b1, trxs1 := block.GenerateTestBlock(nil, nil)
		b2, trxs2 := block.GenerateTestBlock(nil, nil)
		trxs := append(trxs1, trxs2...)
		p1 := NewDownloadResponsePayload(ResponseCodeMoreBlocks, 1, util.RandomPeerID(), 100, []*block.Block{b1, b2}, trxs)
		assert.NoError(t, p1.SanityCheck())
		assert.Equal(t, p1.(*DownloadResponsePayload).To(), 101)
		assert.False(t, p1.(*DownloadResponsePayload).IsRequestNotProcessed())
	})
}
