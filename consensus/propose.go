package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
)

type proposeState struct {
	*consensus
}

func (s *proposeState) enter() {
	s.decide()
}

func (s *proposeState) decide() {
	proposer := s.proposer(s.round)
	if proposer.Address().EqualsTo(s.signer.Address()) {
		s.logger.Info("Our turn to propose", "proposer", proposer.Address())
		s.createProposal(s.height, s.round)
	} else {
		s.logger.Debug("Not our turn to propose", "proposer", proposer.Address())
	}

	s.enterNewState(s.prepareState)
}

func (s *proposeState) createProposal(height int, round int) {
	block, err := s.state.ProposeBlock(round)
	if err != nil {
		s.logger.Error("We can't propose a block. Why?", "err", err)
		return
	}
	if err := s.state.ValidateBlock(block); err != nil {
		s.logger.Error("Our block is invalid. Why?", "err", err)
		return
	}

	proposal := proposal.NewProposal(height, round, block)
	s.signer.SignMsg(proposal)
	s.doSetProposal(proposal)

	s.logger.Info("Proposal signed and broadcasted", "proposal", proposal)

	s.broadcastProposal(proposal)
}

func (s *proposeState) onAddVote(v *vote.Vote) {
	panic("Unreachable")
}

func (s *proposeState) onSetProposal(p *proposal.Proposal) {
	panic("Unreachable")
}

func (s *proposeState) onTimedout(t *ticker) {
	panic("Unreachable")
}
func (s *proposeState) name() string {
	return "propose"
}
