package consensus

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type prepareState struct {
	*consensus
	hasVoted bool
}

func (s *prepareState) enter() {
	s.hasVoted = false

	changeProperTimeout := s.config.CalculateChangeProposerTimeout(s.round)
	queryProposalTimeout := changeProperTimeout / 2
	s.scheduleTimeout(queryProposalTimeout, s.height, s.round, tickerTargetQueryProposal)
	s.scheduleTimeout(changeProperTimeout, s.height, s.round, tickerTargetChangeProposer)

	s.decide()
}

func (s *prepareState) decide() {
	s.vote()

	prepares := s.log.PrepareVoteSet(s.round)
	prepareQH := prepares.QuorumHash()
	if prepareQH != nil {
		s.logger.Debug("prepare has quorum", "prepareQH", prepareQH)
		s.enterNewState(s.precommitState)
	} else {
		// Liveness on PBFT
		//
		// If a replica receives a set of f+1 valid change-proposer votes for this round,
		// it sends a change-proposer vote, even if its timer has not expired;
		// this prevents it from starting the change-proposer state too late.
		voteset := s.log.ChangeProposerVoteSet(s.round)
		if voteset.BlockHashHasOneThirdOfTotalPower(hash.UndefHash) {
			s.enterNewState(s.changeProposerState)
		}
	}
}

func (s *prepareState) vote() {
	if s.hasVoted {
		return
	}

	roundProposal := s.log.RoundProposal(s.round)
	if roundProposal == nil {
		s.logger.Warn("no proposal yet.")
		return
	}

	// Everything is good
	s.logger.Info("proposal approved", "proposal", roundProposal)
	s.signAddVote(vote.VoteTypePrepare, roundProposal.Block().Hash())
	s.hasVoted = true
}

func (s *prepareState) onAddVote(v *vote.Vote) {
	s.doAddVote(v)
	if v.Round() == s.round {
		s.decide()
	}
}

func (s *prepareState) onSetProposal(p *proposal.Proposal) {
	s.doSetProposal(p)
	if p.Round() == s.round {
		s.decide()
	}
}

func (s *prepareState) onTimeout(t *ticker) {
	if t.Target == tickerTargetQueryProposal {
		s.queryProposal()
		s.decide()
	} else if t.Target == tickerTargetChangeProposer {
		s.enterNewState(s.changeProposerState)
	} else {
		s.logger.Debug("invalid ticker", "ticker", t)
	}
}

func (s *prepareState) name() string {
	return "prepare"
}
