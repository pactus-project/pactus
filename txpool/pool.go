package txpool

import (
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/linkedmap"
)

type pool struct {
	list   *linkedmap.LinkedMap[tx.ID, *tx.Tx]
	minFee amount.Amount
}

func newPool(maxSize int, minFee amount.Amount) pool {
	return pool{
		list:   linkedmap.New[tx.ID, *tx.Tx](maxSize),
		minFee: minFee,
	}
}

func (p *pool) calculateDynamicFee() amount.Amount {
	totalSize := p.list.Capacity()
	currentSize := p.list.Size()
	usageRatio := float64(currentSize) / float64(totalSize)

	switch {
	case usageRatio > 0.90:
		return p.minFee * 10000
	case usageRatio > 0.80:
		return p.minFee * 1000
	case usageRatio > 0.70:
		return p.minFee * 100
	case usageRatio > 0.50:
		return p.minFee * 10
	default:
		return p.minFee
	}
}
