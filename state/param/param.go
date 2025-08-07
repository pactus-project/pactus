package param

import (
	"time"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
)

type Params struct {
	BlockVersion              uint8
	BlockIntervalInSecond     int
	MaxTransactionsPerBlock   int
	CommitteeSize             int
	BlockReward               amount.Amount
	TransactionToLiveInterval uint32
	BondInterval              uint32
	UnbondInterval            uint32
	SortitionInterval         uint32
	MinimumStake              amount.Amount
	MaximumStake              amount.Amount
	FoundationReward          amount.Amount
}

func FromGenesis(genDoc *genesis.GenesisParams) *Params {
	return &Params{
		// genesis parameters
		BlockVersion:              genDoc.BlockVersion,
		BlockIntervalInSecond:     genDoc.BlockIntervalInSecond,
		CommitteeSize:             genDoc.CommitteeSize,
		BlockReward:               genDoc.BlockReward,
		TransactionToLiveInterval: genDoc.TransactionToLiveInterval,
		BondInterval:              genDoc.BondInterval,
		UnbondInterval:            genDoc.UnbondInterval,
		SortitionInterval:         genDoc.SortitionInterval,
		MaximumStake:              genDoc.MaximumStake,
		MinimumStake:              genDoc.MinimumStake,

		// chain parameters
		MaxTransactionsPerBlock: 1000,
		FoundationReward:        amount.Amount(300_000_000),
	}
}

func (p *Params) BlockInterval() time.Duration {
	return time.Duration(p.BlockIntervalInSecond) * time.Second
}
