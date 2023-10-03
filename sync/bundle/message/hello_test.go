package message

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
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
		valKey := ts.RandValKey()
		m := NewHelloMessage(ts.RandPeerID(), "Oscar", 100, 0, ts.RandHash(), ts.RandHash())
		m.Sign([]*bls.ValidatorKey{valKey})
		m.Signature = ts.RandBLSSignature()

		assert.ErrorIs(t, crypto.ErrInvalidSignature, m.BasicCheck())
	})

	t.Run("Signature is nil", func(t *testing.T) {
		valKey := ts.RandValKey()
		m := NewHelloMessage(ts.RandPeerID(), "Oscar", 100, 0, ts.RandHash(), ts.RandHash())
		m.Sign([]*bls.ValidatorKey{valKey})
		m.Signature = nil

		assert.Equal(t, errors.Code(m.BasicCheck()), errors.ErrInvalidSignature)
	})

	t.Run("PublicKeys are empty", func(t *testing.T) {
		valKey := ts.RandValKey()
		m := NewHelloMessage(ts.RandPeerID(), "Oscar", 100, 0, ts.RandHash(), ts.RandHash())
		m.Sign([]*bls.ValidatorKey{valKey})
		m.PublicKeys = make([]*bls.PublicKey, 0)

		assert.Equal(t, errors.Code(m.BasicCheck()), errors.ErrInvalidPublicKey)
	})

	t.Run("Ok", func(t *testing.T) {
		valKey := ts.RandValKey()
		m := NewHelloMessage(ts.RandPeerID(), "Alice", 100, 0, ts.RandHash(), ts.RandHash())
		m.Sign([]*bls.ValidatorKey{valKey})

		assert.NoError(t, m.BasicCheck())
		assert.Contains(t, m.String(), "Alice")
	})
}
