package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
)

type prepareState struct {
	*consensus
	hasVoted bool
}

func (s *prepareState) enter() {
	s.hasVoted = false
}

func (s *prepareState) decide() {
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

func (s *prepareState) onAddVote(v *vote.Vote) {
	s.doAddVote(v)
	s.decide()
}

func (s *prepareState) onSetProposal(p *proposal.Proposal) {
	s.doSetProposal(p)
	s.decide()
}

func (s *prepareState) onTimedout(t *ticker) {
	if t.Target != tickerTargetChangeProposer {
		s.logger.Debug("Invalid ticker", "ticker", t)
		return
	}
	s.enterNewState(s.changeProposerState)
}

func (s *prepareState) name() string {
	return "prepare"
}
