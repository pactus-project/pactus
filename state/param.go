package state

import "time"

type Params struct {
	BlockTime    time.Duration
	MaximumPower int
}

func NewParams() *Params {
	return &Params{
		BlockTime:    10 * time.Second,
		MaximumPower: 5,
	}
}
