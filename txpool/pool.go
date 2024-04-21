package txpool

import (
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/linkedmap"
)

type pool struct {
	list     *linkedmap.LinkedMap[tx.ID, *tx.Tx]
	minValue amount.Amount
}

func newPool(maxSize int, minValue amount.Amount) pool {
	return pool{
		list:     linkedmap.New[tx.ID, *tx.Tx](maxSize),
		minValue: minValue,
	}
}
