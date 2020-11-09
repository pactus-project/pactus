package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/util"
)

func (cs *Consensus) ScheduleNewHeight() {
	cs.lk.RLock()
	defer cs.lk.RUnlock()

	stateHeight := cs.state.LastBlockHeight()
	consHeight := cs.hrs.Height()

	if consHeight > stateHeight+1 {
		cs.logger.Panic("Consensus can't be further than state")
		return
	}

	if stateHeight == 0 {
		cs.updateRoundStep(0, hrs.StepTypeNewHeight)
	} else {
		cs.updateRoundStep(0, hrs.StepTypeCommit)
	}

	cs.isCommitted = true
	cs.updateHeight(stateHeight)

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
		cs.logger.Debug("NewHeight: State is not in same height as consensus", "state", cs.state)
		return
	}

	// Apply last committed block
	if cs.votes.lockedProposal != nil {
		vs := cs.votes.Precommits(cs.hrs.Round())
		if vs == nil {
			cs.logger.Warn("NewHeight: Entering new height without having last commit")
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
