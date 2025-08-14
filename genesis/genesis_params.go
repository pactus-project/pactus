package genesis

import (
	"time"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
)

type GenesisParams struct {
	BlockVersion              protocol.Version `cbor:"1,keyasint"  json:"block_version"`
	BlockIntervalInSecond     int              `cbor:"2,keyasint"  json:"block_interval_in_second"`
	CommitteeSize             int              `cbor:"3,keyasint"  json:"committee_size"`
	BlockReward               amount.Amount    `cbor:"4,keyasint"  json:"block_reward"`
	TransactionToLiveInterval uint32           `cbor:"5,keyasint"  json:"transaction_to_live_interval"`
	BondInterval              uint32           `cbor:"6,keyasint"  json:"bond_interval"`
	UnbondInterval            uint32           `cbor:"7,keyasint"  json:"unbond_interval"`
	SortitionInterval         uint32           `cbor:"8,keyasint"  json:"sortition_interval"`
	FeeFraction               float64          `cbor:"9,keyasint"  json:"fee_fraction"` // Deprecated: Replaced by fix fee
	MinimumFee                amount.Amount    `cbor:"10,keyasint" json:"minimum_fee"`  // Deprecated: Replaced by fix fee
	MaximumFee                amount.Amount    `cbor:"11,keyasint" json:"maximum_fee"`  // Deprecated: Replaced by fix fee
	MinimumStake              amount.Amount    `cobr:"12,keyasint" json:"minimum_stake"`
	MaximumStake              amount.Amount    `cbor:"13,keyasint" json:"maximum_stake"`
}

func DefaultGenesisParams() *GenesisParams {
	return &GenesisParams{
		BlockVersion:              1,
		BlockIntervalInSecond:     10,
		CommitteeSize:             51,
		BlockReward:               1000000000,
		TransactionToLiveInterval: 8640,   // one day
		BondInterval:              360,    // one hour
		UnbondInterval:            181440, // 21 days
		SortitionInterval:         17,
		MinimumStake:              1000000000,
		MaximumStake:              1000000000000,

		// Deprecated: Replaced by fix fee
		FeeFraction: 0.0,
		MinimumFee:  0,
		MaximumFee:  0,
	}
}

func (p *GenesisParams) BlockInterval() time.Duration {
	return time.Duration(p.BlockIntervalInSecond) * time.Second
}
