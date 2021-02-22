package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
)

func (cs *consensus) enterCommit(round int) {
	if cs.isCommitted || round > cs.hrs.Round() {
		cs.logger.Debug("Commit: Committed or invalid round", "round", round)
		return
	}

	precommits := cs.pendingVotes.PrecommitVoteSet(round)
	if !precommits.HasQuorum() {
		cs.logger.Warn("Commit: No quorum for precommit stage")
		return
	}
	cs.updateStep(hrs.StepTypeCommit)

	blockHash := precommits.QuorumBlock()
	if blockHash == nil || blockHash.IsUndef() {
		cs.logger.Error("Commit: Block is invalid")
		return
	}

	// Additional check. blockHash should be same for both prepares and precommits
	prepares := cs.pendingVotes.PrepareVoteSet(round)
	hash := prepares.QuorumBlock()
	if hash == nil || !blockHash.EqualsTo(*hash) {
		cs.logger.Warn("Commit: Commit without prepare quorum")
	}

	// For any reason, we are don't have proposal
	roundProposal := cs.pendingVotes.RoundProposal(round)
	if roundProposal == nil {
		cs.requestForProposal()

		cs.logger.Warn("Commit: No proposal, send proposal request.")
		return
	}

	// Proposal is not for quorum block
	// It is impossible, but good to keep this check
	if !roundProposal.IsForBlock(blockHash) {
		cs.logger.Error("Commit: Proposal is invalid.", "proposal", roundProposal)
		return
	}

	commitBlock := roundProposal.Block()
	commit := precommits.ToCommit()
	height := cs.hrs.Height()
	if commit == nil {
		cs.logger.Error("Commit: Invalid precommits", "precommits", precommits)
		return
	}

	if err := cs.state.CommitBlock(height, commitBlock, *commit); err != nil {
		cs.logger.Warn("Commit: committing block failed", "block", commitBlock, "err", err)
		return
	}

	cs.isCommitted = true
	cs.logger.Info("Commit: Block committed, Schedule new height", "block", blockHash.Fingerprint())

	cs.scheduleNewHeight()

	// Now broadcast the committed block
	cs.broadcastBlock(height, &commitBlock, commit)
}
