package message

import (
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/errors"
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
		b := ts.GenerateTestBlock(nil, nil)
		c := block.NewCertificate(-1, nil, nil, nil)
		d, _ := b.Bytes()
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, ResponseCodeMoreBlocks.String(), sid, 100, [][]byte{d}, c)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidRound)
		assert.Equal(t, m.Reason, ResponseCodeMoreBlocks.String())
	})

	t.Run("Unexpected block for height zero", func(t *testing.T) {
		b := ts.GenerateTestBlock(nil, nil)
		d, _ := b.Bytes()
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, ResponseCodeMoreBlocks.String(), sid, 0, [][]byte{d}, nil)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidHeight)
		assert.Equal(t, m.Reason, ResponseCodeMoreBlocks.String())
	})

	t.Run("OK", func(t *testing.T) {
		b1 := ts.GenerateTestBlock(nil, nil)
		b2 := ts.GenerateTestBlock(nil, nil)
		d1, _ := b1.Bytes()
		d2, _ := b2.Bytes()
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, ResponseCodeMoreBlocks.String(), sid, 100,
			[][]byte{d1, d2}, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Zero(t, m.LastCertificateHeight())
		assert.Contains(t, m.String(), "100")
		assert.Equal(t, m.Reason, ResponseCodeMoreBlocks.String())
	})
}

func TestLatestBlocksResponseCode(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("busy", func(t *testing.T) {
		m := NewBlocksResponseMessage(ResponseCodeRejected, ResponseCodeRejected.String(), 1, 0, nil, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Zero(t, m.From)
		assert.Zero(t, m.To())
		assert.Zero(t, m.Count())
		assert.True(t, m.IsRequestRejected())
		assert.Equal(t, m.Reason, ResponseCodeRejected.String())
	})

	t.Run("rejected", func(t *testing.T) {
		m := NewBlocksResponseMessage(ResponseCodeRejected, ResponseCodeRejected.String(), 1, 0, nil, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Zero(t, m.From)
		assert.Zero(t, m.To())
		assert.Zero(t, m.Count())
		assert.True(t, m.IsRequestRejected())
		assert.Equal(t, m.Reason, ResponseCodeRejected.String())
	})

	t.Run("OK - MoreBlocks", func(t *testing.T) {
		b1 := ts.GenerateTestBlock(nil, nil)
		b2 := ts.GenerateTestBlock(nil, nil)
		d1, _ := b1.Bytes()
		d2, _ := b2.Bytes()
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, ResponseCodeMoreBlocks.String(), 1, 100, [][]byte{d1, d2}, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Equal(t, m.From, uint32(100))
		assert.Equal(t, m.To(), uint32(101))
		assert.Equal(t, m.Count(), uint32(2))
		assert.Zero(t, m.LastCertificateHeight())
		assert.False(t, m.IsRequestRejected())
		assert.Equal(t, m.Reason, ResponseCodeMoreBlocks.String())
	})

	t.Run("OK - Synced", func(t *testing.T) {
		cert := ts.GenerateTestCertificate(ts.RandomHash())

		m := NewBlocksResponseMessage(ResponseCodeSynced, ResponseCodeSynced.String(), 1, 100, nil, cert)
		assert.NoError(t, m.SanityCheck())
		assert.Equal(t, m.From, uint32(100))
		assert.Zero(t, m.To())
		assert.Zero(t, m.Count())
		assert.Equal(t, m.LastCertificateHeight(), uint32(100))
		assert.False(t, m.IsRequestRejected())
		assert.Equal(t, m.Reason, ResponseCodeSynced.String())
	})
}
