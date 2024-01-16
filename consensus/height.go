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
	sateHeight := s.bcState.LastBlockHeight()
	validators := s.bcState.CommitteeValidators()
	s.log.MoveToNewHeight(validators)

	s.validators = validators
	s.height = sateHeight + 1
	s.round = 0
	s.active = s.bcState.IsInCommittee(s.valKey.Address())
	s.logger.Info("entering new height", "height", s.height, "active", s.active)

	if s.active {
		sleep := s.bcState.LastBlockTime().Add(s.bcState.Params().BlockInterval()).Sub(util.Now())
		s.scheduleTimeout(sleep, s.height, s.round, tickerTargetNewHeight)
	}
}

func (s *newHeightState) onAddVote(_ *vote.Vote) {
	prepares := s.log.PrepareVoteSet(s.round)
	if prepares.HasQuorumHash() {
		// Added logic to detect when the network majority has voted for a block,
		// but the new height timer has not yet started. This situation can occur if the system
		// time is lagging behind the network time.
		s.logger.Warn("detected network majority voting for a block, but the new height timer has not started yet. " +
			"System time may be behind the network.")
		s.enterNewState(s.proposeState)
	}
}

func (s *newHeightState) onSetProposal(_ *proposal.Proposal) {
	// Ignore proposal
}

func (s *newHeightState) onTimeout(t *ticker) {
	if t.Target == tickerTargetNewHeight {
		s.enterNewState(s.proposeState)
	}
}

func (s *newHeightState) name() string {
	return "new-height"
}
