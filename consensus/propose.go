package consensus

import (
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/vote"
)

type proposeState struct {
	*consensus
}

func (s *proposeState) enter() {
	s.execute()
}

func (s *proposeState) execute() {
	proposer := s.proposer(s.round)
	if proposer.Address().EqualsTo(s.signer.Address()) {
		s.logger.Info("Our turn to propose", "proposer", proposer.Address())
		s.createProposal(s.height, s.round)
	} else {
		s.queryProposal()
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
	if err := s.state.ValidateBlock(*block); err != nil {
		s.logger.Error("Our block is invalid. Why?", "err", err)
		return
	}

	proposal := proposal.NewProposal(height, round, *block)
	s.signer.SignMsg(proposal)
	s.setProposal(proposal)

	s.logger.Info("Proposal signed and broadcasted", "proposal", proposal)

	s.broadcastProposal(proposal)
}

func (s *proposeState) timedout(t *ticker) {
	s.execute()
}

func (s *proposeState) voteAdded(v *vote.Vote) {
}

func (s *proposeState) name() string {
	return proposeName
}
