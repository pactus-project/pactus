package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
)

func TestHelloType(t *testing.T) {
	m := &HelloMessage{}
	assert.Equal(t, m.Type(), MessageTypeHello)
}

func TestHelloMessage(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		m := NewHelloMessage(util.RandomPeerID(), "Oscar", -1, 0, hash.GenerateTestHash())
		signer.SignMsg(m)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("Invalid signature", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		m := NewHelloMessage(util.RandomPeerID(), "Oscar", 100, 0, hash.GenerateTestHash())
		signer.SignMsg(m)

		m.PeerID = util.RandomPeerID()
		assert.Error(t, m.SanityCheck())
	})

	t.Run("Ok", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		m := NewHelloMessage(util.RandomPeerID(), "Alice", 100, 0, hash.GenerateTestHash())
		signer.SignMsg(m)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "Alice")
	})
}
