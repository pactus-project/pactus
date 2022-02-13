package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/util"
)

func TestAleykType(t *testing.T) {
	p := &AleykPayload{}
	assert.Equal(t, p.Type(), PayloadTypeAleyk)
}

func TestAleykPayload(t *testing.T) {
	t.Run("Invalid target", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		p := NewAleykPayload(util.RandomPeerID(), "Oscar", 100, 0, "", ResponseCodeRejected, "rejected")
		signer.SignMsg(p)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid height", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		p := NewAleykPayload(util.RandomPeerID(), "Oscar", -1, 0, util.RandomPeerID(), ResponseCodeRejected, "rejected")
		signer.SignMsg(p)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("Invalid signature", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		p := NewAleykPayload(util.RandomPeerID(), "Oscar", 100, 0, util.RandomPeerID(), ResponseCodeRejected, "rejected")
		signer.SignMsg(p)

		p.PeerID = util.RandomPeerID()
		assert.Error(t, p.SanityCheck())
	})

	t.Run("Ok", func(t *testing.T) {
		signer := bls.GenerateTestSigner()
		p := NewAleykPayload(util.RandomPeerID(), "Alice", 100, 0, util.RandomPeerID(), ResponseCodeRejected, "welcome")
		signer.SignMsg(p)

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), "Alice")
	})
}
