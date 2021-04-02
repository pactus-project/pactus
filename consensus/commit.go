package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
)

func (cs *consensus) enterCommit(round int) {
	if cs.isCommitted || round > cs.hrs.Round() {
		cs.logger.Trace("Commit: Committed or invalid round", "round", round)
		return
	}

	precommits := cs.pendingVotes.PrecommitVoteSet(round)
	if !precommits.HasQuorum() {
		cs.logger.Warn("Commit: No quorum for precommit stage")
		return
	}

	blockHash := precommits.QuorumBlock()
	if blockHash == nil || blockHash.IsUndef() {
		cs.logger.Error("Commit: quorum block  hash is invalid", "hash", blockHash)
		return
	}
	cs.hrs.UpdateStep(hrs.StepTypeCommit)

	// Additional check. blockHash should be same for both prepares and precommits
	prepares := cs.pendingVotes.PrepareVoteSet(round)
	hash := prepares.QuorumBlock()
	if hash == nil || !blockHash.EqualsTo(*hash) {
		cs.logger.Warn("Commit: Commit without prepare quorum", "hash", hash)
	}

	// For any reason, we don't have proposal
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

	certBlock := roundProposal.Block()
	cert := precommits.ToCertificate()
	height := cs.hrs.Height()
	if cert == nil {
		cs.logger.Error("Commit: Invalid precommits", "precommits", precommits)
		return
	}

	if err := cs.state.CommitBlock(height, certBlock, *cert); err != nil {
		cs.logger.Warn("Commit: committing block failed", "block", certBlock, "err", err)
		return
	}

	cs.isCommitted = true
	cs.logger.Info("Commit: Block committed, Schedule new height", "block", blockHash.Fingerprint())

	cs.scheduleNewHeight()

	// Now broadcast the committed block
	cs.broadcastBlock(height, &certBlock, cert)
}
