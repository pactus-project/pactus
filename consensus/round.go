package consensus

import (
	"gitlab.com/zarb-chain/zarb-go/consensus/hrs"
)

func (cs *Consensus) enterNewRound(height int, round int) {
	if cs.hrs.InvalidHeight(height) {
		cs.logger.Debug("NewRound with invalid args", "height", height, "round", round)
		return
	}
	switch cs.hrs.Step() {
	case hrs.StepTypeCommit:
	case hrs.StepTypeNewHeight:
		{
			if round != 0 || cs.hrs.Round() != 0 {
				cs.logger.Debug("NewRound with invalid args", "height", height, "round", round)
				return
			}
		}
	case hrs.StepTypePrecommit:
	case hrs.StepTypePrecommitWait:
		{
			if round == 0 || cs.hrs.Round() != round-1 {
				cs.logger.Debug("NewRound with invalid args", "height", height, "round", round)
				return
			}
		}
	default:
		cs.logger.Debug("NewRound with invalid args", "height", height, "round", round)
		return
	}

	// make sure we have quorom nil votes for previous round
	if round > 0 {
		if !cs.votes.Prevotes(round - 1).HasQuorum() {
			cs.logger.Error("NewRound: No prevote quorom for previous round")
			return
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

	cs.updateRoundStep(round, hrs.StepTypeNewRound)

	cs.logger.Debug("Resetting Proposal info")
	cs.votes.lockedProposal = nil

	cs.enterPropose(height, round)
}
