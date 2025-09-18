package consensusv2

import (
	"time"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type newHeightState struct {
	*consensusV2
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
	s.cpRound = 0
	s.cpDecidedCert = nil
	s.cpWeakValidity = hash.UndefHash
	s.active = s.bcState.IsInCommittee(s.valKey.Address())
	s.logger.Info("entering new height", "height", s.height, "active", s.active)

	sleep := time.Until(s.bcState.LastBlockTime().Add(s.bcState.Params().BlockInterval()))
	s.scheduleTimeout(sleep, s.height, s.round, tickerTargetNewHeight)
}

func (s *newHeightState) onAddVote(_ *vote.Vote) {
	precommits := s.log.PrecommitVoteSet(s.round)
	if precommits.Has2FP1Votes() {
		// Detect when the network majority has voted (2f+1 precommits) but the new height
		// timer hasn't started yet. This edge case occurs when the local system time is
		// lagging behind the network time, causing the scheduled timeout to be delayed.
		// In this scenario, we should immediately transition to the propose state rather
		// than waiting for the timer to expire.
		s.logger.Warn("network majority has voted but new-height timer has not started yet. " +
			"System time may be lagging behind network time.")
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
