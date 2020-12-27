package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *consensus) enterPrecommit(round int) {
	if cs.isPreCommitted || round != cs.hrs.Round() {
		cs.logger.Debug("Precommit: Precommitted before or invalid round", "round", round)
		return
	}

	preVotes := cs.pendingVotes.PrepareVoteSet(round)
	if !preVotes.HasQuorum() {
		cs.logger.Debug("Precommit: Entering without prepare quorum")
		return
	}

	// Make sure we have passed prepared stage before entring precommit stage
	cs.updateStep(hrs.StepTypePrecommit)

	blockHash := preVotes.QuorumBlock()
	if blockHash == nil {
		cs.logger.Info("Precommit: No quorum for prepare")
		cs.signAddVote(vote.VoteTypePrecommit, crypto.UndefHash)
		return
	}

	if blockHash.IsUndef() {
		cs.logger.Info("Precommit: Undef quorum for prepare")
		cs.signAddVote(vote.VoteTypePrecommit, crypto.UndefHash)
		return
	}

	roundProposal := cs.pendingVotes.RoundProposal(round)
	if roundProposal == nil {
		cs.requestForProposal()
		cs.logger.Debug("Precommit: No proposal, send proposal request.")
		return
	}

	if !roundProposal.IsForBlock(blockHash) {
		cs.pendingVotes.SetRoundProposal(round, nil)
		cs.requestForProposal()
		cs.logger.Warn("Precommit: Invalid proposal, send proposal request.")
		return
	}

	// Everything is good
	cs.isPreCommitted = true
	cs.logger.Info("Precommit: Proposal is locked", "proposal", roundProposal)
	cs.signAddVote(vote.VoteTypePrecommit, *blockHash)
}
