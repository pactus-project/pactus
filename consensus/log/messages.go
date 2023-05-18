package log

import (
	"github.com/pactus-project/pactus/consensus/voteset"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/logger"
)

type Messages struct {
	prepareVotes        *voteset.VoteSet
	precommitVotes      *voteset.VoteSet
	changeProposerVotes *voteset.VoteSet
	proposal            *proposal.Proposal
}

func (m *Messages) addVote(v *vote.Vote) error {
	vs := m.voteSet(v.Type())
	return vs.AddVote(v)
}

func (m *Messages) HasVote(hash hash.Hash) bool {
	votes := m.AllVotes()
	for _, v := range votes {
		if v.Hash().EqualsTo(hash) {
			return true
		}
	}
	return false
}

func (m *Messages) AllVotes() []*vote.Vote {
	votes := []*vote.Vote{}
	votes = append(votes, m.prepareVotes.AllVotes()...)
	votes = append(votes, m.precommitVotes.AllVotes()...)
	votes = append(votes, m.changeProposerVotes.AllVotes()...)

	return votes
}

func (m *Messages) voteSet(voteType vote.Type) *voteset.VoteSet {
	switch voteType {
	case vote.VoteTypePrepare:
		return m.prepareVotes
	case vote.VoteTypePrecommit:
		return m.precommitVotes
	case vote.VoteTypeChangeProposer:
		return m.changeProposerVotes
	}

	logger.Panic("unexpected vote type", "voteType", voteType)
	return nil
}
