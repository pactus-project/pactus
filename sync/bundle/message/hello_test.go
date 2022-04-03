package message

import (
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
)

func TestHelloType(t *testing.T) {
	m := &HelloMessage{}
	assert.Equal(t, m.Type(), MessageTypeHello)
}

func TestHelloMessage(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		m := NewHelloMessage(peer.ID("oscar-peer-id"), "Oscar", -1, 0, hash.GenerateTestHash())
		signer.SignMsg(m)

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidHeight)
	})

	t.Run("Invalid signature", func(t *testing.T) {
		signer1 := bls.GenerateTestSigner()
		signer2 := bls.GenerateTestSigner()
		m := NewHelloMessage(peer.ID("oscar-peer-id"), "Oscar", 100, 0, hash.GenerateTestHash())
		signer1.SignMsg(m)
		m.SetPublicKey(signer2.PublicKey())

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidSignature)
	})

	t.Run("Ok", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		m := NewHelloMessage(peer.ID("alice-peer-id"), "Alice", 100, 0, hash.GenerateTestHash())
		signer.SignMsg(m)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "Alice")
	})
}
