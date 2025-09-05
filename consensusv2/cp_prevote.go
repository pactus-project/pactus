package consensusv2

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
)

type cpPreVoteState struct {
	*changeProposer
}

func (s *cpPreVoteState) enter() {
	s.scheduleTimeout(s.config.QueryVoteTimeout, s.height, s.round, tickerTargetQueryVote)

	s.decide()
}

func (s *cpPreVoteState) decide() {
	if s.cpRound == 0 {
		s.decideFirstRound()
	} else {
		s.decideNextRounds()
	}

	s.cpStrongTermination()
	s.absoluteCommit()
}

func (s *cpPreVoteState) decideFirstRound() {
	precommits := s.log.PrecommitVoteSet(s.round)
	precommitQH := precommits.QuorumHash()
	if precommitQH != nil {
		s.cpWeakValidity = precommitQH
		votes := precommits.BlockVotes(*precommitQH)

		cert := s.makeVoteCertificate(votes)
		just := &vote.JustInitNo{
			QCert: cert,
		}
		s.signAddCPPreVote(*precommitQH, s.cpRound, vote.CPValueNo, just)
		s.enterNewState(s.cpMainVoteState)

		return
	}

	if !s.precommitState.hasVoted {
		just := &vote.JustInitYes{}
		s.signAddCPPreVote(hash.UndefHash, s.cpRound, vote.CPValueYes, just)
		s.enterNewState(s.cpMainVoteState)

		return
	}

	cpPreVotes := s.log.CPPreVoteVoteSet(s.round)
	totalVotedPower := cpPreVotes.VotedPower(0) + precommits.VotedPower()

	if totalVotedPower >= 3 {
		just := &vote.JustInitYes{}
		s.signAddCPPreVote(hash.UndefHash, s.cpRound, vote.CPValueYes, just)
		s.enterNewState(s.cpMainVoteState)

		return
	}

	// Waiting for more votes...
	// Transition from Synchronous to Asynchronous Consensus....
}

func (s *cpPreVoteState) decideNextRounds() {
	cpMainVotes := s.log.CPMainVoteVoteSet(s.round)
	switch {
	case cpMainVotes.HasAnyVoteFor(s.cpRound-1, vote.CPValueNo):
		s.logger.Debug("cp: one main-vote for zero", "b", "0")

		vote0 := cpMainVotes.GetRandomVote(s.cpRound-1, vote.CPValueNo)
		just0 := &vote.JustPreVoteHard{
			QCert: vote0.CPJust().(*vote.JustMainVoteNoConflict).QCert,
		}
		s.signAddCPPreVote(*s.cpWeakValidity, s.cpRound, vote.CPValueNo, just0)

	case cpMainVotes.HasAnyVoteFor(s.cpRound-1, vote.CPValueYes):
		s.logger.Debug("cp: one main-vote for one", "b", "1")

		vote1 := cpMainVotes.GetRandomVote(s.cpRound-1, vote.CPValueYes)
		just1 := &vote.JustPreVoteHard{
			QCert: vote1.CPJust().(*vote.JustMainVoteNoConflict).QCert,
		}
		s.signAddCPPreVote(hash.UndefHash, s.cpRound, vote.CPValueYes, just1)

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

	s.enterNewState(s.cpMainVoteState)
}

func (s *cpPreVoteState) onAddVote(_ *vote.Vote) {
	s.decide()
}

func (*cpPreVoteState) name() string {
	return "cp:pre-vote"
}
