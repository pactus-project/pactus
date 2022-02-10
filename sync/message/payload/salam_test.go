package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func TestSalamType(t *testing.T) {
	p := &SalamPayload{}
	assert.Equal(t, p.Type(), PayloadTypeSalam)
}

func TestSalamPayload(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		sig := signer.SignData(signer.PublicKey().RawBytes())
		p := NewSalamPayload("Oscar", signer.PublicKey(), sig, -1, 0, hash.GenerateTestHash())

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid signature", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		sig := signer.SignData(nil)
		p := NewSalamPayload("Oscar", signer.PublicKey(), sig, -1, 0, hash.GenerateTestHash())

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Ok", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		sig := signer.SignData(signer.PublicKey().RawBytes())
		p := NewSalamPayload("Alice", signer.PublicKey(), sig, 0, 0, hash.GenerateTestHash())

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "Alice")
	})
}
