package consensus

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
	cpMainVotes := s.log.CPMainVoteVoteSet(s.round)
	if cpMainVotes.HasTwoThirdOfTotalPower(s.cpRound) {
		if cpMainVotes.HasQuorumVotesFor(s.cpRound, vote.CPValueOne) {
			// decided for yes, and proceeds to the next round
			s.logger.Info("binary agreement decided", "value", 1, "round", s.cpRound)

			votes := cpMainVotes.BinaryVotes(s.cpRound, vote.CPValueOne)
			cert := s.makeCertificate(votes)
			just := &vote.JustDecided{
				QCert: cert,
			}
			s.signAddCPDecidedVote(hash.UndefHash, s.cpRound, vote.CPValueOne, just)
			s.cpDecide(vote.CPValueOne)
		} else if cpMainVotes.HasQuorumVotesFor(s.cpRound, vote.CPValueZero) {
			// decided for no and proceeds to the next round
			s.logger.Info("binary agreement decided", "value", 0, "round", s.cpRound)

			votes := cpMainVotes.BinaryVotes(s.cpRound, vote.CPValueZero)
			cert := s.makeCertificate(votes)
			just := &vote.JustDecided{
				QCert: cert,
			}
			s.signAddCPDecidedVote(*s.cpWeakValidity, s.cpRound, vote.CPValueZero, just)
			s.cpDecide(vote.CPValueZero)
		} else {
			// conflicting votes
			s.logger.Debug("conflicting main votes", "round", s.cpRound)
			s.cpRound++
			s.enterNewState(s.cpPreVoteState)
		}
	}
}

func (s *cpDecideState) onAddVote(v *vote.Vote) {
	if v.Type() == vote.VoteTypeCPMainVote {
		s.decide()
	}

	s.strongTermination()
}

func (*cpDecideState) name() string {
	return "cp:decide"
}
