package state

import "time"

type Params struct {
	BlockTime time.Duration
}

func NewParams() *Params {
	return &Params{
		BlockTime: 10 * time.Second,
	}
}
