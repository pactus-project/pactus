package consensus

import (
	"github.com/zarbchain/zarb-go/vote"
)

type commitState struct {
	*consensus
}

func (s *commitState) enter() {
	s.execute()
}

func (s *commitState) execute() {
	precommits := s.pendingVotes.PrecommitVoteSet(s.round)
	precommitQH := precommits.QuorumHash()

	// For any reason, we don't have proposal
	roundProposal := s.pendingVotes.RoundProposal(s.round)
	if roundProposal == nil {
		s.requestForProposal()

		s.logger.Warn("No proposal, send proposal request.")
		return
	}

	// Proposal is not for quorum block
	// It is impossible, but good to keep this check
	if !roundProposal.IsForBlock(*precommitQH) {
		s.pendingVotes.SetRoundProposal(s.round, nil)
		s.logger.Error("Proposal is invalid.", "proposal", roundProposal)
		return
	}

	certBlock := roundProposal.Block()
	cert := precommits.ToCertificate()
	if cert == nil {
		s.logger.Error("Invalid precommits", "precommitQH", precommitQH)
		return
	}

	if err := s.state.CommitBlock(s.height, certBlock, *cert); err != nil {
		s.logger.Warn("committing block failed", "block", certBlock, "err", err)
		return
	}

	s.logger.Info("Block committed, Schedule new height", "precommitQH", precommitQH)
	// Now we can broadcast the committed block
	s.broadcastBlock(s.height, &certBlock, cert)

	s.enterNewState(s.newHeightState)
}

func (s *commitState) voteAdded(v *vote.Vote) {
	s.execute()
}

func (s *commitState) timedout(t *ticker) {
}

func (s *commitState) name() string {
	return commitName
}
