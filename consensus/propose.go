package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *Consensus) proposer(round int) *validator.Validator {
	return cs.valset.Proposer(round)
}

func (cs *Consensus) isProposer(address crypto.Address, round int) bool {
	return cs.proposer(round).Address().EqualsTo(address)
}

func (cs *Consensus) setProposal(proposal *vote.Proposal) {
	if cs.invalidHeight(proposal.Height()) {
		cs.logger.Debug("Propose: Invalid height or committed", "proposal", proposal, "committed", cs.isCommitted)
		return
	}

	roundProposal := cs.votes.RoundProposal(proposal.Round())
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

	cs.logger.Info("propose: Proposal set", "proposal", proposal)
	cs.votes.SetRoundProposal(proposal.Round(), proposal)
	// Maybe received proposal after prevote, (maybe because of network latency?)
	// Enter prevote
	cs.enterPrevote(proposal.Height(), proposal.Round())
}

func (cs *Consensus) enterPropose(height int, round int) {
	if cs.invalidHeightRoundStep(height, round, hrs.StepTypePropose) {
		cs.logger.Debug("Propose: Invalid height/round/step or committed before", "height", height, "round", round, "committed", cs.isCommitted)
		return
	}

	if cs.privValidator == nil {
		cs.logger.Debug("Propose: This node is not a validator")
		return
	}
	cs.logger.Debug("Propose: This node is a validator")

	address := cs.privValidator.Address()
	if !cs.valset.Contains(address) {
		cs.logger.Info("Propose: This node is not in validator set", "addr", address)
		return
	}

	cs.updateRoundStep(round, hrs.StepTypePropose)

	if cs.isProposer(address, round) {
		cs.logger.Info("Propose: Our turn to propose", "proposer", address)
		cs.createProposal(height, round)
	} else {
		cs.logger.Debug("Propose: Not our turn to propose", "proposer", cs.proposer(round).Address())
	}

	cs.scheduleTimeout(cs.config.Propose(round), height, round, hrs.StepTypePrevote)
}

func (cs *Consensus) createProposal(height int, round int) {
	if cs.privValidator == nil {
		cs.logger.Error("Propose: privValidator is nil")
		return
	}

	block := cs.state.ProposeBlock()
	if err := cs.state.ValidateBlock(block); err != nil {
		cs.logger.Error("Propose: Our block is invalid. Why?", "error", err)
		return
	}

	proposal := vote.NewProposal(height, round, block)
	if cs.config.FuzzTesting {
		if n := util.RandInt(5); n == 3 {
			// Randomly send invalid proposal
			proposal, _ = vote.GenerateTestProposal(cs.hrs.Height(), cs.hrs.Round())
		}
	}
	cs.privValidator.SignMsg(proposal)
	cs.setProposal(proposal)

	cs.logger.Info("Propose: Proposal signed and sent", "proposal", proposal)

	// Broadcast proposal
	msg := message.NewProposalMessage(*proposal)
	cs.broadcastCh <- msg
}
