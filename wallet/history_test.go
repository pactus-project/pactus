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

	history.addTransaction("addr-1", Activity{TxID: "id-1-2", BlockTime: 2}, Transaction{})
	history.addTransaction("addr-1", Activity{TxID: "id-1-4", BlockTime: 4}, Transaction{})
	history.addTransaction("addr-1", Activity{TxID: "id-1-3", BlockTime: 3}, Transaction{})
	history.addTransaction("addr-1", Activity{TxID: "id-1-1", BlockTime: 1}, Transaction{})
	history.addTransaction("addr-2", Activity{TxID: "id-2-1", BlockTime: 6}, Transaction{})

	assert.Equal(t, history.Activities["addr-1"], []Activity{
		{TxID: "id-1-1", Status: "", BlockTime: 1, PayloadType: "", Amount: 0},
		{TxID: "id-1-2", Status: "", BlockTime: 2, PayloadType: "", Amount: 0},
		{TxID: "id-1-3", Status: "", BlockTime: 3, PayloadType: "", Amount: 0},
		{TxID: "id-1-4", Status: "", BlockTime: 4, PayloadType: "", Amount: 0}})
}
