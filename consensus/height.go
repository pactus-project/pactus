package consensus

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
)

type newHeightState struct {
	*consensus
}

func (s *newHeightState) enter() {
	s.decide()
}

func (s *newHeightState) decide() {
	sateHeight := s.state.LastBlockHeight()
	validators := s.state.CommitteeValidators()
	s.log.MoveToNewHeight(validators)

	s.validators = validators
	s.height = sateHeight + 1
	s.round = 0
	s.active = s.state.IsInCommittee(s.valKey.Address())
	s.logger.Info("entering new height", "height", s.height, "active", s.active)

	sleep := s.state.LastBlockTime().Add(s.state.Params().BlockInterval()).Sub(util.Now())
	s.scheduleTimeout(sleep, s.height, s.round, tickerTargetNewHeight)
}

func (s *newHeightState) onAddVote(_ *vote.Vote) {
	// Ignore votes
}

func (s *newHeightState) onSetProposal(_ *proposal.Proposal) {
	// Ignore proposal
}

func (s *newHeightState) onTimeout(t *ticker) {
	if t.Target == tickerTargetNewHeight {
		if s.active {
			s.enterNewState(s.proposeState)
		}
	}
}

func (s *newHeightState) name() string {
	return "new-height"
}
