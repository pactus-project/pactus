package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestBlockAnnounceType(t *testing.T) {
	p := &BlockAnnouncePayload{}
	assert.Equal(t, p.Type(), PayloadTypeBlockAnnounce)
}

func TestBlockAnnouncePayload(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		b, _ := block.GenerateTestBlock(nil, nil)
		c := block.GenerateTestCertificate(b.Hash())
		p := NewBlockAnnouncePayload(-1, b, c)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid certificate", func(t *testing.T) {
		b, _ := block.GenerateTestBlock(nil, nil)
		c := block.GenerateTestCertificate(crypto.UndefHash)
		p := NewBlockAnnouncePayload(100, b, c)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		b, _ := block.GenerateTestBlock(nil, nil)
		c := block.GenerateTestCertificate(b.Hash())
		p := NewBlockAnnouncePayload(100, b, c)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "100")
	})
}
