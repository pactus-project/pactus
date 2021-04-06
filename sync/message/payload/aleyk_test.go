package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestAleykType(t *testing.T) {
	p := &AleykPayload{}
	assert.Equal(t, p.Type(), PayloadTypeAleyk)
}

func TestAleykPayload(t *testing.T) {
	p := NewAleykPayload(ResponseCodeRejected, "busy",
		"devil", crypto.GenerateTestSigner().PublicKey(), -1, 0)

	assert.Error(t, p.SanityCheck())
}
