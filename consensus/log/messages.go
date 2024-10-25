package log

import (
	"fmt"

	"github.com/pactus-project/pactus/consensus/voteset"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type Messages struct {
	prepareVotes   *voteset.BlockVoteSet  // Prepare votes
	precommitVotes *voteset.BlockVoteSet  // Precommit votes
	cpPreVotes     *voteset.BinaryVoteSet // Change proposer Pre-votes
	cpMainVotes    *voteset.BinaryVoteSet // Change proposer Main-votes
	cpDecidedVotes *voteset.BinaryVoteSet // Change proposer Decided-votes
	proposal       *proposal.Proposal
}

func (m *Messages) addVote(vte *vote.Vote) (bool, error) {
	switch vte.Type() {
	case vote.VoteTypePrepare:
		return m.prepareVotes.AddVote(vte)
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
	votes := []*vote.Vote{}
	votes = append(votes, m.prepareVotes.AllVotes()...)
	votes = append(votes, m.precommitVotes.AllVotes()...)
	votes = append(votes, m.cpPreVotes.AllVotes()...)
	votes = append(votes, m.cpMainVotes.AllVotes()...)
	votes = append(votes, m.cpDecidedVotes.AllVotes()...)

	return votes
}
