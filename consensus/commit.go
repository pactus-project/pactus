package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/message"
)

func (cs *Consensus) enterCommit(height int, round int) {
	if cs.hrs.InvalidHeight(height) || cs.isCommitted {
		cs.logger.Debug("Commit with invalid args or committed before", "height", height, "committed", cs.isCommitted)
		return
	}

	preVotes := cs.votes.Prevotes(round)
	preCommits := cs.votes.Precommits(round)

	if !preCommits.HasQuorum() {
		cs.logger.Error("Commit witout quorom for precommit stage")
		return
	}

	blockHash := preCommits.QuorumBlock()
	if blockHash == nil || blockHash.IsUndef() {
		cs.logger.Error("Commit is for invalid block")
		return
	}

	// Additional check. blockHash should be same for both prevotes and precommits
	prevoteBlockHash := preVotes.QuorumBlock()
	if prevoteBlockHash == nil || !blockHash.EqualsTo(*prevoteBlockHash) {
		cs.logger.Error("Commit witout quorom for prevote stage")
		return
	}

	if cs.votes.lockedProposal == nil {
		// For any reason, we are not locked, try to found the locked proposal
		roundProposal := cs.votes.RoundProposal(round)
		if roundProposal != nil && roundProposal.IsForBlock(blockHash) {
			cs.votes.lockedProposal = roundProposal
		} else {
			cs.logger.Error("We don't have commit proposal.")
			return
		}
	}

	// Locked proposal is not for quorom block
	// It is impossible, but good to keep this check
	if !cs.votes.lockedProposal.IsForBlock(blockHash) {
		cs.votes.lockedProposal = nil
		cs.logger.Error("Commit proposal is invalid.", "proposal", cs.votes.lockedProposal)
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

	// Block is invalid
	// It is impossible, but good to keep this check
	block := cs.votes.lockedProposal.Block()
	if err := cs.state.ValidateBlock(block); err != nil {
		cs.votes.lockedProposal = nil
		cs.logger.Error("Commit: invalid block", "block", block, "err", err)
		return
	}

	commit := preCommits.ToCommit()
	if commit != nil {
		if err := cs.state.ApplyBlock(block, *commit); err != nil {
			cs.logger.Error("Commit: Applying block failed", "block", block, "err", err)
			return
		}

		// Npw broadcast the committed block
		msg := message.NewBlockMessage(height, block, *commit)
		cs.broadcastCh <- msg
	}

	cs.updateRoundStep(round, hrs.StepTypeCommit)
	cs.isCommitted = true

	cs.logger.Info("Commit: Block stored", "block", blockHash.Fingerprint())
	cs.scheduleNewHeight()

}
