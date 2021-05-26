package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

func TestQueryTransactionsType(t *testing.T) {
	p := &QueryTransactionsPayload{}
	assert.Equal(t, p.Type(), PayloadTypeQueryTransactions)
}

func TestQueryTransactionsPayload(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		p := NewQueryTransactionsPayload(nil)

		assert.Error(t, p.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		id := crypto.GenerateTestHash()
		p := NewQueryTransactionsPayload([]tx.ID{id})

		assert.NoError(t, p.SanityCheck())
		assert.Contains(t, p.Fingerprint(), id.Fingerprint())
	})
}
