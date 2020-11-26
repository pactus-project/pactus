package state

import "time"

type Params struct {
	BlockTime                time.Duration
	MaximumPower             int
	SubsidyReductionInterval int
	MaximumMemoLength        int
	FeeFraction              float64
	MinimumFee               int64
}

func NewParams() *Params {
	return &Params{
		BlockTime:                10 * time.Second,
		MaximumPower:             5,
		SubsidyReductionInterval: 210000,
		MaximumMemoLength:        1024,
		FeeFraction:              0.001,
		MinimumFee:               1000,
	}
}
