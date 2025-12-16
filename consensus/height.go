package consensus

import (
	"time"

	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type newHeightState struct {
	*consensus
}

func (s *newHeightState) enter() {
	s.decide()
}

func (s *newHeightState) decide() {
	sateHeight := s.bcState.LastBlockHeight()
	validators := s.bcState.CommitteeValidators()
	s.log.MoveToNewHeight(validators)

	s.validators = validators
	s.height = sateHeight + 1
	s.round = 0
	s.active = s.bcState.IsInCommittee(s.valKey.Address())
	s.logger.Info("entering new height", "height", s.height, "active", s.active)

	sleep := time.Until(s.bcState.LastBlockTime().Add(s.bcState.Params().BlockInterval()))
	s.scheduleTimeout(sleep, s.height, s.round, tickerTargetNewHeight)
}

func (s *newHeightState) onAddVote(_ *vote.Vote) {
	prepares := s.log.PrepareVoteSet(s.round)
	if prepares.HasQuorumHash() {
		// Add logic to detect when the network majority has voted for a block,
		// but the new height timer has not yet started. This situation can occur if the system
		// time is lagging behind the network time.
		s.logger.Warn("detected network majority voting for a block, but the new height timer has not started yet. " +
			"system time may be behind the network.")
		s.enterNewState(s.proposeState)
	}
}

func (*newHeightState) onSetProposal(_ *proposal.Proposal) {
	// Ignore proposal
}

func (s *newHeightState) onTimeout(t *ticker) {
	if t.Target == tickerTargetNewHeight {
		if s.active {
			s.enterNewState(s.proposeState)
		}
	}
}

func (*newHeightState) name() string {
	return "new-height"
}
