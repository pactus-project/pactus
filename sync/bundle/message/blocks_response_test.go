package message

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestLatestBlocksResponseType(t *testing.T) {
	m := &BlocksResponseMessage{}
	assert.Equal(t, m.Type(), MessageTypeBlocksResponse)
}

func TestBlocksResponseMessage(t *testing.T) {
	t.Run("Invalid certificate", func(t *testing.T) {
		b := block.GenerateTestBlock(nil, nil)
		c := block.NewCertificate(-1, nil, nil, nil)
		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, 1, 100, []*block.Block{b}, c)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidRound)
	})

	t.Run("OK", func(t *testing.T) {
		b1 := block.GenerateTestBlock(nil, nil)
		b2 := block.GenerateTestBlock(nil, nil)
		m := NewBlocksResponseMessage(ResponseCodeSynced, 1, 100, []*block.Block{b1, b2}, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Zero(t, m.LastCertificateHeight())
		assert.Contains(t, m.Fingerprint(), "100")
	})
}

func TestLatestBlocksResponseCode(t *testing.T) {
	t.Run("busy", func(t *testing.T) {
		m := NewBlocksResponseMessage(ResponseCodeBusy, 1, 0, nil, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Zero(t, m.From)
		assert.Zero(t, m.To())
		assert.Zero(t, m.Count())
		assert.True(t, m.IsRequestRejected())
	})

	t.Run("rejected", func(t *testing.T) {
		m := NewBlocksResponseMessage(ResponseCodeRejected, 1, 0, nil, nil)

		assert.NoError(t, m.SanityCheck())
		assert.Zero(t, m.From)
		assert.Zero(t, m.To())
		assert.Zero(t, m.Count())
		assert.True(t, m.IsRequestRejected())
	})

	t.Run("OK - MoreBlocks", func(t *testing.T) {
		b1 := block.GenerateTestBlock(nil, nil)
		b2 := block.GenerateTestBlock(nil, nil)

		m := NewBlocksResponseMessage(ResponseCodeMoreBlocks, 1, 100, []*block.Block{b1, b2}, nil)
		assert.NoError(t, m.SanityCheck())
		assert.Equal(t, m.From, uint32(100))
		assert.Equal(t, m.To(), uint32(101))
		assert.Equal(t, m.Count(), uint32(2))
		assert.Zero(t, m.LastCertificateHeight())
		assert.False(t, m.IsRequestRejected())
	})

	t.Run("OK - Synced", func(t *testing.T) {
		cert := block.GenerateTestCertificate(hash.GenerateTestHash())

		m := NewBlocksResponseMessage(ResponseCodeSynced, 1, 100, nil, cert)
		assert.NoError(t, m.SanityCheck())
		assert.Equal(t, m.From, uint32(100))
		assert.Zero(t, m.To())
		assert.Zero(t, m.Count())
		assert.Equal(t, m.LastCertificateHeight(), uint32(100))
		assert.False(t, m.IsRequestRejected())
	})
}
