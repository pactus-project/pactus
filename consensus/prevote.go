package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *consensus) enterPrevote(height int, round int) {
	if cs.invalidHeightRoundStep(height, round, hrs.StepTypePrevoteWait) {
		cs.logger.Debug("Prevote: Invalid height/round/step or committed before", "height", height, "round", round, "committed", cs.isCommitted)
		return
	}
	cs.updateRoundStep(round, hrs.StepTypePrevote)

	if cs.votes.lockedProposal != nil {
		cs.logger.Error("Prevote: A block is locked. Unlock it")
		cs.votes.lockedProposal = nil
	}

	roundProposal := cs.votes.RoundProposal(round)
	if roundProposal == nil {
		cs.requestForProposal()

		cs.logger.Warn("Prevote: No proposal")
		cs.signAddVote(vote.VoteTypePrevote, crypto.UndefHash)
		return
	}

	cs.logger.Info("Prevote: Proposal is validated", "proposal", roundProposal)
	cs.signAddVote(vote.VoteTypePrevote, roundProposal.Block().Hash())
}

func (cs *consensus) enterPrevoteWait(height int, round int) {
	if cs.invalidHeightRoundStep(height, round, hrs.StepTypePrevoteWait) {
		cs.logger.Debug("PrevoteWait: Invalid height/round/step or committed before", "height", height, "round", round, "committed", cs.isCommitted)
		return
	}

	cs.updateRoundStep(round, hrs.StepTypePrevoteWait)

	cs.logger.Info("PrevoteWait: Wait for some more prevotes") //then enter precommit
	cs.scheduleTimeout(cs.config.Prevote(round), height, round, hrs.StepTypePrecommit)
}
