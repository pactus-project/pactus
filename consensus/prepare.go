package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *consensus) enterPrepare(height int, round int) {
	if cs.invalidHeightRoundStep(height, round, hrs.StepTypePrepareWait) {
		cs.logger.Debug("Prepare: Invalid height/round/step or committed before", "height", height, "round", round, "committed", cs.isCommitted)
		return
	}
	cs.updateRoundStep(round, hrs.StepTypePrepare)

	if cs.votes.lockedProposal != nil {
		cs.logger.Error("Prepare: A block is locked. Unlock it")
		cs.votes.lockedProposal = nil
	}

	roundProposal := cs.votes.RoundProposal(round)
	if roundProposal == nil {
		cs.requestForProposal()

		cs.logger.Warn("Prepare: No proposal")
		cs.signAddVote(vote.VoteTypePrepare, crypto.UndefHash)
		return
	}

	cs.logger.Info("Prepare: Proposal is validated", "proposal", roundProposal)
	cs.signAddVote(vote.VoteTypePrepare, roundProposal.Block().Hash())
}

func (cs *consensus) enterPrepareWait(height int, round int) {
	if cs.invalidHeightRoundStep(height, round, hrs.StepTypePrepareWait) {
		cs.logger.Debug("PrepareWait: Invalid height/round/step or committed before", "height", height, "round", round, "committed", cs.isCommitted)
		return
	}

	cs.updateRoundStep(round, hrs.StepTypePrepareWait)

	cs.logger.Info("PrepareWait: Wait for some more prepares") //then enter precommit
	cs.scheduleTimeout(cs.config.PrepareTimeout(round), height, round, hrs.StepTypePrecommit)
}
