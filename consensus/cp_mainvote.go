package consensus

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type cpMainVoteState struct {
	*consensus
}

func (s *cpMainVoteState) enter() {
	s.decide()
}

func (s *cpMainVoteState) decide() {
	cpPreVotes := s.log.CPPreVoteVoteSet(s.round)
	if cpPreVotes.HasTwoThirdOfTotalPower(s.cpRound) {
		if cpPreVotes.HasQuorumVotesFor(s.cpRound, vote.CPValueOne) {
			s.logger.Info("cp: quorum for pre-votes", "v", "1")

			votes := cpPreVotes.BinaryVotes(s.cpRound, vote.CPValueOne)
			cert := s.makeCertificate(votes)
			just := &vote.JustMainVoteNoConflict{
				QCert: cert,
			}
			s.signAddCPMainVote(hash.UndefHash, s.cpRound, vote.CPValueOne, just)
			s.enterNewState(s.cpDecideState)
		} else if cpPreVotes.HasQuorumVotesFor(s.cpRound, vote.CPValueZero) {
			s.logger.Info("cp: quorum for pre-votes", "v", "0")

			votes := cpPreVotes.BinaryVotes(s.cpRound, vote.CPValueZero)
			cert := s.makeCertificate(votes)
			just := &vote.JustMainVoteNoConflict{
				QCert: cert,
			}
			s.signAddCPMainVote(*s.cpWeakValidity, s.cpRound, vote.CPValueZero, just)
			s.enterNewState(s.cpDecideState)
		} else {
			s.logger.Info("cp: no-quorum for pre-votes", "v", "abstain")

			vote0 := cpPreVotes.GetRandomVote(s.cpRound, vote.CPValueZero)
			vote1 := cpPreVotes.GetRandomVote(s.cpRound, vote.CPValueOne)

			just := &vote.JustMainVoteConflict{
				Just0: vote0.CPJust(),
				Just1: vote1.CPJust(),
			}
			s.signAddCPMainVote(*s.cpWeakValidity, s.cpRound, vote.CPValueAbstain, just)
			s.enterNewState(s.cpDecideState)
		}
	}
}

func (s *cpMainVoteState) onAddVote(v *vote.Vote) {
	if v.Type() == vote.VoteTypeCPPreVote {
		s.decide()
	}
}

func (s *cpMainVoteState) onSetProposal(_ *proposal.Proposal) {
	// Ignore proposal
}

func (s *cpMainVoteState) onTimeout(_ *ticker) {
	// Ignore timeouts
}

func (s *cpMainVoteState) name() string {
	return "cp:main-vote"
}
