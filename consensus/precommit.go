package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *Consensus) enterPrecommit(height int, round int) {
	if cs.hrs.InvalidHeightRoundStep(height, round, hrs.StepTypePrecommit) {
		cs.logger.Debug("Precommit with invalid args", "height", height, "round", round)
		return
	}

	if !cs.votes.Prevotes(round).HasQuorum() {
		cs.logger.Error("Precommit: Entering precommit witout having quorom for prevote stage")
		return
	}

	// Now, update state and vote!
	cs.updateRoundStep(round, hrs.StepTypePrecommit)

	blockHash := cs.votes.Prevotes(round).QuorumBlock()
	if blockHash == nil || blockHash.IsUndef() {
		cs.logger.Info("Precommit: Voted for nil")
		cs.signAddVote(vote.VoteTypePrecommit, crypto.UndefHash)
		return
	}

	roundProposal := cs.votes.RoundProposal(round)
	if roundProposal == nil {
		cs.logger.Error("Precommit: We don't have proposal yet")
		cs.signAddVote(vote.VoteTypePrevote, crypto.UndefHash)
	}

	if !roundProposal.IsForBlock(blockHash) {
		cs.logger.Error("Precommit: We have unknown proposal")
		cs.signAddVote(vote.VoteTypePrevote, crypto.UndefHash)
		return
	}

	if err := cs.state.ValidateBlock(roundProposal.Block()); err != nil {
		cs.logger.Warn("Precommit: Voted for nil, invalid block", "Proposal", roundProposal, "err", err)
		cs.signAddVote(vote.VoteTypePrevote, crypto.UndefHash)
		return

	}

	// Everything is good
	cs.votes.lockedProposal = roundProposal
	cs.logger.Info("Precommit: Proposal is locked", "proposal", roundProposal)
	cs.signAddVote(vote.VoteTypePrecommit, *blockHash)
}

func (cs *Consensus) enterPrecommitWait(height int, round int) {
	if cs.hrs.InvalidHeightRoundStep(height, round, hrs.StepTypePrecommitWait) {
		cs.logger.Debug("PrecommitWait with invalid args", "height", height, "round", round)
		return
	}
	cs.updateRoundStep(round, hrs.StepTypePrecommitWait)

	if !cs.votes.Precommits(round).HasQuorum() {
		cs.logger.Error("PrecommitWait, but Precommits does not have any +2/3 votes")
	}

	cs.scheduleTimeout(cs.config.Consensus.Precommit(round), height, round, hrs.StepTypeNewRound)
	cs.logger.Info("Wait for some more precommits") // ,then enter enterNewRound
}
