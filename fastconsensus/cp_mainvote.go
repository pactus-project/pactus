package fastconsensus

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type cpMainVoteState struct {
	*changeProposer
}

func (s *cpMainVoteState) enter() {
	s.decide()
}

func (s *cpMainVoteState) decide() {
	s.strongCommit()
	s.cpStrongTermination()
	s.checkForWeakValidity()
	s.detectByzantineProposal()

	cpPreVotes := s.log.CPPreVoteVoteSet(s.round)
	if cpPreVotes.HasTwoFPlusOneVotes(s.cpRound) {
		if cpPreVotes.HasTwoFPlusOneVotesFor(s.cpRound, vote.CPValueYes) {
			s.logger.Debug("cp: quorum for pre-votes", "value", "yes")

			votes := cpPreVotes.BinaryVotes(s.cpRound, vote.CPValueYes)
			cert := s.makeVoteCertificate(votes)
			just := &vote.JustMainVoteNoConflict{
				QCert: cert,
			}
			s.signAddCPMainVote(hash.UndefHash, s.cpRound, vote.CPValueYes, just)
			s.enterNewState(s.cpDecideState)
		} else if cpPreVotes.HasTwoFPlusOneVotesFor(s.cpRound, vote.CPValueNo) {
			s.logger.Debug("cp: quorum for pre-votes", "value", "no")

			votes := cpPreVotes.BinaryVotes(s.cpRound, vote.CPValueNo)
			cert := s.makeVoteCertificate(votes)
			just := &vote.JustMainVoteNoConflict{
				QCert: cert,
			}
			s.signAddCPMainVote(*s.cpWeakValidity, s.cpRound, vote.CPValueNo, just)
			s.enterNewState(s.cpDecideState)
		} else {
			s.logger.Debug("cp: no-quorum for pre-votes", "value", "abstain")

			vote0 := cpPreVotes.GetRandomVote(s.cpRound, vote.CPValueNo)
			vote1 := cpPreVotes.GetRandomVote(s.cpRound, vote.CPValueYes)

			just := &vote.JustMainVoteConflict{
				Just0: vote0.CPJust(),
				Just1: vote1.CPJust(),
			}

			s.signAddCPMainVote(*s.cpWeakValidity, s.cpRound, vote.CPValueAbstain, just)
			s.enterNewState(s.cpDecideState)
		}
	}
}

func (s *cpMainVoteState) checkForWeakValidity() {
	if s.cpWeakValidity == nil {
		preVotes := s.log.CPPreVoteVoteSet(s.round)
		randVote := preVotes.GetRandomVote(s.cpRound, vote.CPValueNo)
		if randVote != nil {
			bh := randVote.BlockHash()
			s.cpWeakValidity = &bh
		}
	}
}

func (s *cpMainVoteState) detectByzantineProposal() {
	if s.cpWeakValidity != nil {
		roundProposal := s.log.RoundProposal(s.round)

		if roundProposal != nil &&
			roundProposal.Block().Hash() != *s.cpWeakValidity {
			s.logger.Warn("double proposal detected",
				"prepared", s.cpWeakValidity,
				"roundProposal", roundProposal.Block().Hash())

			s.log.SetRoundProposal(s.round, nil)
			s.queryProposal()
		}
	}
}

func (s *cpMainVoteState) onAddVote(_ *vote.Vote) {
	s.decide()
}

func (s *cpMainVoteState) onSetProposal(_ *proposal.Proposal) {
	// Ignore proposal
}

func (s *cpMainVoteState) onTimeout(_ *ticker) {
	// Ignore timeouts
}

func (s *cpMainVoteState) name() string {
	return "cp:main-vote"
}
