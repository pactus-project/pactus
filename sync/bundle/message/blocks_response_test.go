package message

import (
	"testing"

	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestLatestBlocksResponseType(t *testing.T) {
	m := &BlocksResponseMessage{}
	assert.Equal(t, m.Type(), TypeBlocksResponse)
}

func TestBlocksResponseMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	sid := 123
	t.Run("Invalid certificate", func(t *testing.T) {
		blk, _ := ts.GenerateTestBlock(ts.RandHeight())
		cert := certificate.NewBlockCertificate(0, 0, false)
		d, _ := blk.Bytes()
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, ResponseCodeMoreBlocks.String(),
			sid, ts.RandHeight(), [][]byte{d}, cert)
		err := m.BasicCheck()

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
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, ResponseCodeMoreBlocks.String(),
			sid, 100, [][]byte{d1, d2}, cert2)

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), "100")
		assert.Equal(t, m.Reason, ResponseCodeMoreBlocks.String())
	})
}

func TestLatestBlocksResponseCode(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("rejected", func(t *testing.T) {
		reason := ts.RandString(16)
		m := NewBlocksResponseMessage(ResponseCodeRejected, reason, 1, 0, nil, nil)

		assert.NoError(t, m.BasicCheck())
		assert.Zero(t, m.From)
		assert.Zero(t, m.To())
		assert.Zero(t, m.Count())
		assert.True(t, m.IsRequestRejected())
		assert.Equal(t, m.Reason, reason)
	})

	t.Run("OK - MoreBlocks", func(t *testing.T) {
		height := ts.RandHeight()
		blk1, _ := ts.GenerateTestBlock(height)
		blk2, _ := ts.GenerateTestBlock(height + 1)
		d1, _ := blk1.Bytes()
		d2, _ := blk2.Bytes()
		reason := ts.RandString(16)
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, reason, 1, 100, [][]byte{d1, d2}, nil)

		assert.NoError(t, m.BasicCheck())
		assert.Equal(t, m.From, uint32(100))
		assert.Equal(t, m.To(), uint32(101))
		assert.Equal(t, m.Count(), uint32(2))
		assert.False(t, m.IsRequestRejected())
		assert.Equal(t, m.Reason, reason)
	})

	t.Run("OK - Synced", func(t *testing.T) {
		height := ts.RandHeight()
		_, cert := ts.GenerateTestBlock(height)

		reason := ts.RandString(16)
		m := NewBlocksResponseMessage(ResponseCodeSynced, reason, 1, 100, nil, cert)

		assert.NoError(t, m.BasicCheck())
		assert.Equal(t, m.From, uint32(100))
		assert.Zero(t, m.To())
		assert.Zero(t, m.Count())
		assert.False(t, m.IsRequestRejected())
		assert.Equal(t, m.Reason, reason)
	})
}
