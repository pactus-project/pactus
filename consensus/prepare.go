package consensus

import (
	"github.com/zarbchain/zarb-go/vote"
)

type prepareState struct {
	*consensus
	hasVoted bool
}

func (s *prepareState) enter() {
	s.hasVoted = false
	s.execute()
}

func (s *prepareState) execute() {
	s.vote()

	prepares := s.pendingVotes.PrepareVoteSet(s.round)
	prepareQH := prepares.QuorumHash()
	if prepareQH != nil {
		s.logger.Debug("prepare has quorum", "prepareQH", prepareQH)
		s.enterNewState(s.precommitState)
	}
}

func (s *prepareState) vote() {
	if s.hasVoted {
		return
	}

	roundProposal := s.pendingVotes.RoundProposal(s.round)
	if roundProposal == nil {
		s.queryProposal()
		s.logger.Warn("No proposal yet.")
		return
	}

	// Everything is good
	s.logger.Info("Proposal approved", "proposal", roundProposal)
	s.signAddVote(vote.VoteTypePrepare, roundProposal.Block().Hash())
	s.hasVoted = true
}

func (s *prepareState) voteAdded(v *vote.Vote) {
	s.execute()
}

func (s *prepareState) timedout(t *ticker) {
}

func (s *prepareState) name() string {
	return prepareName
}
