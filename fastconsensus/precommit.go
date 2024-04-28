package fastconsensus

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
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
	s.strongCommit()

	precommits := s.log.PrecommitVoteSet(s.round)
	precommitQH := precommits.QuorumHash()
	if precommitQH != nil {
		s.logger.Debug("pre-commit has quorum", "hash", precommitQH)

		roundProposal := s.log.RoundProposal(s.round)
		if roundProposal == nil {
			// There is a consensus about a proposal that we don't have yet.
			// Ask peers for this proposal.
			s.logger.Info("query for a decided proposal", "precommitQH", precommitQH)
			s.queryProposal()

			return
		}

		votes := precommits.BlockVotes(*precommitQH)
		s.blockCert = s.makeBlockCertificate(votes, false)

		s.enterNewState(s.commitState)
	}
}

func (s *precommitState) vote() {
	if s.hasVoted {
		return
	}

	roundProposal := s.log.RoundProposal(s.round)
	if roundProposal == nil {
		s.logger.Debug("no proposal yet")

		return
	}

	// Everything is good
	s.signAddPrecommitVote(roundProposal.Block().Hash())
	s.hasVoted = true
}

func (s *precommitState) onAddVote(_ *vote.Vote) {
	s.decide()
}

func (s *precommitState) onSetProposal(_ *proposal.Proposal) {
	s.decide()
}

func (s *precommitState) onTimeout(_ *ticker) {
	// Ignore timeouts
}

func (s *precommitState) name() string {
	return "precommit"
}
