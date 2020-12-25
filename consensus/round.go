package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
)

func (cs *consensus) enterNewRound(height int, round int) {
	if cs.invalidHeight(height) {
		cs.logger.Debug("NewRound: Invalid height or committed before", "height", height, "round", round, "committed", cs.isCommitted)
		return
	}

	if round < cs.hrs.Round() {
		cs.logger.Debug("NewRound: Try to enter prior round", "height", height, "round", round)
		return
	}

	// make sure we have quorum votes for previous round
	if round > 0 {
		prepares := cs.votes.PrepareVoteSet(round - 1)
		precommits := cs.votes.PrecommitVoteSet(round - 1)
		if !prepares.HasQuorum() {
			cs.logger.Debug("NewRound: No prepare quorum for previous round")
		}
		if !precommits.HasQuorum() {
			cs.logger.Error("NewRound: No precommit quorum for previous round")
			return
		}
		// Normally when there is no proposal for this round, every one should vote for nil
		prepareBlockHash := prepares.QuorumBlock()
		precommitBlockHash := precommits.QuorumBlock()
		if prepareBlockHash == nil || !prepareBlockHash.IsUndef() {
			cs.logger.Warn("NewRound: Suspicious prepares", "blockHash", prepareBlockHash)
		}
		if precommitBlockHash == nil || !precommitBlockHash.IsUndef() {
			cs.logger.Warn("NewRound: Suspicious precommits", "blockHash", precommitBlockHash)
		}
	}

	cs.votes.lockedProposal = nil
	cs.updateRoundStep(round, hrs.StepTypeNewRound)
	cs.logger.Info("NewRound: Entering new round", "round", round)

	cs.enterPropose(height, round)
}
