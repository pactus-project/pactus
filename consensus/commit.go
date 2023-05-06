package consensus

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type commitState struct {
	*consensus
}

func (s *commitState) enter() {
	s.decide()
}

func (s *commitState) decide() {
	precommits := s.log.PrecommitVoteSet(s.round)
	precommitQH := precommits.QuorumHash()

	// For any reason, we don't have proposal
	roundProposal := s.log.RoundProposal(s.round)
	if roundProposal == nil {
		s.logger.Warn("no proposal, send proposal request.")
		s.queryProposal()
		return
	}

	// Proposal is not for quorum block
	// It is impossible, but good to keep this check
	if !roundProposal.IsForBlock(*precommitQH) {
		s.log.SetRoundProposal(s.round, nil)
		s.logger.Error("proposal is invalid", "proposal", roundProposal)
		return
	}

	certBlock := roundProposal.Block()
	cert := precommits.ToCertificate()
	if cert == nil {
		s.logger.Error("invalid precommits", "precommitQH", precommitQH)
		return
	}

	if err := s.state.CommitBlock(s.height, certBlock, cert); err != nil {
		s.logger.Warn("committing block failed", "block", certBlock, "err", err)
		return
	}

	s.logger.Info("block committed, schedule new height", "precommitQH", precommitQH)

	// Now we can announce the committed block and certificate
	s.announceNewBlock(s.height, certBlock, cert)

	s.enterNewState(s.newHeightState)
}

func (s *commitState) onAddVote(v *vote.Vote) {
	s.doAddVote(v)
	s.decide()
}

func (s *commitState) onSetProposal(p *proposal.Proposal) {
	s.doSetProposal(p)
	s.decide()
}

func (s *commitState) onTimeout(_ *ticker) {
	// Ignore timeouts
}

func (s *commitState) name() string {
	return "commit"
}
