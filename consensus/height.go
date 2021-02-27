package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/util"
)

func (cs *consensus) MoveToNewHeight() {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	cs.scheduleNewHeight()
}

func (cs *consensus) scheduleNewHeight() {
	sleep := cs.state.LastBlockTime().Add(cs.state.BlockTime()).Sub(util.Now())
	cs.logger.Debug("NewHeight is scheduled", "seconds", sleep.Seconds())
	cs.scheduleTimeout(sleep, cs.hrs.Height(), 0, hrs.StepTypeNewHeight)
}

func (cs *consensus) enterNewHeight() {
	sateHeight := cs.state.LastBlockHeight()
	if cs.hrs.Height() == sateHeight+1 {
		cs.logger.Debug("NewHeight: Duplicated entry")
		return
	}

	// Apply last committed block, We might have more votes now
	if cs.hrs.Height() == sateHeight && cs.hrs.Round() >= 0 {
		vs := cs.pendingVotes.PrecommitVoteSet(cs.hrs.Round())
		if vs == nil {
			cs.logger.Warn("NewHeight: Entering new height without last commit")
		} else {
			// Update last commit here, consensus had enough time to populate more votes
			lastCommit := vs.ToCommit()
			if lastCommit != nil {
				if err := cs.state.UpdateLastCommit(lastCommit); err != nil {
					cs.logger.Warn("NewHeight: Updating last commit failed", "err", err)
				}
			}
		}
	}

	vals := cs.state.ValidatorSet().Validators()
	cs.pendingVotes.MoveToNewHeight(sateHeight+1, vals)

	cs.updateHeight(sateHeight + 1)
	cs.updateRound(-1)
	cs.updateStep(hrs.StepTypeNewHeight)
	cs.logger.Info("NewHeight: Entering new height", "height", sateHeight+1)

	cs.enterNewRound(0)
}
