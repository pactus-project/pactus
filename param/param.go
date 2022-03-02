package param

import "time"

type Params struct {
	BlockVersion               int     `cbor:"1,keyasint"`
	BlockTimeInSecond          int     `cbor:"2,keyasint"`
	CommitteeSize              int     `cbor:"3,keyasint"`
	BlockReward                int64   `cbor:"4,keyasint"`
	TransactionToLiveInterval  int     `cbor:"5,keyasint"`
	BondInterval               int     `cbor:"6,keyasint"`
	UnbondInterval             int     `cbor:"7,keyasint"`
	MaximumTransactionPerBlock int     `cbor:"8,keyasint"`
	MaximumMemoLength          int     `cbor:"9,keyasint"`
	FeeFraction                float64 `cbor:"10,keyasint"`
	MinimumFee                 int64   `cbor:"11,keyasint"`
}

func DefaultParams() Params {
	return Params{
		BlockVersion:               1,
		BlockTimeInSecond:          10,
		CommitteeSize:              21,
		BlockReward:                100000000,
		TransactionToLiveInterval:  8640,   // one day
		BondInterval:               360,    // one hour
		UnbondInterval:             181440, // 21 days
		MaximumTransactionPerBlock: 1000,
		MaximumMemoLength:          1024,
		FeeFraction:                0.001,
		MinimumFee:                 1000,
	}
}

func (p *Params) BlockTime() time.Duration {
	return time.Duration(p.BlockTimeInSecond) * time.Second
}
