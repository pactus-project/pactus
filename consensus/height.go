package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/util"
)

func (cs *Consensus) scheduleNewHeight() {
	sleep := cs.state.LastBlockTime().Add(cs.state.BlockTime()).Sub(util.Now())
	cs.logger.Debug("NewHeight is scheduled", "seconds", sleep.Seconds())
	cs.scheduleTimeout(sleep, cs.hrs.Height(), cs.hrs.Round(), hrs.StepTypeNewHeight)
}

func (cs *Consensus) enterNewHeight(height int) {
	if cs.hrs.Height() != height-1 || (cs.hrs.Step() != hrs.StepTypeNewHeight && cs.hrs.Step() != hrs.StepTypeCommit) {
		cs.logger.Debug("NewHeight with invalid args", "height", height)
		return
	}
	if cs.state.LastBlockHeight() != cs.hrs.Height() {
		cs.logger.Debug("State is not in same height as consensus", "state", cs.state)
		return
	}

	if height > 1 {
		vs := cs.votes.Precommits(cs.commitRound)
		if vs == nil {
			cs.logger.Warn("Entering new height without having last commit")
		} else {
			// Update last commit here, consensus had enough time to populate votes
			cs.lastCommit = vs.ToCommit()
		}
	}

	cs.commitRound = -1
	cs.votes.Reset(height)
	cs.updateHeight(height)
	cs.updateRoundStep(0, hrs.StepTypeNewHeight)

	cs.enterNewRound(height, 0)
}
