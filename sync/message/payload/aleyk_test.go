package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

func TestAleykType(t *testing.T) {
	p := &AleykPayload{}
	assert.Equal(t, p.Type(), PayloadTypeAleyk)
}

func TestAleykPayload(t *testing.T) {
	t.Run("Invalid target", func(t *testing.T) {
		p := NewAleykPayload("", ResponseCodeRejected, "rejected",
			"Eve", crypto.GenerateTestSigner().PublicKey(), 100, 0)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid height", func(t *testing.T) {
		p := NewAleykPayload(util.RandomPeerID(), ResponseCodeRejected, "rejected",
			"Eve", crypto.GenerateTestSigner().PublicKey(), -1, 0)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Ok", func(t *testing.T) {
		p := NewAleykPayload(util.RandomPeerID(), ResponseCodeRejected, "welcome",
			"Alice", crypto.GenerateTestSigner().PublicKey(), 100, 0)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "Alice")
	})
}
