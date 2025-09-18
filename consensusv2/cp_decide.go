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
	cpMainVotes := s.log.CPMainVoteVoteSet(s.round)
	if cpMainVotes.Has2FP1Votes(s.cpRound) {
		switch {
		case cpMainVotes.Has2FP1VotesFor(s.cpRound, vote.CPValueNo):
			s.logger.Panic("unreachable state: decide on 'no (biased)' should be handled in pre-vote state")

		case cpMainVotes.Has2FP1VotesFor(s.cpRound, vote.CPValueYes):
			s.logger.Info("binary agreement decided", "value", "yes", "round", s.cpRound)

			// decided for yes, and proceeds to the next round
			votes := cpMainVotes.BinaryVotes(s.cpRound, vote.CPValueYes)
			cert := s.makeCertificate(votes)
			just := &vote.JustDecided{
				QCert: cert,
			}
			s.signAddCPDecidedVote(hash.UndefHash, s.cpRound, vote.CPValueYes, just)

		default:
			// conflicting votes
			s.logger.Debug("conflicting main votes", "round", s.cpRound)
			s.cpRound++
			s.enterNewState(s.cpPreVoteState)
		}
	}

	s.cpStrongTermination()
	s.absoluteCommit()
}

func (s *cpDecideState) onAddVote(_ *vote.Vote) {
	s.decide()
}

func (*cpDecideState) name() string {
	return "cp:decide"
}
