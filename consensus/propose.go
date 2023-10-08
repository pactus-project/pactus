package consensus

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type proposeState struct {
	*consensus
}

func (s *proposeState) enter() {
	s.decide()
}

func (s *proposeState) decide() {
	proposer := s.proposer(s.round)
	if proposer.Address() == s.valKey.Address() {
		s.logger.Info("our turn to propose", "proposer", proposer.Address())
		s.createProposal(s.height, s.round)
	} else {
		s.logger.Debug("not our turn to propose", "proposer", proposer.Address())
	}

	s.cpRound = 0
	s.cpDecided = -1
	s.cpWeakValidity = nil
	s.enterNewState(s.prepareState)
}

func (s *proposeState) createProposal(height uint32, round int16) {
	block, err := s.state.ProposeBlock(s.valKey, s.rewardAddr)
	if err != nil {
		s.logger.Error("unable to propose a block!", "error", err)
		return
	}
	if err := s.state.ValidateBlock(block); err != nil {
		s.logger.Error("proposed block is invalid!", "error", err)
		return
	}

	proposal := proposal.NewProposal(height, round, block)
	sig := s.valKey.Sign(proposal.SignBytes())
	proposal.SetSignature(sig)

	s.log.SetRoundProposal(round, proposal)

	s.broadcastProposal(proposal)

	s.logger.Info("proposal signed and broadcasted", "proposal", proposal)
}

func (s *proposeState) onAddVote(_ *vote.Vote) {
	panic("Unreachable")
}

func (s *proposeState) onSetProposal(_ *proposal.Proposal) {
	panic("Unreachable")
}

func (s *proposeState) onTimeout(_ *ticker) {
	panic("Unreachable")
}

func (s *proposeState) name() string {
	return "propose"
}
