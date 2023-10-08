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
	tickerTargetQueryVotes     = tickerTarget(4)
)

func (rs tickerTarget) String() string {
	switch rs {
	case tickerTargetNewHeight:
		return "new-height"
	case tickerTargetChangeProposer:
		return "change-proposer"
	case tickerTargetQueryProposal:
		return "query-proposal"
	case tickerTargetQueryVotes:
		return "query-votes"
	default:
		return "Unknown"
	}
}

type ticker struct {
	Duration time.Duration
	Height   uint32
	Round    int16
	Target   tickerTarget
}

func (ti ticker) String() string {
	return fmt.Sprintf("%v@ %d/%d/%s", ti.Duration, ti.Height, ti.Round, ti.Target)
}
