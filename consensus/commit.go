package consensus

type commitState struct {
	*consensus
}

func (s *commitState) enter() {
	s.decide()
}

func (s *commitState) decide() {
	// For any reason, we don't have proposal
	roundProposal := s.log.RoundProposal(s.round)
	if roundProposal == nil {
		s.logger.Warn("no proposal, send proposal request.")
		s.queryProposal()
		return
	}

	certBlock := roundProposal.Block()
	precommits := s.log.PrecommitVoteSet(s.round)
	votes := precommits.BlockVotes(certBlock.Hash())
	cert := s.makeCertificate(votes)
	if err := s.state.CommitBlock(s.height, certBlock, cert); err != nil {
		s.logger.Error("committing block failed", "block", certBlock, "err", err)
		return
	}

	s.logger.Info("block committed, schedule new height", "hash", certBlock.Hash().ShortString())

	// Now we can announce the committed block and certificate
	s.announceNewBlock(s.height, certBlock, cert)

	s.enterNewState(s.newHeightState)
}

func (s *commitState) timeout(_ *ticker) {
	// Ignore timeouts
}

func (s *commitState) name() string {
	return "commit"
}
