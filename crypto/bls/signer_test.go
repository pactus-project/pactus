package bls

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

type testSignableMsg struct {
	sig *Signature
	pub *PublicKey
}

func (t *testSignableMsg) SignBytes() []byte {
	return []byte("zarb")
}
func (t *testSignableMsg) SetSignature(sig crypto.Signature) {
	t.sig = sig.(*Signature)
}
func (t *testSignableMsg) SetPublicKey(pub crypto.PublicKey) {
	t.pub = pub.(*PublicKey)
}

func TestSignable(t *testing.T) {
	signable := new(testSignableMsg)
	s := GenerateTestSigner()
	s.SignMsg(signable)

	assert.True(t, s.Address().EqualsTo(s.PublicKey().Address()))
	assert.True(t, signable.pub.EqualsTo(s.PublicKey()))
	assert.NoError(t, signable.pub.Verify(signable.SignBytes(), signable.sig))

	assert.True(t, s.SignData([]byte("zarb")).EqualsTo(signable.sig))

}
