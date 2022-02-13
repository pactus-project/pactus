package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
)

func TestSalamType(t *testing.T) {
	p := &SalamPayload{}
	assert.Equal(t, p.Type(), PayloadTypeSalam)
}

func TestSalamPayload(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		p := NewSalamPayload(util.RandomPeerID(), "Oscar", -1, 0, hash.GenerateTestHash())
		signer.SignMsg(p)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid signature", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		p := NewSalamPayload(util.RandomPeerID(), "Oscar", 100, 0, hash.GenerateTestHash())
		signer.SignMsg(p)

		p.PeerID = util.RandomPeerID()
		assert.Error(t, p.SanityCheck())
	})

	t.Run("Ok", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		p := NewSalamPayload(util.RandomPeerID(), "Alice", 100, 0, hash.GenerateTestHash())
		signer.SignMsg(p)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "Alice")
	})
}
