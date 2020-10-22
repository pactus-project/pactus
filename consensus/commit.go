package consensus

import "github.com/zarbchain/zarb-go/consensus/hrs"

func (cs *Consensus) enterCommit(height int, round int) {
	if cs.hrs.InvalidHeight(height) || cs.commitRound != -1 {
		cs.logger.Debug("Commit with invalid args or committed before", "height", height)
		return
	}

	if !cs.votes.Precommits(round).HasQuorum() {
		cs.logger.Error("Commit witout quorom for precommit stage")
		return
	}

	blockHash := cs.votes.Precommits(round).QuorumBlock()
	if blockHash == nil || blockHash.IsUndef() {
		cs.logger.Error("Commit is for invalid block")
		return
	}

	// Additional check. blockHash should be same for both prevotes and precommits
	prevoteBlockHash := cs.votes.Prevotes(round).QuorumBlock()
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
	if err := cs.state.ValidateBlock(roundProposal.Block()); err != nil {
		cs.votes.lockedProposal = nil
		cs.logger.Warn("Commit: invalid block", "Proposal", roundProposal, "err", err)
		return
	}

	block := cs.votes.lockedProposal.Block()
	cs.state.ApplyBlock(block, round)
	cs.updateRoundStep(round, hrs.StepTypeCommit)
	cs.commitRound = round

	// Using `~` to show `block` after `consensus` in logger output
	cs.logger.Info("Commit: Block stored", "~block", blockHash.Fingerprint())

	cs.scheduleNewHeight()
}
