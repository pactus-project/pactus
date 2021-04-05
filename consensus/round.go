package consensus

import (
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/vote"
)

type newRoundState struct {
	*consensus
}

func (s *newRoundState) enter() {
	sleep := s.config.ChangeProposerTimeout
	s.scheduleTimeout(sleep, s.height, s.round, tickerTargetChangeProposer)
	s.logger.Debug("Change proposer timer started...", "timeout", sleep.Seconds())

	s.decide()
}

func (s *newRoundState) decide() {
	s.logger.Info("Entering new round", "round", s.round)
	s.enterNewState(s.proposeState)
}

func (s *newRoundState) onAddVote(v *vote.Vote) {
	panic("Unreachable")
}

func (s *newRoundState) onSetProposal(p *proposal.Proposal) {
	panic("Unreachable")
}

func (s *newRoundState) onTimedout(t *ticker) {
	panic("Unreachable")
}

func (s *newRoundState) name() string {
	return "new-round"
}
