package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto"
)

type prepareState struct {
	*consensus
	hasVoted bool
}

func (s *prepareState) enter() {
	s.hasVoted = false
}

func (s *prepareState) decide() {
	// Liveness on PBFT
	//
	// If a replica receives a set of f+1 valid change-proposer votes for the next round
	// it sends a change-proposer vote for this round, even if its timer has not expired;
	// this prevents it from starting the next change-proposer state too late.
	voteset := s.pendingVotes.ChangeProposerVoteSet(s.round)
	if voteset.BlockHashHasOneThirdOfTotalPower(crypto.UndefHash) {
		s.enterNewState(s.changeProposerState)
		return
	}

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
	if v.Round() == s.round &&
		v.VoteType() == vote.VoteTypePrepare {
		s.decide()
	}
}

func (s *prepareState) onSetProposal(p *proposal.Proposal) {
	s.doSetProposal(p)
	if p.Round() == s.round {
		s.decide()
	}
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
