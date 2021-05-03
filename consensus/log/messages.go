package log

import (
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/consensus/vote_set"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
)

type Messages struct {
	prepareVotes        *vote_set.VoteSet
	precommitVotes      *vote_set.VoteSet
	changeProposerVotes *vote_set.VoteSet
	proposal            *proposal.Proposal
}

func (m *Messages) addVote(v *vote.Vote) error {
	vs := m.voteSet(v.VoteType())
	return vs.AddVote(v)
}

func (m *Messages) HasVote(hash crypto.Hash) bool {
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

func (m *Messages) voteSet(voteType vote.VoteType) *vote_set.VoteSet {
	switch voteType {
	case vote.VoteTypePrepare:
		return m.prepareVotes
	case vote.VoteTypePrecommit:
		return m.precommitVotes
	case vote.VoteTypeChangeProposer:
		return m.changeProposerVotes
	}

	logger.Panic("Unexpected vote type %d", voteType)
	return nil
}
