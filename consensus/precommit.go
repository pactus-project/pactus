package consensus

import (
	"time"

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

	cs.updateStep(hrs.StepTypePrecommit)

	blockHash := prepares.QuorumBlock()
	roundProposal := cs.pendingVotes.RoundProposal(round)
	if roundProposal == nil && blockHash != nil && !blockHash.IsUndef() {
		cs.requestForProposal()
		cs.logger.Debug("Precommit: No proposal, send proposal request.")
		return
	}

	if blockHash == nil && cs.isPrepared {
		// Byzantine node might send different valid proposals to different nodes
		// We wait for a while here to make sure we are not behind the partition
		// At the same time we broadcast our proposal
		// Then decide to vote
		cs.logger.Info("Precommit: It looks some peers don't have proposal yet.")
		cs.broadcastProposal(roundProposal)
		time.Sleep(cs.config.PrecommitTimeout(round))
	}

	// PreCommit vote CAN NOT be updated, unlike prepare
	cs.isPreCommitted = true

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

	// Everything is good
	cs.logger.Info("Precommit: Signed proposal", "proposal", roundProposal)
	cs.signAddVote(vote.VoteTypePrecommit, *blockHash)
}
