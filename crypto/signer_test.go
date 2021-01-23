package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSignableMsg struct {
	sig *Signature
	pub *PublicKey
}

func (t *testSignableMsg) SignBytes() []byte {
	return []byte("zarb")
}
func (t *testSignableMsg) SetSignature(sig Signature) {
	t.sig = &sig
}
func (t *testSignableMsg) SetPublicKey(pub PublicKey) {
	t.pub = &pub
}

func TestSignable(t *testing.T) {
	signable := new(testSignableMsg)
	s := GenerateTestSigner()
	s.SignMsg(signable)

	assert.True(t, s.Address().EqualsTo(s.PublicKey().Address()))
	assert.True(t, signable.pub.EqualsTo(s.PublicKey()))
	assert.True(t, signable.pub.Verify(signable.SignBytes(), *signable.sig))

	assert.True(t, s.SignData([]byte("zarb")).EqualsTo(*signable.sig))

}
