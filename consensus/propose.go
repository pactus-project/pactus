package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
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
	if cs.hrs.InvalidHeightRound(proposal.Height(), proposal.Round()) {
		cs.logger.Error("Proposal received from wrong height/round", "proposal", proposal)
		return
	}

	roundProposal := cs.votes.RoundProposal(cs.hrs.Round())
	if roundProposal != nil {
		cs.logger.Trace("We already have a proposal for this round", "proposal", proposal)
		return
	}

	if err := proposal.SanityCheck(); err != nil {
		cs.logger.Error("Proposal is invalid", "proposal", proposal)
		return
	}

	if err := proposal.Verify(cs.proposer().PublicKey()); err != nil {
		cs.logger.Error("Proposal has invalid signature", "proposal", proposal, "err", err)
		return
	}

	cs.state.SyncTxPool(proposal.Block())

	cs.votes.SetRoundProposal(cs.hrs.Round(), proposal)
	// Maybe received proposal after prevote, (Due to network latency maybe?)
	// Enter prevote
	cs.enterPrevote(proposal.Height(), proposal.Round())
}

func (cs *Consensus) enterPropose(height int, round int) {
	if cs.hrs.InvalidHeightRoundStep(height, round, hrs.StepTypePropose) {
		cs.logger.Debug("Propose with invalid args", "height", height, "round", round)
		return
	}

	if cs.privValidator == nil {
		cs.logger.Debug("This node is not a validator")
		return
	}
	cs.logger.Debug("This node is a validator")

	address := cs.privValidator.Address()
	if !cs.valset.Contains(address) {
		cs.logger.Error("This node is not in validator set", "addr", address)
		return
	}

	cs.updateRoundStep(round, hrs.StepTypePropose)

	if cs.isProposer(address) {
		cs.logger.Info("Our turn to propose", "proposer", address)
		cs.createProposal(height, round)
	} else {
		cs.logger.Debug("Not our turn to propose", "proposer", cs.proposer().Address())
	}

	cs.scheduleTimeout(cs.config.Consensus.Propose(round), height, round, hrs.StepTypePrevote)
}

func (cs *Consensus) createProposal(height int, round int) {
	if cs.privValidator == nil {
		cs.logger.Error("privValidator is nil")
		return
	}
	if height > 1 && cs.lastCommit == nil {
		cs.logger.Error("We don't have lastCommit, Unsafe restart?")
		return
	}

	proposerAddr := cs.privValidator.Address()
	block := cs.state.ProposeBlock(height, proposerAddr, cs.lastCommit)

	if err := cs.state.ValidateBlock(&block); err != nil {
		cs.logger.Error("Our block is invalid. Why?")
		return
	}

	proposal := vote.NewProposal(height, round, block)
	cs.privValidator.SignMsg(proposal)
	cs.setProposal(proposal)

	cs.logger.Info("Proposal signed and set", "proposal", proposal)

	// Broadcast proposal
	go cs.syncer.BroadcastProposal(proposal)

}
