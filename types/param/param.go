package param

import (
	"time"

	"github.com/pactus-project/pactus/types/amount"
)

type Params struct {
	BlockVersion              uint8
	BlockIntervalInSecond     int
	CommitteeSize             int
	BlockReward               amount.Amount
	TransactionToLiveInterval uint32
	BondInterval              uint32
	UnbondInterval            uint32
	SortitionInterval         uint32
	MaxTransactionsPerBlock   int
	MinimumStake              amount.Amount
	MaximumStake              amount.Amount
}

func DefaultParams() *Params {
	return &Params{
		BlockVersion:              1,
		BlockIntervalInSecond:     10,
		CommitteeSize:             51,
		BlockReward:               1000000000,
		TransactionToLiveInterval: 8640,   // one day
		BondInterval:              360,    // one hour
		UnbondInterval:            181440, // 21 days
		SortitionInterval:         17,
		MaxTransactionsPerBlock:   1000,
		MinimumStake:              1000000000,
		MaximumStake:              1000000000000,
	}
}

func (p *Params) BlockInterval() time.Duration {
	return time.Duration(p.BlockIntervalInSecond) * time.Second
}
