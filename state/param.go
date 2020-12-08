package state

import "time"

type Params struct {
	BlockTime                 time.Duration
	MaximumPower              int
	SubsidyReductionInterval  int
	MaximumMemoLength         int
	FeeFraction               float64
	MinimumFee                int64
	TransactionToLiveInterval int
}

func MainnetParams() Params {
	return Params{
		BlockTime:                 10 * time.Second,
		MaximumPower:              21,
		SubsidyReductionInterval:  2100000,
		MaximumMemoLength:         1024,
		FeeFraction:               0.001,
		MinimumFee:                1000,
		TransactionToLiveInterval: 500,
	}
}

func TestnetParams() Params {
	return Params{
		BlockTime:                 5 * time.Second,
		MaximumPower:              5,
		SubsidyReductionInterval:  210000,
		MaximumMemoLength:         1024,
		FeeFraction:               0.001,
		MinimumFee:                1000,
		TransactionToLiveInterval: 500,
	}
}

func TestParams() Params {
	return Params{
		BlockTime:                 1 * time.Second,
		MaximumPower:              5,
		SubsidyReductionInterval:  210000,
		MaximumMemoLength:         1024,
		FeeFraction:               0.001,
		MinimumFee:                1000,
		TransactionToLiveInterval: 4,
	}
}
