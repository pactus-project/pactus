package bls_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

type testSignableMsg struct {
	sig *bls.Signature
	pub *bls.PublicKey
}

func (t *testSignableMsg) SignBytes() []byte {
	return []byte("zarb")
}
func (t *testSignableMsg) SetSignature(sig crypto.Signature) {
	t.sig = sig.(*bls.Signature)
}
func (t *testSignableMsg) SetPublicKey(pub crypto.PublicKey) {
	t.pub = pub.(*bls.PublicKey)
}

func TestSignable(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	signable := new(testSignableMsg)
	s := ts.RandomSigner()
	s.SignMsg(signable)

	assert.True(t, s.Address().EqualsTo(s.PublicKey().Address()))
	assert.True(t, signable.pub.EqualsTo(s.PublicKey()))
	assert.NoError(t, signable.pub.Verify(signable.SignBytes(), signable.sig))

	assert.True(t, s.SignData([]byte("zarb")).EqualsTo(signable.sig))
}
