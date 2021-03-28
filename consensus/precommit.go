package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *consensus) enterPrecommit(round int) {
	if cs.status.IsPreCommitted() || round > cs.hrs.Round() {
		cs.logger.Debug("Precommit: Precommitted or invalid round/step", "round", round)
		return
	}

	prepares := cs.pendingVotes.PrepareVoteSet(round)
	if !prepares.HasQuorum() {
		cs.logger.Debug("Precommit: Entering without prepare quorum")
		return
	}

	blockHash := prepares.QuorumBlock()
	roundProposal := cs.pendingVotes.RoundProposal(round)
	if roundProposal == nil && blockHash != nil && !blockHash.IsUndef() {
		// There is a consensus about a proposal which we don't have it yet.
		// Ask peers for this proposal
		cs.requestForProposal()
		cs.logger.Debug("Precommit: No proposal, send proposal request.")
		return
	}

	if roundProposal != nil && blockHash == nil {
		// We have a valid proposal, but there is no consensus about it
		//
		// If we are behind the partition, it might be easy to find it here
		// There should be some null-votes here.
		// If weight of null-votes are greather than `1f` (`f` stands for faulty)
		// Then we broadcast our proposal and return here
		//
		// Note: Byzantine node might send different valid proposals to different nodes
		//
		cs.logger.Info("Precommit: Some peers don't have proposal yet.")

		if prepares.HasOneThirdOfTotalPower(crypto.UndefHash) {
			cs.logger.Debug("Precommit: Broadcst proposal.", "proposal", roundProposal)
			cs.broadcastProposal(roundProposal)
			return
		}
	}

	cs.hrs.UpdateStep(hrs.StepTypePrecommit)

	if blockHash == nil {
		cs.logger.Info("Precommit: No quorum for prepare")
		cs.signAddVote(vote.VoteTypePrecommit, round, crypto.UndefHash)
		return
	}

	if blockHash.IsUndef() {
		cs.logger.Info("Precommit: Undef quorum for prepare")
		cs.signAddVote(vote.VoteTypePrecommit, round, crypto.UndefHash)
		return
	}

	if !roundProposal.IsForBlock(blockHash) {
		cs.pendingVotes.SetRoundProposal(round, nil)
		cs.logger.Warn("Precommit: Invalid proposal.")
		cs.signAddVote(vote.VoteTypePrecommit, round, crypto.UndefHash)
		return
	}

	// Everything is good
	cs.status.SetPreCommitted(true)
	cs.logger.Info("Precommit: Proposal approved", "proposal", roundProposal)
	cs.signAddVote(vote.VoteTypePrecommit, round, *blockHash)
}
