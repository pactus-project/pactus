package consensus

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
)

type cpPreVoteState struct {
	*changeProposer
}

func (s *cpPreVoteState) enter() {
	s.decide()
}

func (s *cpPreVoteState) decide() {
	if s.cpRound == 0 {
		// broadcast the initial value
		prepares := s.log.PrepareVoteSet(s.round)
		preparesQH := prepares.QuorumHash()
		if preparesQH != nil {
			s.cpWeakValidity = preparesQH
			cert := s.makeCertificate(prepares.BlockVotes(*preparesQH))
			just := &vote.JustInitNo{
				QCert: cert,
			}
			s.signAddCPPreVote(*s.cpWeakValidity, s.cpRound, 0, just)
		} else {
			just := &vote.JustInitYes{}
			s.signAddCPPreVote(hash.UndefHash, s.cpRound, 1, just)
		}
		s.scheduleTimeout(s.config.QueryVoteTimeout, s.height, s.round, tickerTargetQueryVote)
	} else {
		cpMainVotes := s.log.CPMainVoteVoteSet(s.round)
		switch {
		case cpMainVotes.HasAnyVoteFor(s.cpRound-1, vote.CPValueYes):
			s.logger.Debug("cp: one main-vote for one", "b", "1")

			vote1 := cpMainVotes.GetRandomVote(s.cpRound-1, vote.CPValueYes)
			just1 := &vote.JustPreVoteHard{
				QCert: vote1.CPJust().(*vote.JustMainVoteNoConflict).QCert,
			}
			s.signAddCPPreVote(hash.UndefHash, s.cpRound, vote.CPValueYes, just1)

		case cpMainVotes.HasAnyVoteFor(s.cpRound-1, vote.CPValueNo):
			s.logger.Debug("cp: one main-vote for zero", "b", "0")

			vote0 := cpMainVotes.GetRandomVote(s.cpRound-1, vote.CPValueNo)
			just0 := &vote.JustPreVoteHard{
				QCert: vote0.CPJust().(*vote.JustMainVoteNoConflict).QCert,
			}
			s.signAddCPPreVote(*s.cpWeakValidity, s.cpRound, vote.CPValueNo, just0)

		case cpMainVotes.HasAllVotesFor(s.cpRound-1, vote.CPValueAbstain):
			s.logger.Debug("cp: all main-votes are abstain", "b", "0 (biased)")

			votes := cpMainVotes.BinaryVotes(s.cpRound-1, vote.CPValueAbstain)
			cert := s.makeCertificate(votes)
			just := &vote.JustPreVoteSoft{
				QCert: cert,
			}
			s.signAddCPPreVote(*s.cpWeakValidity, s.cpRound, vote.CPValueNo, just)

		default:
			s.logger.Panic("protocol violated. We have combination of votes for one and zero")
		}
	}

	s.enterNewState(s.cpMainVoteState)
}

func (*cpPreVoteState) onAddVote(_ *vote.Vote) {
	panic("Unreachable")
}

func (*cpPreVoteState) name() string {
	return "cp:pre-vote"
}
