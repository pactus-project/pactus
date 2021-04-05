package consensus

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

type prepareState struct {
	*consensus
	hasTimedout bool
}

func (s *prepareState) enter() {
	s.vote()

	sleep := s.config.PrepareTimeout(s.round)
	s.scheduleTimeout(sleep, s.height, s.round, tickerTargetPrepare)
	s.logger.Trace("Prepare scheduled", "timeout", sleep.Seconds())
}

func (s *prepareState) execute() {
	prepares := s.pendingVotes.PrepareVoteSet(s.round)
	if prepares.HasAccumulatedTwoThirdOfTotalPower() {
		s.enterNewState(s.precommitState)
	}
}

func (s *prepareState) vote() {
	roundProposal := s.pendingVotes.RoundProposal(s.round)
	if roundProposal == nil {
		s.logger.Warn("No proposal")
		s.signAddVote(vote.VoteTypePrepare, crypto.UndefHash)
		return
	}

	// Everything is good
	s.logger.Info("Proposal approved", "proposal", roundProposal)
	s.signAddVote(vote.VoteTypePrepare, roundProposal.Block().Hash())
}

func (s *prepareState) voteAdded(v *vote.Vote) {
	if s.hasTimedout {
		s.execute()
	}

	prepares := s.pendingVotes.PrepareVoteSet(s.round)
	prepareQH := prepares.QuorumHash()
	if prepareQH != nil {
		s.logger.Debug("prepare has quorum", "prepareQH", prepareQH)
		s.execute()
	}
}

func (s *prepareState) timedout(t *ticker) {
	if t.Target != tickerTargetPrepare {
		s.logger.Debug("Invalid ticker", "ticker", t)
		return
	}

	s.hasTimedout = true
	s.execute()
}

func (s *prepareState) name() string {
	return prepareName
}
