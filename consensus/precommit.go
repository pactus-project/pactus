package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *consensus) enterPrecommit(round int) {
	if cs.isPreCommitted || round != cs.hrs.Round() {
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

	if blockHash == nil && cs.isPrepared {
		// We have a valid proposal, but there is no consensus about it
		//
		// If we are behind the partition, it might be easy to find it here
		// There should be some null-votes here
		// If number of null-votes are greather tha `1f` (`f` stands for faulty)
		// Then we broadcast our proposal and return here
		//
		// Note: Byzantine node might send different valid proposals to different nodes
		//
		cs.logger.Info("Precommit: Some peers don't have proposal yet.")

		votes := prepares.AllVotes()
		count := 0
		for _, v := range votes {
			if v.BlockHash().IsUndef() {
				count++
			}
		}

		if count > len(votes)/3 {
			cs.logger.Debug("Precommit: Broadcst proposal.")
			cs.broadcastProposal(roundProposal)
			return
		}
	}

	cs.updateStep(hrs.StepTypePrecommit)

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

	if !roundProposal.IsForBlock(blockHash) {
		cs.pendingVotes.SetRoundProposal(round, nil)
		cs.logger.Warn("Precommit: Invalid proposal.")
		cs.signAddVote(vote.VoteTypePrecommit, crypto.UndefHash)
		return
	}

	cs.isPreCommitted = true

	// Everything is good
	cs.logger.Info("Precommit: Proposal signed", "proposal", roundProposal)
	cs.signAddVote(vote.VoteTypePrecommit, *blockHash)
}
