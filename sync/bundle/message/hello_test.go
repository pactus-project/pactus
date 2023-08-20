package message

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"

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
		signer := ts.RandomSigner()
		m := NewHelloMessage(ts.RandomPeerID(), "Oscar", 100, 0, ts.RandomHash(), ts.RandomHash())
		m.Sign(signer)
		m.Signature = ts.RandomBLSSignature()

		assert.Equal(t, errors.Code(m.BasicCheck()), errors.ErrInvalidSignature)
	})

	t.Run("Signature is nil", func(t *testing.T) {
		signer := ts.RandomSigner()
		m := NewHelloMessage(ts.RandomPeerID(), "Oscar", 100, 0, ts.RandomHash(), ts.RandomHash())
		m.Sign(signer)
		m.Signature = nil

		assert.Equal(t, errors.Code(m.BasicCheck()), errors.ErrInvalidSignature)
	})

	t.Run("PublicKeys are empty", func(t *testing.T) {
		signer := ts.RandomSigner()
		m := NewHelloMessage(ts.RandomPeerID(), "Oscar", 100, 0, ts.RandomHash(), ts.RandomHash())
		m.Sign(signer)
		m.PublicKeys = make([]*bls.PublicKey, 0)

		assert.Equal(t, errors.Code(m.BasicCheck()), errors.ErrInvalidPublicKey)
	})

	t.Run("Ok", func(t *testing.T) {
		signer := ts.RandomSigner()
		m := NewHelloMessage(ts.RandomPeerID(), "Alice", 100, 0, ts.RandomHash(), ts.RandomHash())
		m.Sign(signer)

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), "Alice")
	})
}
