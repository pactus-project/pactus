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
		p1 := NewQueryTransactionsPayload(nil)
		assert.Error(t, p1.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		p2 := NewQueryTransactionsPayload([]tx.ID{crypto.GenerateTestHash()})
		assert.NoError(t, p2.SanityCheck())
	})
}
