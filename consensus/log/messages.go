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

func (m *Messages) addVote(v *vote.Vote) (bool, error) {
	switch v.Type() {
	case vote.VoteTypePrepare:
		return m.prepareVotes.AddVote(v)
	case vote.VoteTypePrecommit:
		return m.precommitVotes.AddVote(v)
	case vote.VoteTypeCPPreVote:
		return m.cpPreVotes.AddVote(v)
	case vote.VoteTypeCPMainVote:
		return m.cpMainVotes.AddVote(v)
	case vote.VoteTypeCPDecided:
		return m.cpDecidedVotes.AddVote(v)
	}

	return false, fmt.Errorf("unexpected vote type: %v", v.Type())
}

func (m *Messages) HasVote(hash hash.Hash) bool {
	votes := m.AllVotes()
	for _, v := range votes {
		if v.Hash() == hash {
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
