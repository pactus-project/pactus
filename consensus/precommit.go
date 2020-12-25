package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *consensus) enterPrecommit(height int, round int) {
	if cs.invalidHeightRoundStep(height, round, hrs.StepTypePrecommitWait) {
		cs.logger.Debug("Precommit: Invalid height/round/step or committed before", "height", height, "round", round, "committed", cs.isCommitted)
		return
	}

	if cs.votes.lockedProposal != nil {
		cs.logger.Debug("Precommit: we have locked before")
		return
	}

	preVotes := cs.votes.PrepareVoteSet(round)
	if !preVotes.HasQuorum() {
		cs.logger.Debug("Precommit: Entering without prepare quorum")
		return
	}

	// Now, update state and vote!
	cs.updateRoundStep(round, hrs.StepTypePrecommit)

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

	roundProposal := cs.votes.RoundProposal(round)
	if roundProposal == nil {
		cs.requestForProposal()
		cs.logger.Debug("Precommit: No proposal, send proposal request.")
		return
	}

	if !roundProposal.IsForBlock(blockHash) {
		cs.votes.SetRoundProposal(round, nil)
		cs.requestForProposal()
		cs.logger.Warn("Precommit: Invalid proposal, send proposal request.")
		return
	}

	// Everything is good
	cs.votes.lockedProposal = roundProposal
	cs.logger.Info("Precommit: Proposal is locked", "proposal", roundProposal)
	cs.signAddVote(vote.VoteTypePrecommit, *blockHash)
}

func (cs *consensus) enterPrecommitWait(height int, round int) {
	if cs.invalidHeightRoundStep(height, round, hrs.StepTypePrecommitWait) {
		cs.logger.Debug("PrecommitWait: Invalid height/round/step or committed before", "height", height, "round", round, "committed", cs.isCommitted)
		return
	}
	cs.updateRoundStep(round, hrs.StepTypePrecommitWait)

	cs.logger.Info("PrecommitWait: Wait for some more precommits") // ,then enter enterNewRound
	cs.scheduleTimeout(cs.config.PrecommitTimeout(round), height, round, hrs.StepTypeNewRound)
}
