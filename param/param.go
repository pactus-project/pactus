package param

import "time"

type Params struct {
	BlockTime                  time.Duration
	MaximumTransactionPerBlock int
	MaximumPower               int
	SubsidyReductionInterval   int
	MaximumMemoLength          int
	FeeFraction                float64
	MinimumFee                 int64
	TransactionToLiveInterval  int
}

func NewParams() Params {
	return Params{
		BlockTime:                  10 * time.Second,
		MaximumTransactionPerBlock: 1000,
		MaximumPower:               21,
		SubsidyReductionInterval:   2100000,
		MaximumMemoLength:          1024,
		FeeFraction:                0.001,
		MinimumFee:                 1000,
		TransactionToLiveInterval:  500,
	}
}
