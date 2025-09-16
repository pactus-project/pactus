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
		if cpMainVotes.Has2FP1VotesFor(s.cpRound, vote.CPValueNo) {
			s.logger.Panic("unreachable state: decide on 'no' should be handled in main-vote state")
		} else if cpMainVotes.Has2FP1VotesFor(s.cpRound, vote.CPValueYes) {
			// decided for yes, and proceeds to the next round
			s.logger.Info("binary agreement decided", "value", "yes", "round", s.cpRound)

			votes := cpMainVotes.BinaryVotes(s.cpRound, vote.CPValueYes)
			cert := s.makeCertificate(votes)
			just := &vote.JustDecided{
				QCert: cert,
			}
			s.signAddCPDecidedVote(hash.UndefHash, s.cpRound, vote.CPValueYes, just)
		} else {
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
