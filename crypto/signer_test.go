package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSignable struct {
	sig *Signature
}

func (t *testSignable) SignBytes() []byte {
	return []byte("zarb")
}
func (t *testSignable) SetSignature(sig *Signature) {
	t.sig = sig
}

func TestSigner(t *testing.T) {
	_, _, priv := GenerateTestKeyPair()
	s := NewSigner(priv)
	pub := s.PublicKey()
	assert.True(t, pub.EqualsTo(priv.PublicKey()))
	assert.True(t, s.Address().EqualsTo(priv.PublicKey().Address()))
	assert.True(t, s.Sign([]byte("zarb")).EqualsTo(*priv.Sign([]byte("zarb"))))
}

func TestSignable(t *testing.T) {
	signable := new(testSignable)
	_, pub, priv := GenerateTestKeyPair()
	s := NewSigner(priv)
	s.SignMsg(signable)

	assert.True(t, pub.Verify(signable.SignBytes(), signable.sig))
}
