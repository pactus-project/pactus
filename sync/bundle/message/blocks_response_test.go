package message

import (
	"testing"

	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestLatestBlocksResponseType(t *testing.T) {
	msg := &BlocksResponseMessage{}
	assert.Equal(t, TypeBlocksResponse, msg.Type())
}

func TestBlocksResponseMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sid := 123
	t.Run("Invalid certificate", func(t *testing.T) {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight())
		cert := certificate.NewBlockCertificate(0, 0)
		d, _ := blk.Bytes()
		msg := NewBlocksResponseMessage(ResponseCodeMoreBlocks, ResponseCodeMoreBlocks.String(),
			sid, ts.RandHeight(), [][]byte{d}, cert)
		err := msg.BasicCheck()

		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "height is not positive: 0",
		})
	})

	t.Run("OK", func(t *testing.T) {
		height := ts.RandHeight()
		blk1, _ := ts.GenerateTestBlock(height)
		blk2, cert2 := ts.GenerateTestBlock(height + 1)
		d1, _ := blk1.Bytes()
		d2, _ := blk2.Bytes()
		msg := NewBlocksResponseMessage(ResponseCodeMoreBlocks, ResponseCodeMoreBlocks.String(),
			sid, 100, [][]byte{d1, d2}, cert2)

		assert.NoError(t, msg.BasicCheck())
		assert.Contains(t, msg.String(), "100")
		assert.Equal(t, ResponseCodeMoreBlocks.String(), msg.Reason)
	})
}

func TestLatestBlocksResponseCode(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("rejected", func(t *testing.T) {
		reason := ts.RandString(16)
		msg := NewBlocksResponseMessage(ResponseCodeRejected, reason, 1, 0, nil, nil)

		assert.NoError(t, msg.BasicCheck())
		assert.Zero(t, msg.From)
		assert.Zero(t, msg.To())
		assert.Zero(t, msg.Count())
		assert.True(t, msg.IsRequestRejected())
		assert.Equal(t, reason, msg.Reason)
	})

	t.Run("OK - MoreBlocks", func(t *testing.T) {
		height := ts.RandHeight()
		blk1, _ := ts.GenerateTestBlock(height)
		blk2, _ := ts.GenerateTestBlock(height + 1)
		d1, _ := blk1.Bytes()
		d2, _ := blk2.Bytes()
		reason := ts.RandString(16)
		msg := NewBlocksResponseMessage(ResponseCodeMoreBlocks, reason, 1, 100, [][]byte{d1, d2}, nil)

		assert.NoError(t, msg.BasicCheck())
		assert.Equal(t, uint32(100), msg.From)
		assert.Equal(t, uint32(101), msg.To())
		assert.Equal(t, uint32(2), msg.Count())
		assert.False(t, msg.IsRequestRejected())
		assert.Equal(t, reason, msg.Reason)
	})

	t.Run("OK - Synced", func(t *testing.T) {
		height := ts.RandHeight()
		_, cert := ts.GenerateTestBlock(height)

		reason := ts.RandString(16)
		msg := NewBlocksResponseMessage(ResponseCodeSynced, reason, 1, 100, nil, cert)

		assert.NoError(t, msg.BasicCheck())
		assert.Equal(t, uint32(100), msg.From)
		assert.Zero(t, msg.To())
		assert.Zero(t, msg.Count())
		assert.False(t, msg.IsRequestRejected())
		assert.Equal(t, reason, msg.Reason)
	})
}
