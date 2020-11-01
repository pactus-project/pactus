package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

func (cs *Consensus) proposer() *validator.Validator {
	return cs.valset.Proposer(cs.hrs.Round())
}

func (cs *Consensus) isProposer(address crypto.Address) bool {
	return cs.proposer().Address().EqualsTo(address)
}

func (cs *Consensus) setProposal(proposal *vote.Proposal) {
	if cs.invalidHeightRound(proposal.Height(), proposal.Round()) {
		cs.logger.Info("Proposal received from wrong height/round", "proposal", proposal)
		return
	}

	roundProposal := cs.votes.RoundProposal(cs.hrs.Round())
	if roundProposal != nil {
		cs.logger.Trace("propose: This round has proposal", "proposal", proposal)
		return
	}

	if err := proposal.SanityCheck(); err != nil {
		cs.logger.Error("propose: Proposal is invalid", "proposal", proposal, "err", err)
		return
	}

	if err := proposal.Verify(cs.proposer().PublicKey()); err != nil {
		cs.logger.Error("propose: Proposal has invalid signature", "proposal", proposal, "err", err)
		return
	}

	cs.logger.Info("propose: Proposal set", "proposal", proposal)
	cs.votes.SetRoundProposal(cs.hrs.Round(), proposal)
	// Maybe received proposal after prevote, (maybe because of network latency?)
	// Enter prevote
	cs.enterPrevote(proposal.Height(), proposal.Round())
}

func (cs *Consensus) enterPropose(height int, round int) {
	if cs.invalidHeightRoundStep(height, round, hrs.StepTypePropose) {
		cs.logger.Debug("Propose with invalid args", "height", height, "round", round)
		return
	}

	if cs.privValidator == nil {
		cs.logger.Debug("Propose: This node is not a validator")
		return
	}
	cs.logger.Debug("Propose: This node is a validator")

	address := cs.privValidator.Address()
	if !cs.valset.Contains(address) {
		cs.logger.Error("Propose: This node is not in validator set", "addr", address)
		return
	}

	cs.updateRoundStep(round, hrs.StepTypePropose)

	if cs.isProposer(address) {
		cs.logger.Info("Propose: Our turn to propose", "proposer", address)
		cs.createProposal(height, round)
	} else {
		cs.logger.Debug("Propose: Not our turn to propose", "proposer", cs.proposer().Address())
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
