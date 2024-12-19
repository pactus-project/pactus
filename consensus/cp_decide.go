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
		if cpMainVotes.HasQuorumVotesFor(s.cpRound, vote.CPValueYes) {
			// decided for yes, and proceeds to the next round
			s.logger.Info("binary agreement decided", "value", 1, "round", s.cpRound)

			votes := cpMainVotes.BinaryVotes(s.cpRound, vote.CPValueYes)
			cert := s.makeVoteCertificate(votes)
			just := &vote.JustDecided{
				QCert: cert,
			}
			s.signAddCPDecidedVote(hash.UndefHash, s.cpRound, vote.CPValueYes, just)
			s.cpStrongTermination(s.round, s.cpRound)
		} else if cpMainVotes.HasQuorumVotesFor(s.cpRound, vote.CPValueNo) {
			// decided for no and proceeds to the next round
			s.logger.Info("binary agreement decided", "value", 0, "round", s.cpRound)

			votes := cpMainVotes.BinaryVotes(s.cpRound, vote.CPValueNo)
			cert := s.makeVoteCertificate(votes)
			just := &vote.JustDecided{
				QCert: cert,
			}
			s.signAddCPDecidedVote(*s.cpWeakValidity, s.cpRound, vote.CPValueNo, just)
			s.cpStrongTermination(s.round, s.cpRound)
		} else {
			// conflicting votes
			s.logger.Debug("conflicting main votes", "round", s.cpRound)
			s.cpRound++
			s.enterNewState(s.cpPreVoteState)
		}
	}
}

func (s *cpDecideState) onAddVote(vte *vote.Vote) {
	if vte.Type() == vote.VoteTypeCPMainVote {
		s.decide()
	}

	if vte.IsCPVote() {
		s.cpStrongTermination(vte.Round(), vte.CPRound())
	}
}

func (*cpDecideState) name() string {
	return "cp:decide"
}
