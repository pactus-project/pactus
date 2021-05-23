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
		p := NewDownloadResponsePayload(ResponseCodeMoreBlocks, 1, "", 100, []*block.Block{b}, trxs)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid from", func(t *testing.T) {
		b, trxs := block.GenerateTestBlock(nil, nil)
		p := NewDownloadResponsePayload(ResponseCodeMoreBlocks, 1, util.RandomPeerID(), -1, []*block.Block{b}, trxs)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		b1, trxs1 := block.GenerateTestBlock(nil, nil)
		b2, trxs2 := block.GenerateTestBlock(nil, nil)
		trxs := append(trxs1, trxs2...)
		p := NewDownloadResponsePayload(ResponseCodeMoreBlocks, 1, util.RandomPeerID(), 100, []*block.Block{b1, b2}, trxs)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "100")
	})
}

func TestDownloadResponseCode(t *testing.T) {
	t.Run("busy", func(t *testing.T) {
		p := NewDownloadResponsePayload(ResponseCodeBusy, 1, util.RandomPeerID(), 0, nil, nil)

		assert.NoError(t, p.SanityCheck())
		assert.Equal(t, p.(*DownloadResponsePayload).To(), 0)
		assert.True(t, p.(*DownloadResponsePayload).IsRequestNotProcessed())
	})

	t.Run("rejected", func(t *testing.T) {
		p := NewDownloadResponsePayload(ResponseCodeRejected, 1, util.RandomPeerID(), 0, nil, nil)

		assert.NoError(t, p.SanityCheck())
		assert.Equal(t, p.(*DownloadResponsePayload).To(), 0)
		assert.True(t, p.(*DownloadResponsePayload).IsRequestNotProcessed())
	})

	t.Run("OK", func(t *testing.T) {
		b1, trxs1 := block.GenerateTestBlock(nil, nil)
		b2, trxs2 := block.GenerateTestBlock(nil, nil)
		trxs := append(trxs1, trxs2...)
		p := NewDownloadResponsePayload(ResponseCodeMoreBlocks, 1, util.RandomPeerID(), 100, []*block.Block{b1, b2}, trxs)

		assert.NoError(t, p.SanityCheck())
		assert.Equal(t, p.(*DownloadResponsePayload).To(), 101)
		assert.False(t, p.(*DownloadResponsePayload).IsRequestNotProcessed())
	})
}
