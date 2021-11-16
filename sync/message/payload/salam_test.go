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
		p := NewSalamPayload("Eve", bls.GenerateTestSigner().PublicKey(), hash.GenerateTestHash(), -1, 0)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Ok", func(t *testing.T) {
		p := NewSalamPayload("Alice", bls.GenerateTestSigner().PublicKey(), hash.GenerateTestHash(), 0, 0)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "Alice")
	})
}
