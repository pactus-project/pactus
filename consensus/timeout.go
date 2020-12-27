package consensus

import (
	"fmt"
	"time"

	"github.com/zarbchain/zarb-go/consensus/hrs"
)

type timeout struct {
	Duration time.Duration
	Height   int
	Round    int
	Step     hrs.StepType
}

func (ti timeout) String() string {
	return fmt.Sprintf("%v@ %d/%d/%s", ti.Duration, ti.Height, ti.Round, ti.Step)
}
