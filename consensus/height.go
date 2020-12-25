package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/util"
)

func (cs *consensus) MoveToNewHeight() {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	stateHeight := cs.state.LastBlockHeight()
	consHeight := cs.hrs.Height()

	if consHeight > stateHeight+1 {
		cs.logger.Panic("Consensus can't be further than state by more than one block")
		return
	}

	cs.isCommitted = true
	if stateHeight == 0 {
		cs.updateRoundStep(0, hrs.StepTypeNewHeight)
	} else {
		cs.updateRoundStep(0, hrs.StepTypeCommit)
	}
	cs.updateHeight(stateHeight)
	cs.scheduleNewHeight()
}

func (cs *consensus) scheduleNewHeight() {
	sleep := cs.state.LastBlockTime().Add(cs.state.BlockTime()).Sub(util.Now())
	cs.logger.Debug("NewHeight is scheduled", "seconds", sleep.Seconds())
	cs.scheduleTimeout(sleep, cs.hrs.Height(), 0, hrs.StepTypeNewHeight)
}

func (cs *consensus) enterNewHeight(height int) {
	if cs.hrs.Height() != height-1 || (cs.hrs.Step() != hrs.StepTypeNewHeight && cs.hrs.Step() != hrs.StepTypeCommit) {
		cs.logger.Debug("NewHeight: Invalid args", "height", height)
		return
	}
	if cs.state.LastBlockHeight() != cs.hrs.Height() {
		cs.logger.Debug("NewHeight: State is not in same height as consensus", "state", cs.state)
		return
	}

	// Apply last committed block
	if cs.votes.lockedProposal != nil {
		vs := cs.votes.PrecommitVoteSet(cs.hrs.Round())
		if vs == nil {
			cs.logger.Warn("NewHeight: Entering new height without last commit")
		} else {
			// TODO: add test for me
			// Update last commit here, consensus had enough time to populate more votes
			lastCommit := vs.ToCommit()
			if lastCommit != nil {
				if err := cs.state.UpdateLastCommit(lastCommit); err != nil {
					cs.logger.Warn("NewHeight: Updating last commit failed", "err", err)
				}
			}
		}
	}

	cs.isCommitted = false
	cs.updateHeight(height)
	cs.updateRoundStep(0, hrs.StepTypeNewHeight)
	cs.logger.Info("NewHeight: Entering new height", "height", height)

	cs.enterNewRound(height, 0)
}
