package consensus

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
)

type cpMainVoteState struct {
	*changeProposer
}

func (s *cpMainVoteState) enter() {
	s.decide()
}

func (s *cpMainVoteState) decide() {
	s.checkForWeakValidity()
	s.detectByzantineProposal()

	cpPreVotes := s.log.CPPreVoteVoteSet(s.round)
	if cpPreVotes.HasTwoThirdOfTotalPower(s.cpRound) {
		if cpPreVotes.HasQuorumVotesFor(s.cpRound, vote.CPValueYes) {
			s.logger.Debug("cp: quorum for pre-votes", "v", "1")

			votes := cpPreVotes.BinaryVotes(s.cpRound, vote.CPValueYes)
			cert := s.makeVoteCertificate(votes)
			just := &vote.JustMainVoteNoConflict{
				QCert: cert,
			}
			s.signAddCPMainVote(hash.UndefHash, s.cpRound, vote.CPValueYes, just)
			s.enterNewState(s.cpDecideState)
		} else if cpPreVotes.HasQuorumVotesFor(s.cpRound, vote.CPValueNo) {
			s.logger.Debug("cp: quorum for pre-votes", "v", "0")

			votes := cpPreVotes.BinaryVotes(s.cpRound, vote.CPValueNo)
			cert := s.makeVoteCertificate(votes)
			just := &vote.JustMainVoteNoConflict{
				QCert: cert,
			}
			s.signAddCPMainVote(*s.cpWeakValidity, s.cpRound, vote.CPValueNo, just)
			s.enterNewState(s.cpDecideState)
		} else {
			s.logger.Debug("cp: no-quorum for pre-votes", "v", "abstain")

			vote0 := cpPreVotes.GetRandomVote(s.cpRound, vote.CPValueNo)
			vote1 := cpPreVotes.GetRandomVote(s.cpRound, vote.CPValueYes)

			just := &vote.JustMainVoteConflict{
				JustNo:  vote0.CPJust(),
				JustYes: vote1.CPJust(),
			}

			s.signAddCPMainVote(*s.cpWeakValidity, s.cpRound, vote.CPValueAbstain, just)
			s.enterNewState(s.cpDecideState)
		}
	}
}

func (s *cpMainVoteState) checkForWeakValidity() {
	if s.cpWeakValidity == nil {
		preVotes := s.log.CPPreVoteVoteSet(s.round)
		preVotesZero := preVotes.BinaryVotes(s.cpRound, vote.CPValueNo)

		for _, v := range preVotesZero {
			bh := v.BlockHash()
			s.cpWeakValidity = &bh

			break
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

func (s *cpMainVoteState) onAddVote(vte *vote.Vote) {
	if vte.Type() == vote.VoteTypeCPPreVote {
		s.decide()
	}

	if vte.IsCPVote() {
		s.cpStrongTermination(vte.Round(), vte.CPRound())
	}
}

func (*cpMainVoteState) name() string {
	return "cp:main-vote"
}
