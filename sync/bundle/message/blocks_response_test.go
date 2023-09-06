package message

import (
	"testing"

	"github.com/pactus-project/pactus/types/certificate"
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
		b := ts.GenerateTestBlock(nil)
		c := certificate.NewCertificate(0, 0, nil, nil, nil)
		d, _ := b.Bytes()
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, ResponseCodeMoreBlocks.String(), sid, 100, [][]byte{d}, c)
		err := m.BasicCheck()

		assert.ErrorIs(t, err, certificate.BasicCheckError{
			Reason: "height is not positive: 0",
		})
	})

	t.Run("Unexpected block for height zero", func(t *testing.T) {
		b := ts.GenerateTestBlock(nil)
		d, _ := b.Bytes()
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, ResponseCodeMoreBlocks.String(), sid, 0, [][]byte{d}, nil)

		assert.Equal(t, errors.Code(m.BasicCheck()), errors.ErrInvalidHeight)
	})

	t.Run("OK", func(t *testing.T) {
		b1 := ts.GenerateTestBlock(nil)
		b2 := ts.GenerateTestBlock(nil)
		d1, _ := b1.Bytes()
		d2, _ := b2.Bytes()
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, ResponseCodeMoreBlocks.String(), sid, 100,
			[][]byte{d1, d2}, nil)

		assert.NoError(t, m.BasicCheck())
		assert.Zero(t, m.LastCertificateHeight())
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
		b1 := ts.GenerateTestBlock(nil)
		b2 := ts.GenerateTestBlock(nil)
		d1, _ := b1.Bytes()
		d2, _ := b2.Bytes()
		reason := ts.RandString(16)
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, reason, 1, 100, [][]byte{d1, d2}, nil)

		assert.NoError(t, m.BasicCheck())
		assert.Equal(t, m.From, uint32(100))
		assert.Equal(t, m.To(), uint32(101))
		assert.Equal(t, m.Count(), uint32(2))
		assert.Zero(t, m.LastCertificateHeight())
		assert.False(t, m.IsRequestRejected())
		assert.Equal(t, m.Reason, reason)
	})

	t.Run("OK - Synced", func(t *testing.T) {
		cert := ts.GenerateTestCertificate()

		reason := ts.RandString(16)
		m := NewBlocksResponseMessage(ResponseCodeSynced, reason, 1, 100, nil, cert)

		assert.NoError(t, m.BasicCheck())
		assert.Equal(t, m.From, uint32(100))
		assert.Zero(t, m.To())
		assert.Zero(t, m.Count())
		assert.Equal(t, m.LastCertificateHeight(), uint32(100))
		assert.False(t, m.IsRequestRejected())
		assert.Equal(t, m.Reason, reason)
	})
}
