package consensus

import (
	"github.com/zarbchain/zarb-go/vote"
)

type precommitState struct {
	*consensus
	hasVoted bool
}

func (s *precommitState) enter() {
	s.hasVoted = false
	s.vote()
}

func (s *precommitState) execute() {
	s.vote()

	precommits := s.pendingVotes.PrecommitVoteSet(s.round)
	precommitQH := precommits.QuorumHash()
	if precommitQH != nil {
		s.logger.Debug("precommit has quorum", "precommitQH", precommitQH)
		s.enterNewState(s.commitState)
	}
}

func (s *precommitState) vote() {
	prepares := s.pendingVotes.PrepareVoteSet(s.round)
	prepareQH := prepares.QuorumHash()
	roundProposal := s.pendingVotes.RoundProposal(s.round)
	if roundProposal == nil {
		// There is a consensus about a proposal which we don't have it yet.
		// Ask peers for this proposal
		s.queryProposal()
		s.logger.Debug("No proposal yet.")
		return
	}

	if !roundProposal.IsForBlock(*prepareQH) {
		s.pendingVotes.SetRoundProposal(s.round, nil)
		s.queryProposal()
		s.logger.Error("Proposal is invalid.", "proposal", roundProposal)
		return
	}

	// Everything is good
	s.logger.Info("Proposal approved", "proposal", roundProposal)
	s.signAddVote(vote.VoteTypePrepare, *prepareQH)
	s.signAddVote(vote.VoteTypePrecommit, *prepareQH)
	s.hasVoted = true
}

func (s *precommitState) voteAdded(v *vote.Vote) {
	s.execute()
}

func (s *precommitState) timedout(t *ticker) {
}

func (s *precommitState) name() string {
	return precommitName
}
