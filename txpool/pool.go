package txpool

import (
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/linkedmap"
)

type payloadPool struct {
	list     *linkedmap.LinkedMap[tx.ID, *tx.Tx]
	minValue amount.Amount
}

func newPayloadPool(maxSize int, minValue amount.Amount) payloadPool {
	return payloadPool{
		list:     linkedmap.New[tx.ID, *tx.Tx](maxSize),
		minValue: minValue,
	}
}
