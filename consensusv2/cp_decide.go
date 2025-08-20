package consensusv2

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
)

type cpDecideState struct {
	*changeProposer
}

func (s *cpDecideState) enter() {
	s.decide()
}

func (s *cpDecideState) decide() {
	s.absoluteCommit()
	s.cpStrongTermination()

	cpMainVotes := s.log.CPMainVoteVoteSet(s.round)
	if cpMainVotes.HasTwoFPlusOneVotes(s.cpRound) {
		if cpMainVotes.HasTwoFPlusOneVotesFor(s.cpRound, vote.CPValueNo) {
			panic("unreachable")
		} else if cpMainVotes.HasTwoFPlusOneVotesFor(s.cpRound, vote.CPValueYes) {
			// decided for yes, and proceeds to the next round
			s.logger.Info("binary agreement decided", "value", "yes", "round", s.cpRound)

			votes := cpMainVotes.BinaryVotes(s.cpRound, vote.CPValueYes)
			cert := s.makeVoteCertificate(votes)
			just := &vote.JustDecided{
				QCert: cert,
			}
			s.signAddCPDecidedVote(hash.UndefHash, s.cpRound, vote.CPValueYes, just)
			s.cpStrongTermination()
		} else {
			// conflicting votes
			s.logger.Debug("conflicting main votes", "round", s.cpRound)
			s.cpRound++
			s.enterNewState(s.cpPreVoteState)
		}
	}
}

func (s *cpDecideState) onAddVote(_ *vote.Vote) {
	s.decide()
}

func (*cpDecideState) name() string {
	return "cp:decide"
}
