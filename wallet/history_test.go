package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddTransaction(t *testing.T) {
	history := history{
		Activities:   map[string][]Activity{},
		Transactions: map[string]Transaction{},
	}

	history.addTransaction("addr-1", "id-1-2", Activity{}, Transaction{BlockTime: 2})
	history.addTransaction("addr-1", "id-1-4", Activity{}, Transaction{BlockTime: 4})
	history.addTransaction("addr-1", "id-1-3", Activity{}, Transaction{BlockTime: 3})
	history.addTransaction("addr-1", "id-1-1", Activity{}, Transaction{BlockTime: 1})
	history.addTransaction("addr-2", "id-2-1", Activity{}, Transaction{BlockTime: 6})

	assert.Equal(t, history.Activities["addr-1"], []string{"id-1-1", "id-1-2", "id-1-3", "id-1-4"})
}
