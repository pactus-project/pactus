package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *Consensus) enterPrecommit(height int, round int) {
	if cs.invalidHeightRoundStep(height, round, hrs.StepTypePrecommit) {
		cs.logger.Debug("Precommit: Invalid height/round/step or committed before", "height", height, "round", round, "committed", cs.isCommitted)
		return
	}

	if cs.votes.lockedProposal != nil {
		cs.logger.Debug("Precommit: we have locked before")
		return
	}

	preVotes := cs.votes.Prevotes(round)
	if !preVotes.HasQuorum() {
		cs.logger.Error("Precommit: Entering precommit witout having quorom for prevote stage")
		return
	}

	// Now, update state and vote!
	cs.updateRoundStep(round, hrs.StepTypePrecommit)

	blockHash := preVotes.QuorumBlock()
	if blockHash == nil || blockHash.IsUndef() {
		cs.logger.Info("Precommit: No quorum for prevote")
		cs.signAddVote(vote.VoteTypePrecommit, crypto.UndefHash)
		return
	}

	roundProposal := cs.votes.RoundProposal(round)
	if roundProposal == nil {
		cs.requestForProposal()

		cs.logger.Error("Precommit: No proposal")
		cs.signAddVote(vote.VoteTypePrevote, crypto.UndefHash)
		return
	}

	if !roundProposal.IsForBlock(blockHash) {
		cs.logger.Error("Precommit: Unknown proposal")
		cs.signAddVote(vote.VoteTypePrevote, crypto.UndefHash)
		return
	}

	if err := cs.state.ValidateBlock(roundProposal.Block()); err != nil {
		cs.logger.Warn("Precommit: Invalid block", "proposal", roundProposal, "err", err)
		cs.signAddVote(vote.VoteTypePrevote, crypto.UndefHash)
		return

	}

	// Everything is good
	cs.votes.lockedProposal = roundProposal
	cs.logger.Info("Precommit: Proposal is locked", "proposal", roundProposal)
	cs.signAddVote(vote.VoteTypePrecommit, *blockHash)
}

func (cs *Consensus) enterPrecommitWait(height int, round int) {
	if cs.invalidHeightRoundStep(height, round, hrs.StepTypePrecommitWait) {
		cs.logger.Debug("PrecommitWait: Invalid height/round/step or committed before", "height", height, "round", round, "committed", cs.isCommitted)
		return
	}
	cs.updateRoundStep(round, hrs.StepTypePrecommitWait)

	if !cs.votes.Precommits(round).HasQuorum() {
		cs.logger.Error("PrecommitWait: Precommits does not have any +2/3 votes")
	}

	cs.logger.Info("PrecommitWait: Wait for some more precommits") // ,then enter enterNewRound
	cs.scheduleTimeout(cs.config.Precommit(round), height, round, hrs.StepTypeNewRound)
}
