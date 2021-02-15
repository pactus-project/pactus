package param

import "time"

type Params struct {
	BlockVersion               int     `cbor:"1,keyasint"`
	BlockTimeInSecond          int     `cbor:"2,keyasint"`
	CommitteeSize              int     `cbor:"3,keyasint"`
	SubsidyReductionInterval   int     `cbor:"4,keyasint"`
	TransactionToLiveInterval  int     `cbor:"5,keyasint"`
	UnbondInterval             int     `cbor:"6,keyasint"`
	MaximumTransactionPerBlock int     `cbor:"7,keyasint"`
	MaximumMemoLength          int     `cbor:"8,keyasint"`
	FeeFraction                float64 `cbor:"9,keyasint"`
	MinimumFee                 int64   `cbor:"10,keyasint"`
}

func DefaultParams() Params {
	return Params{
		BlockVersion:               1,
		BlockTimeInSecond:          10,
		CommitteeSize:              21,
		SubsidyReductionInterval:   4200000, // 16 months
		TransactionToLiveInterval:  8640,    // one days
		UnbondInterval:             181440,  // 21 days
		MaximumTransactionPerBlock: 1000,
		MaximumMemoLength:          1024,
		FeeFraction:                0.001,
		MinimumFee:                 1000,
	}
}

func (p *Params) BlockTime() time.Duration {
	return time.Duration(p.BlockTimeInSecond) * time.Second
}
