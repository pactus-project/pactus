package consensus

import (
	"github.com/zarbchain/zarb-go/types/proposal"
	"github.com/zarbchain/zarb-go/types/vote"
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
		s.logger.Info("our turn to propose", "proposer", proposer.Address())
		s.createProposal(s.height, s.round)
	} else {
		s.logger.Debug("not our turn to propose", "proposer", proposer.Address())
	}

	s.enterNewState(s.prepareState)
}

func (s *proposeState) createProposal(height int32, round int16) {
	block, err := s.state.ProposeBlock(round)
	if err != nil {
		s.logger.Error("unable to propose a block. Why?", "err", err)
		return
	}
	if err := s.state.ValidateBlock(block); err != nil {
		s.logger.Error("proposed block is invalid. Why?", "err", err)
		return
	}

	proposal := proposal.NewProposal(height, round, block)
	s.signer.SignMsg(proposal)
	s.doSetProposal(proposal)

	s.logger.Info("proposal signed and broadcasted", "proposal", proposal)

	s.broadcastProposal(proposal)
}

func (s *proposeState) onAddVote(v *vote.Vote) {
	panic("Unreachable")
}

func (s *proposeState) onSetProposal(p *proposal.Proposal) {
	panic("Unreachable")
}

func (s *proposeState) onTimeout(t *ticker) {
	panic("Unreachable")
}
func (s *proposeState) name() string {
	return "propose"
}
