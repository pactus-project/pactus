package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
)

func TestBlockAnnounceType(t *testing.T) {
	m := &BlockAnnounceMessage{}
	assert.Equal(t, m.Type(), MessageTypeBlockAnnounce)
}

func TestBlockAnnounceMessage(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		b := block.GenerateTestBlock(nil, nil)
		c := block.GenerateTestCertificate(b.Hash())
		m := NewBlockAnnounceMessage(-1, b, c)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("Invalid certificate", func(t *testing.T) {
		b := block.GenerateTestBlock(nil, nil)
		c := block.NewCertificate(-1, nil, nil, nil)
		m := NewBlockAnnounceMessage(100, b, c)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		b := block.GenerateTestBlock(nil, nil)
		c := block.GenerateTestCertificate(b.Hash())
		m := NewBlockAnnounceMessage(100, b, c)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "100")
	})
}
