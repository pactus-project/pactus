package consensus

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/message"
)

func (cs *consensus) enterCommit(round int) {
	if cs.isCommitted || round != cs.hrs.Round() {
		cs.logger.Debug("Precommit: Precommitted before or invalid round", "round", round)
		return
	}
	cs.updateStep(hrs.StepTypeCommit)

	preVotes := cs.pendingVotes.PrepareVoteSet(round)
	preCommits := cs.pendingVotes.PrecommitVoteSet(round)
	if !preCommits.HasQuorum() {
		cs.logger.Debug("Commit: No quorum for precommit stage")
		return
	}

	blockHash := preCommits.QuorumBlock()
	if blockHash == nil || blockHash.IsUndef() {
		cs.logger.Error("Commit: Block is invalid")
		return
	}

	// Additional check. blockHash should be same for both prepares and precommits
	hash := preVotes.QuorumBlock()
	if hash == nil || !blockHash.EqualsTo(*hash) {
		cs.logger.Debug("Commit: Commit without prepare quorum")
	}

	// For any reason, we are not locked, try to found the locked proposal
	roundProposal := cs.pendingVotes.RoundProposal(round)
	if roundProposal == nil {
		cs.requestForProposal()

		cs.logger.Debug("Commit: No proposal, send proposal request.")
		return
	}

	// Locked proposal is not for quorum block
	// It is impossible, but good to keep this check
	if !roundProposal.IsForBlock(blockHash) {
		cs.logger.Error("Commit: Proposal is invalid.", "proposal", roundProposal)
		return
	}

	commitBlock := roundProposal.Block()
	commit := preCommits.ToCommit()
	height := cs.hrs.Height()
	if commit == nil {
		cs.logger.Error("Commit: Invalid precommits", "preCommits", preCommits)
		return
	}

	if err := cs.state.ApplyBlock(height, commitBlock, *commit); err != nil {
		cs.logger.Error("Commit: Applying block failed", "block", commitBlock, "err", err)
		return
	}

	cs.isCommitted = true
	cs.logger.Info("Commit: Block committed, Schedule new height", "block", blockHash.Fingerprint())

	cs.scheduleNewHeight()

	// Now broadcast the committed block
	msg := message.NewLatestBlocksMessage(height, []*block.Block{&commitBlock}, nil, commit)
	cs.broadcastCh <- msg
}
