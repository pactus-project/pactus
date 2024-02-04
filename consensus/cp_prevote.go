package consensus

import (
	"time"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
)

var queryVoteInitialTimeout = 2 * time.Second

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
			just := &vote.JustInitZero{
				QCert: cert,
			}
			s.signAddCPPreVote(*s.cpWeakValidity, s.cpRound, 0, just)
		} else {
			just := &vote.JustInitOne{}
			s.signAddCPPreVote(hash.UndefHash, s.cpRound, 1, just)
		}
		s.scheduleTimeout(queryVoteInitialTimeout, s.height, s.round, tickerTargetQueryVotes)
	} else {
		cpMainVotes := s.log.CPMainVoteVoteSet(s.round)
		switch {
		case cpMainVotes.HasAnyVoteFor(s.cpRound-1, vote.CPValueOne):
			s.logger.Debug("cp: one main-vote for one", "b", "1")

			vote1 := cpMainVotes.GetRandomVote(s.cpRound-1, vote.CPValueOne)
			just1 := &vote.JustPreVoteHard{
				QCert: vote1.CPJust().(*vote.JustMainVoteNoConflict).QCert,
			}
			s.signAddCPPreVote(hash.UndefHash, s.cpRound, vote.CPValueOne, just1)

		case cpMainVotes.HasAnyVoteFor(s.cpRound-1, vote.CPValueZero):
			s.logger.Debug("cp: one main-vote for zero", "b", "0")

			vote0 := cpMainVotes.GetRandomVote(s.cpRound-1, vote.CPValueZero)
			just0 := &vote.JustPreVoteHard{
				QCert: vote0.CPJust().(*vote.JustMainVoteNoConflict).QCert,
			}
			s.signAddCPPreVote(*s.cpWeakValidity, s.cpRound, vote.CPValueZero, just0)

		case cpMainVotes.HasAllVotesFor(s.cpRound-1, vote.CPValueAbstain):
			s.logger.Debug("cp: all main-votes are abstain", "b", "0 (biased)")

			votes := cpMainVotes.BinaryVotes(s.cpRound-1, vote.CPValueAbstain)
			cert := s.makeCertificate(votes)
			just := &vote.JustPreVoteSoft{
				QCert: cert,
			}
			s.signAddCPPreVote(*s.cpWeakValidity, s.cpRound, vote.CPValueZero, just)

		default:
			s.logger.Panic("protocol violated. We have combination of votes for one and zero")
		}
	}

	s.enterNewState(s.cpMainVoteState)
}

func (s *cpPreVoteState) onAddVote(_ *vote.Vote) {
	panic("Unreachable")
}

func (s *cpPreVoteState) name() string {
	return "cp:pre-vote"
}
