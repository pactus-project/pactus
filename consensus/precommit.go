package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

type precommitState struct {
	*consensus
	hasVoted bool
}

func (s *precommitState) enter() {
	s.hasVoted = false
	s.decide()
}

func (s *precommitState) decide() {
	s.vote()

	precommits := s.log.PrecommitVoteSet(s.round)
	precommitQH := precommits.QuorumHash()
	if precommitQH != nil {
		s.logger.Debug("precommit has quorum", "precommitQH", precommitQH)
		s.enterNewState(s.commitState)
	} else {
		// Liveness on PBFT
		// ...
		voteset := s.log.ChangeProposerVoteSet(s.round)
		if voteset.BlockHashHasOneThirdOfTotalPower(hash.UndefHash) {
			s.enterNewState(s.changeProposerState)
		}
	}
}

func (s *precommitState) vote() {
	if s.hasVoted {
		return
	}

	prepares := s.log.PrepareVoteSet(s.round)
	prepareQH := prepares.QuorumHash()
	roundProposal := s.log.RoundProposal(s.round)
	if roundProposal == nil {
		// There is a consensus about a proposal which we don't have it yet.
		// Ask peers for this proposal
		s.queryProposal()
		s.logger.Debug("no proposal yet")
		return
	}

	if !roundProposal.IsForBlock(*prepareQH) {
		s.log.SetRoundProposal(s.round, nil)
		s.queryProposal()
		s.logger.Error("proposal is invalid", "proposal", roundProposal)
		return
	}

	// Everything is good
	s.logger.Info("proposal approved", "proposal", roundProposal)
	s.signAddVote(vote.VoteTypePrecommit, *prepareQH)
	s.hasVoted = true
}

func (s *precommitState) onAddVote(v *vote.Vote) {
	s.doAddVote(v)

	if v.Round() == s.round {
		s.decide()
	}
}

func (s *precommitState) onSetProposal(p *proposal.Proposal) {
	s.doSetProposal(p)
	if p.Round() == s.round {
		s.decide()
	}
}

func (s *precommitState) onTimedout(t *ticker) {
	if t.Target == tickerTargetChangeProposer {
		s.enterNewState(s.changeProposerState)
	} else {
		s.logger.Trace("invalid ticker", "ticker", t)
	}
}

func (s *precommitState) name() string {
	return "precommit"
}
