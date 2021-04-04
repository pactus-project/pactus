package consensus

import (
	"fmt"
	"time"
)

type tickerTarget int

const (
	tickerTargetNewHeight = tickerTarget(1)
	tickerTargetPropose   = tickerTarget(2)
	tickerTargetPrepare   = tickerTarget(3)
	tickerTargetPrecommit = tickerTarget(4)
)

func (rs tickerTarget) String() string {
	switch rs {
	case tickerTargetNewHeight:
		return newHeightName
	case tickerTargetPropose:
		return proposeName
	case tickerTargetPrepare:
		return prepareName
	case tickerTargetPrecommit:
		return precommitName
	default:
		return "Unknown"
	}
}

type ticker struct {
	Duration time.Duration
	Height   int
	Round    int
	Target   tickerTarget
}

func (ti ticker) String() string {
	return fmt.Sprintf("%v@ %d/%d/%s", ti.Duration, ti.Height, ti.Round, ti.Target)
}
