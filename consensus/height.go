package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/util"
)

func (cs *Consensus) ScheduleNewHeight() {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	stateHeight := cs.state.LastBlockHeight()
	if cs.hrs.Height() < stateHeight-1 {
		cs.isCommitted = false
		cs.votes.Reset(stateHeight)
		cs.updateHeight(stateHeight)
		cs.updateRoundStep(0, hrs.StepTypeNewHeight)
	}

	cs.scheduleNewHeight()
}

func (cs *Consensus) scheduleNewHeight() {
	sleep := cs.state.LastBlockTime().Add(cs.state.BlockTime()).Sub(util.Now())
	cs.logger.Debug("NewHeight is scheduled", "seconds", sleep.Seconds())
	cs.scheduleTimeout(sleep, cs.hrs.Height(), 0, hrs.StepTypeNewHeight)
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

	// Apply last committed block
	if cs.votes.lockedProposal != nil {
		vs := cs.votes.Precommits(cs.hrs.Round())
		if vs == nil {
			cs.logger.Warn("Entering new height without having last commit")
		} else {
			// Update last commit here, consensus had enough time to populate more votes
			block := cs.votes.lockedProposal.Block()
			lastCommit := vs.ToCommit()
			if lastCommit != nil {
				cs.state.UpdateLastCommit(block.Hash(), *lastCommit)
			}
		}
	}

	cs.isCommitted = false
	cs.votes.Reset(height)
	cs.updateHeight(height)
	cs.updateRoundStep(0, hrs.StepTypeNewHeight)

	cs.enterNewRound(height, 0)
}
