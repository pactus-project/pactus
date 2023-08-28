package consensus

import (
	"github.com/pactus-project/pactus/types/vote"
)

type cpDecideState struct {
	*consensus
}

func (s *cpDecideState) enter() {
	s.decide()
}

func (s *cpDecideState) decide() {
	if s.cpDecided == 1 {
		s.round++
		s.enterNewState(s.proposeState)
	} else if s.cpDecided == 0 {
		s.queryProposal()
		// s.queryVotes()
		s.enterNewState(s.prepareState)
	} else {
		cpMainVotes := s.log.CPMainVoteVoteSet(s.round)
		if cpMainVotes.HasTwoThirdOfTotalPower(s.cpRound) {
			if cpMainVotes.HasQuorumVotesFor(s.cpRound, vote.CPValueOne) {
				// decided for yes, and proceeds to the next round
				s.logger.Info("binary agreement decided", "value", 1, "round", s.cpRound)

				s.cpDecided = 1
			} else if cpMainVotes.HasQuorumVotesFor(s.cpRound, vote.CPValueZero) {
				// decided for no and proceeds to the next round
				s.logger.Info("binary agreement decided", "value", 0, "round", s.cpRound)

				s.cpDecided = 0
			} else {
				// conflicting votes
				s.logger.Info("conflicting main votes", "round", s.cpRound)
			}

			s.cpRound++
			s.enterNewState(s.cpPreVoteState)
		}
	}
}

func (s *cpDecideState) timeout(_ *ticker) {
	// Ignore timeouts
}

func (s *cpDecideState) name() string {
	return "cp:decide"
}
