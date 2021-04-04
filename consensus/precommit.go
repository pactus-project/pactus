package consensus

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

type precommitState struct {
	*consensus
	hasTimedout bool
}

func (s *precommitState) enter() {
	sleep := s.config.PrecommitTimeout(s.round)
	s.scheduleTimeout(sleep, s.height, s.round, tickerTargetPrecommit)
	s.logger.Trace("Precommit scheduled", "timeout", sleep.Seconds())

	s.vote()
}

func (s *precommitState) execute() {
	s.vote()

	precommits := s.pendingVotes.PrecommitVoteSet(s.round)
	precommitQH := precommits.QuorumHash()

	if precommitQH != nil {
		if precommitQH.IsUndef() {
			s.enterNewState(s.newRoundState)
		} else {
			s.enterNewState(s.commitState)
		}
	}
}

func (s *precommitState) vote() {
	prepares := s.pendingVotes.PrepareVoteSet(s.round)
	prepareQH := prepares.QuorumHash()
	roundProposal := s.pendingVotes.RoundProposal(s.round)
	if roundProposal == nil && prepareQH != nil && !prepareQH.IsUndef() {
		// There is a consensus about a proposal which we don't have it yet.
		// Ask peers for this proposal
		s.requestForProposal()
		s.logger.Debug("No proposal, send proposal request.")
		return
	}

	if prepareQH == nil {
		s.logger.Info("No quorum for prepare")
		s.signAddVote(vote.VoteTypePrecommit, crypto.UndefHash)
		return
	}

	if prepareQH.IsUndef() {
		s.logger.Info("Undef quorum for prepare")
		s.signAddVote(vote.VoteTypePrecommit, crypto.UndefHash)
		return
	}

	// TODO: write test for me!
	if !roundProposal.IsForBlock(*prepareQH) {
		s.pendingVotes.SetRoundProposal(s.round, nil)
		s.logger.Warn("Invalid proposal.")
		s.signAddVote(vote.VoteTypePrecommit, crypto.UndefHash)
		return
	}

	// Everything is good
	s.logger.Info("Proposal approved", "proposal", roundProposal)
	s.signAddVote(vote.VoteTypePrepare, *prepareQH)
	s.signAddVote(vote.VoteTypePrecommit, *prepareQH)
}

func (s *precommitState) voteAdded(v *vote.Vote) {
	if s.hasTimedout {
		s.execute()
	}

	precommits := s.pendingVotes.PrecommitVoteSet(s.round)
	precommitQH := precommits.QuorumHash()
	if precommitQH != nil {
		s.logger.Debug("precommit has quorum", "precommitQH", precommitQH)
		s.execute()
	}
}

func (s *precommitState) timedout(t *ticker) {
	if t.Target != tickerTargetPrecommit {
		s.logger.Debug("Invalid ticker", "ticker", t)
		return
	}

	s.hasTimedout = true
	s.execute()
}

func (s *precommitState) name() string {
	return precommitName
}
