package param

import "time"

type Params struct {
	BlockTimeInSecond          int     `cbor:"1,keyasint"`
	CommitteeSize              int     `cbor:"2,keyasint"`
	SubsidyReductionInterval   int     `cbor:"3,keyasint"`
	TransactionToLiveInterval  int     `cbor:"4,keyasint"`
	UnbondInterval             int     `cbor:"5,keyasint"`
	MaximumTransactionPerBlock int     `cbor:"6,keyasint"`
	MaximumMemoLength          int     `cbor:"7,keyasint"`
	FeeFraction                float64 `cbor:"8,keyasint"`
	MinimumFee                 int64   `cbor:"9,keyasint"`
}

func MainnetParams() Params {
	return Params{
		BlockTimeInSecond:          10,
		CommitteeSize:              21,
		SubsidyReductionInterval:   2100000, // 243 days
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
