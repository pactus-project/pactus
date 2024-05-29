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
	capacity := p.list.Capacity()
	size := p.list.Size()
	usageRatio := float64(size) / float64(capacity)

	switch {
	case usageRatio > 0.90:
		return p.minFee * 1000000
	case usageRatio > 0.80:
		return p.minFee * 100000
	case usageRatio > 0.70:
		return p.minFee * 10000
	case usageRatio > 0.60:
		return p.minFee * 1000
	case usageRatio > 0.50:
		return p.minFee * 100
	case usageRatio > 0.40:
		return p.minFee * 10
	default:
		return p.minFee
	}
}
