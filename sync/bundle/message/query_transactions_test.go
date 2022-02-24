package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/tx"
)

func TestQueryTransactionsType(t *testing.T) {
	m := &QueryTransactionsMessage{}
	assert.Equal(t, m.Type(), MessageTypeQueryTransactions)
}

func TestQueryTransactionsMessage(t *testing.T) {
	t.Run("Invalid height", func(t *testing.T) {
		m := NewQueryTransactionsMessage(nil)

		assert.Error(t, m.SanityCheck())
	})

	t.Run("OK", func(t *testing.T) {
		id := hash.GenerateTestHash()
		m := NewQueryTransactionsMessage([]tx.ID{id})

		assert.NoError(t, m.SanityCheck())
		assert.Contains(t, m.Fingerprint(), id.Fingerprint())
	})
}
