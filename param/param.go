package param

import "time"

type Params struct {
	BlockTimeInSecond          int     `cbor:"1,keyasint"`
	MaximumTransactionPerBlock int     `cbor:"2,keyasint"`
	MaximumPower               int     `cbor:"3,keyasint"`
	SubsidyReductionInterval   int     `cbor:"4,keyasint"`
	MaximumMemoLength          int     `cbor:"5,keyasint"`
	FeeFraction                float64 `cbor:"6,keyasint"`
	MinimumFee                 int64   `cbor:"7,keyasint"`
	TransactionToLiveInterval  int     `cbor:"8,keyasint"`
}

func MainnetParams() Params {
	return Params{
		BlockTimeInSecond:          10,
		MaximumTransactionPerBlock: 1000,
		MaximumPower:               21,
		SubsidyReductionInterval:   2100000,
		MaximumMemoLength:          1024,
		FeeFraction:                0.001,
		MinimumFee:                 1000,
		TransactionToLiveInterval:  500,
	}
}

func (p *Params) BlockTime() time.Duration {
	return time.Duration(p.BlockTimeInSecond) * time.Second
}
