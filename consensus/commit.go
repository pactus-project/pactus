package consensus

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/message"
)

func (cs *consensus) enterCommit(height int, round int) {
	if cs.invalidHeight(height) {
		cs.logger.Debug("Commit: Invalid height or committed before", "height", height, "committed", cs.isCommitted)
		return
	}

	preVotes := cs.votes.Prevotes(round)
	preCommits := cs.votes.Precommits(round)

	if !preCommits.HasQuorum() {
		cs.logger.Debug("Commit: No quorum for precommit stage")
		return
	}

	blockHash := preCommits.QuorumBlock()
	if blockHash == nil || blockHash.IsUndef() {
		cs.logger.Error("Commit: Block is invalid")
		return
	}

	// Additional check. blockHash should be same for both prevotes and precommits
	prevoteBlockHash := preVotes.QuorumBlock()
	if prevoteBlockHash == nil || !blockHash.EqualsTo(*prevoteBlockHash) {
		cs.logger.Debug("Commit: Commit without prevote quorum")
	}

	if cs.votes.lockedProposal == nil {
		// For any reason, we are not locked, try to found the locked proposal
		roundProposal := cs.votes.RoundProposal(round)
		if roundProposal == nil {
			cs.requestForProposal()

			cs.logger.Debug("Commit: No proposal, send proposal request.")
			return
		} else if roundProposal.IsForBlock(blockHash) {
			cs.votes.lockedProposal = roundProposal
		} else {
			cs.logger.Error("Commit: Invalid proposal.", "proposal", roundProposal)
			return
		}
	}

	// Locked proposal is not for quorum block
	// It is impossible, but good to keep this check
	if !cs.votes.lockedProposal.IsForBlock(blockHash) {
		cs.votes.lockedProposal = nil
		cs.logger.Error("Commit: Proposal is invalid.", "proposal", cs.votes.lockedProposal)
		return
	}

	// Locked proposal should be same as round proposal
	// It is impossible, but good to keep this check
	roundProposal := cs.votes.RoundProposal(round)
	if roundProposal != nil && !roundProposal.IsForBlock(blockHash) {
		cs.votes.lockedProposal = nil
		cs.logger.Error("Commit: Proposal is not for this round.", "proposal", cs.votes.lockedProposal, "round-proposal", roundProposal)
		return
	}

	commitBlock := cs.votes.lockedProposal.Block()
	commit := preCommits.ToCommit()
	if commit == nil {
		cs.logger.Error("Commit: Invalid precommits", "preCommits", preCommits)
		return
	}

	if err := cs.state.ApplyBlock(height, commitBlock, *commit); err != nil {
		cs.logger.Error("Commit: Applying block failed", "block", commitBlock, "err", err)
		return
	}

	cs.isCommitted = true
	cs.updateRoundStep(round, hrs.StepTypeCommit)
	cs.logger.Info("Commit: Block committed, Schedule new height", "block", blockHash.Fingerprint())

	cs.scheduleNewHeight()

	// Now broadcast the committed block
	msg := message.NewBlocksMessage(height, []*block.Block{&commitBlock}, commit)
	cs.broadcastCh <- msg
}
