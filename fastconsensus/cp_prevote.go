package fastconsensus

import (
	"time"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
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
	s.strongCommit()
	s.cpStrongTermination()

	if s.cpRound == 0 {
		// broadcast the initial value
		prepares := s.log.PrepareVoteSet(s.round)
		prepareQH := prepares.QuorumHash()
		if prepareQH != nil {
			s.cpWeakValidity = prepareQH
			votes := prepares.BlockVotes(*prepareQH)
			cert := s.makeVoteCertificate(votes)
			just := &vote.JustInitNo{
				QCert: cert,
			}
			s.signAddCPPreVote(*s.cpWeakValidity, s.cpRound, 0, just)
		} else {
			if prepares.HasVoted(s.valKey.Address()) {
				preVotes := s.log.CPPreVoteVoteSet(s.round)
				if !preVotes.HasFPlusOneVotesFor(s.cpRound, vote.CPValueYes) {
					s.logger.Debug("we have proposal but not minority of pre-votes for 'Yes'")

					return
				}
			}
			just := &vote.JustInitYes{}
			s.signAddCPPreVote(hash.UndefHash, s.cpRound, 1, just)
		}
		s.scheduleTimeout(queryVoteInitialTimeout, s.height, s.round, tickerTargetQueryVotes)
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
			cert := s.makeVoteCertificate(votes)
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

func (s *cpPreVoteState) onAddVote(_ *vote.Vote) {
	s.decide()
}

func (s *cpPreVoteState) onSetProposal(_ *proposal.Proposal) {
}

func (s *cpPreVoteState) onTimeout(_ *ticker) {
}

func (s *cpPreVoteState) name() string {
	return "cp:pre-vote"
}
