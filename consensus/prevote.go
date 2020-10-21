package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *Consensus) enterPrevote(height int, round int) {
	if cs.hrs.InvalidHeightRoundStep(height, round, hrs.StepTypePrevote) {
		cs.logger.Debug("Prevote with invalid args", "height", height, "round", round)
		return
	}
	cs.updateRoundStep(round, hrs.StepTypePrevote)

	if cs.votes.lockedProposal != nil {
		cs.logger.Error("Prevote: A block is locked. Unlock it")
		cs.votes.lockedProposal = nil
	}

	roundProposal := cs.votes.RoundProposal(round)
	if roundProposal == nil {
		cs.logger.Warn("Prevote: Voted for nil, no proposal.")
		cs.signAddVote(vote.VoteTypePrevote, crypto.UndefHash)
		return
	}

	if err := cs.state.ValidateBlock(roundProposal.Block()); err != nil {
		cs.logger.Warn("Prevote: Voted for nil, invalid block", "Proposal", roundProposal, "err", err)
		cs.signAddVote(vote.VoteTypePrevote, crypto.UndefHash)
		return
	}

	cs.logger.Info("Prevote: Proposal is validated", "proposal", roundProposal)
	cs.signAddVote(vote.VoteTypePrevote, roundProposal.Block().Hash())
}

func (cs *Consensus) enterPrevoteWait(height int, round int) {
	if cs.hrs.InvalidHeightRoundStep(height, round, hrs.StepTypePrevoteWait) {
		cs.logger.Debug("PrevoteWait with invalid args", "height", height, "round", round)
		return
	}

	cs.updateRoundStep(round, hrs.StepTypePrevoteWait)

	if !cs.votes.Prevotes(round).HasQuorum() {
		cs.logger.Error("PrevoteWait: Prevotes does not have any +2/3 votes")
	}

	cs.scheduleTimeout(cs.config.Prevote(round), height, round, hrs.StepTypePrecommit)
	cs.logger.Info("Wait for some more prevotes") //then enter precommit
}
