package log

import (
	"fmt"

	"github.com/pactus-project/pactus/consensusv2/voteset"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type Messages struct {
	precommitVotes *voteset.BlockVoteSet  // Precommit votes
	cpPreVotes     *voteset.BinaryVoteSet // Change proposer Pre-votes
	cpMainVotes    *voteset.BinaryVoteSet // Change proposer Main-votes
	cpDecidedVotes *voteset.BinaryVoteSet // Change proposer Decided-votes
	proposal       *proposal.Proposal
}

func (m *Messages) addVote(vte *vote.Vote) (bool, error) {
	switch vte.Type() {
	case vote.VoteTypePrepare:
		// Deprecated
	case vote.VoteTypePrecommit:
		return m.precommitVotes.AddVote(vte)
	case vote.VoteTypeCPPreVote:
		return m.cpPreVotes.AddVote(vte)
	case vote.VoteTypeCPMainVote:
		return m.cpMainVotes.AddVote(vte)
	case vote.VoteTypeCPDecided:
		return m.cpDecidedVotes.AddVote(vte)
	}

	return false, fmt.Errorf("unexpected vote type: %v", vte.Type())
}

func (m *Messages) HasVote(h hash.Hash) bool {
	votes := m.AllVotes()
	for _, v := range votes {
		if v.Hash() == h {
			return true
		}
	}

	return false
}

func (m *Messages) AllVotes() []*vote.Vote {
	precommit := m.precommitVotes.AllVotes()
	cpPre := m.cpPreVotes.AllVotes()
	cpMain := m.cpMainVotes.AllVotes()
	cpDecided := m.cpDecidedVotes.AllVotes()

	votes := make([]*vote.Vote, 0, len(precommit)+len(cpPre)+len(cpMain)+len(cpDecided))
	votes = append(votes, precommit...)
	votes = append(votes, cpPre...)
	votes = append(votes, cpMain...)
	votes = append(votes, cpDecided...)

	return votes
}
