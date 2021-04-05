package consensus

import (
	"fmt"
	"time"
)

type tickerTarget int

const (
	tickerTargetNewHeight      = tickerTarget(1)
	tickerTargetChangeProposer = tickerTarget(2)
)

func (rs tickerTarget) String() string {
	switch rs {
	case tickerTargetNewHeight:
		return "new-height"
	case tickerTargetChangeProposer:
		return "change-proposer"
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
