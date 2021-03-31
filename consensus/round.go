package consensus

import "github.com/zarbchain/zarb-go/consensus/hrs"

func (cs *consensus) enterNewRound(round int) {
	if round > 0 && round <= cs.hrs.Round() {
		cs.logger.Trace("NewRound: Try to enter prior round", "round", round)
		return
	}

	// make sure we have quorum votes for previous round
	if round > 0 {
		prepares := cs.pendingVotes.PrepareVoteSet(round - 1)
		precommits := cs.pendingVotes.PrecommitVoteSet(round - 1)
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

	cs.status.SetProposed(false)
	cs.status.SetPrepared(false)
	cs.hrs.UpdateRound(round)
	cs.hrs.UpdateStep(hrs.StepTypeNewRound)
	cs.logger.Info("NewRound: Entering new round", "round", round)

	cs.enterPropose(round)
}
