package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
)

func TestHelloType(t *testing.T) {
	p := &HelloPayload{}
	assert.Equal(t, p.Type(), PayloadTypeHello)
}

func TestHelloPayload(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		p := NewHelloPayload(util.RandomPeerID(), "Oscar", -1, 0, hash.GenerateTestHash())
		signer.SignMsg(p)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid signature", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		p := NewHelloPayload(util.RandomPeerID(), "Oscar", 100, 0, hash.GenerateTestHash())
		signer.SignMsg(p)

		p.PeerID = util.RandomPeerID()
		assert.Error(t, p.SanityCheck())
	})

	t.Run("Ok", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		p := NewHelloPayload(util.RandomPeerID(), "Alice", 100, 0, hash.GenerateTestHash())
		signer.SignMsg(p)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "Alice")
	})
}
