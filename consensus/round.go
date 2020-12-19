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

	// make sure we have quorom votes for previous round
	if round > 0 {
		if !cs.votes.Prevotes(round - 1).HasQuorum() {
			cs.logger.Debug("NewRound: No prevote quorom for previous round")
		}
		if !cs.votes.Precommits(round - 1).HasQuorum() {
			cs.logger.Error("NewRound: No precommit quorom for previous round")
			return
		}
		// Normally when there is no proposal for this round, every one should vote for UndefHash
		prevoteBlockHash := cs.votes.Prevotes(round - 1).QuorumBlock()
		precommitBlockHash := cs.votes.Precommits(round - 1).QuorumBlock()
		if prevoteBlockHash == nil || !prevoteBlockHash.IsUndef() {
			cs.logger.Warn("NewRound: Suspicious prevotes", "blockHash", prevoteBlockHash)
		}
		if precommitBlockHash == nil || !precommitBlockHash.IsUndef() {
			cs.logger.Warn("NewRound: Suspicious precommits", "blockHash", precommitBlockHash)
		}
	}

	cs.votes.lockedProposal = nil
	cs.updateRoundStep(round, hrs.StepTypeNewRound)
	cs.logger.Info("NewRound: Entring new round", "round", round)

	cs.enterPropose(height, round)
}
