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
	b, _ := block.GenerateTestBlock(nil, nil)
	c := block.GenerateTestCertificate(crypto.UndefHash)

	p1 := NewBlockAnnouncePayload(-1, b, c)
	assert.Error(t, p1.SanityCheck())

	p2 := NewBlockAnnouncePayload(100, b, c)
	assert.Error(t, p2.SanityCheck())

	c = block.GenerateTestCertificate(b.Hash())
	p3 := NewBlockAnnouncePayload(100, b, c)
	assert.NoError(t, p3.SanityCheck())
}
