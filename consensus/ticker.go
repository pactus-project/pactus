package consensus

import (
	"fmt"
	"time"
)

type tickerTarget int

const (
	tickerTargetNewHeight      = tickerTarget(1)
	tickerTargetChangeProposer = tickerTarget(2)
	tickerTargetQueryProposal  = tickerTarget(3)
)

func (rs tickerTarget) String() string {
	switch rs {
	case tickerTargetNewHeight:
		return "new-height"
	case tickerTargetChangeProposer:
		return "change-proposer"
	case tickerTargetQueryProposal:
		return "query-proposal"
	default:
		return "Unknown"
	}
}

type ticker struct {
	Duration time.Duration
	Height   int32
	Round    int16
	Target   tickerTarget
}

func (ti ticker) String() string {
	return fmt.Sprintf("%v@ %d/%d/%s", ti.Duration, ti.Height, ti.Round, ti.Target)
}
