package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *consensus) proposer(round int) *validator.Validator {
	return cs.state.ValidatorSet().Proposer(round)
}

func (cs *consensus) isProposer(addr crypto.Address, round int) bool {
	return cs.state.ValidatorSet().IsProposer(addr, round)
}

func (cs *consensus) setProposal(proposal *vote.Proposal) {
	if proposal.Height() != cs.hrs.Height() {
		cs.logger.Debug("Propose: Invalid height", "proposal", proposal)
		return
	}
	if proposal.Round() > cs.hrs.Round() {
		cs.logger.Debug("Propose: Invalid round", "proposal", proposal)
		return
	}

	roundProposal := cs.pendingVotes.RoundProposal(proposal.Round())
	if roundProposal != nil {
		cs.logger.Trace("propose: This round has proposal", "proposal", proposal)
		return
	}

	if err := proposal.SanityCheck(); err != nil {
		cs.logger.Error("propose: Proposal is invalid", "proposal", proposal, "err", err)
		return
	}

	proposer := cs.proposer(proposal.Round())
	if err := proposal.Verify(proposer.PublicKey()); err != nil {
		cs.logger.Error("propose: Proposal has invalid signature", "proposal", proposal, "err", err)
		return
	}

	if err := cs.state.ValidateBlock(proposal.Block()); err != nil {
		cs.logger.Warn("propose: Invalid block", "proposal", proposal, "err", err)
		return
	}

	cs.logger.Info("propose: Proposal set", "proposal", proposal)
	cs.pendingVotes.SetRoundProposal(proposal.Round(), proposal)
	// Proposal might be received after prepare or precommit, (maybe because of network latency?)
	cs.enterPrepare(proposal.Round())
	cs.enterPrecommit(proposal.Round())
}

func (cs *consensus) enterPropose(round int) {
	if cs.isProposed || round != cs.hrs.Round() {
		cs.logger.Debug("Propose: Proposed before or invalid round", "round", round)
		return
	}
	cs.updateStep(hrs.StepTypePropose)
	cs.scheduleTimeout(cs.config.PrepareTimeout(round), cs.hrs.Height(), round, hrs.StepTypePrepare)

	address := cs.signer.Address()
	if !cs.state.ValidatorSet().Contains(address) {
		cs.logger.Debug("Propose: This node is not in validator set", "addr", address)
		return
	}

	if cs.isProposer(address, round) {
		cs.logger.Info("Propose: Our turn to propose", "proposer", address)
		cs.createProposal(cs.hrs.Height(), round)
	} else {
		cs.logger.Debug("Propose: Not our turn to propose", "proposer", cs.proposer(round).Address())
	}

	cs.isProposed = true
}

func (cs *consensus) createProposal(height int, round int) {
	block, err := cs.state.ProposeBlock(round)
	if err != nil {
		cs.logger.Error("Propose: We can't propose a block. Why?", "err", err)
		return
	}
	if err := cs.state.ValidateBlock(*block); err != nil {
		cs.logger.Error("Propose: Our block is invalid. Why?", "err", err)
		return
	}

	proposal := vote.NewProposal(height, round, *block)
	cs.signer.SignMsg(proposal)
	cs.setProposal(proposal)

	cs.logger.Info("Proposal signed and broadcasted", "proposal", proposal)

	cs.broadcastProposal(proposal)
}
