package message

import (
	"testing"

	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestHelloType(t *testing.T) {
	m := &HelloMessage{}
	assert.Equal(t, m.Type(), TypeHello)
}

func TestHelloMessage(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid signature", func(t *testing.T) {
		signer1 := ts.RandomSigner()
		signer2 := ts.RandomSigner()
		m := NewHelloMessage(ts.RandomPeerID(), "Oscar", 100, 0, ts.RandomHash(), ts.RandomHash())
		signer1.SignMsg(m)
		m.SetPublicKey(signer2.PublicKey())

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidSignature)
	})

	t.Run("Signature is nil", func(t *testing.T) {
		signer := ts.RandomSigner()
		m := NewHelloMessage(ts.RandomPeerID(), "Oscar", 100, 0, ts.RandomHash(), ts.RandomHash())
		signer.SignMsg(m)
		m.Signature = nil

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidSignature)
	})

	t.Run("PublicKey is nil", func(t *testing.T) {
		signer := ts.RandomSigner()
		m := NewHelloMessage(ts.RandomPeerID(), "Oscar", 100, 0, ts.RandomHash(), ts.RandomHash())
		signer.SignMsg(m)
		m.PublicKey = nil

		assert.Equal(t, errors.Code(m.SanityCheck()), errors.ErrInvalidPublicKey)
	})

	t.Run("Ok", func(t *testing.T) {
		signer := ts.RandomSigner()
		m := NewHelloMessage(ts.RandomPeerID(), "Alice", 100, 0, ts.RandomHash(), ts.RandomHash())
		signer.SignMsg(m)

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), "Alice")
	})
}
