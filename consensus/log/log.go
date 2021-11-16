package log

import (
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/consensus/voteset"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/validator"
)

type Log struct {
	height        int
	validators    []*validator.Validator
	roundMessages []*Messages
}

func NewLog() *Log {
	return &Log{
		roundMessages: make([]*Messages, 0),
	}
}

func (log *Log) RoundMessages(round int) *Messages {
	if round < len(log.roundMessages) {
		return log.roundMessages[round]
	}
	return nil
}

func (log *Log) HasVote(h hash.Hash) bool {
	for _, m := range log.roundMessages {
		if m.HasVote(h) {
			return true
		}
	}
	return false
}

func (log *Log) MustGetRoundMessages(round int) *Messages {
	for i := len(log.roundMessages); i <= round; i++ {
		rv := &Messages{
			prepareVotes:        voteset.NewVoteSet(log.height, i, vote.VoteTypePrepare, log.validators),
			precommitVotes:      voteset.NewVoteSet(log.height, i, vote.VoteTypePrecommit, log.validators),
			changeProposerVotes: voteset.NewVoteSet(log.height, i, vote.VoteTypeChangeProposer, log.validators),
		}

		// extendind votes slice
		log.roundMessages = append(log.roundMessages, rv)
	}

	return log.RoundMessages(round)
}

func (log *Log) AddVote(v *vote.Vote) error {
	m := log.MustGetRoundMessages(v.Round())
	return m.addVote(v)
}

func (log *Log) PrepareVoteSet(round int) *voteset.VoteSet {
	m := log.MustGetRoundMessages(round)
	return m.voteSet(vote.VoteTypePrepare)
}

func (log *Log) PrecommitVoteSet(round int) *voteset.VoteSet {
	m := log.MustGetRoundMessages(round)
	return m.voteSet(vote.VoteTypePrecommit)
}

func (log *Log) ChangeProposerVoteSet(round int) *voteset.VoteSet {
	m := log.MustGetRoundMessages(round)
	return m.voteSet(vote.VoteTypeChangeProposer)
}

func (log *Log) HasRoundProposal(round int) bool {
	return log.RoundProposal(round) != nil
}

func (log *Log) RoundProposal(round int) *proposal.Proposal {
	m := log.RoundMessages(round)
	if m == nil {
		return nil
	}
	return m.proposal
}

func (log *Log) SetRoundProposal(round int, proposal *proposal.Proposal) {
	m := log.MustGetRoundMessages(round)
	m.proposal = proposal
}

func (log *Log) MoveToNewHeight(height int, validators []*validator.Validator) {
	log.height = height
	log.roundMessages = make([]*Messages, 0)
	log.validators = validators
}

func (log *Log) CanVote(addr crypto.Address) bool {
	for _, val := range log.validators {
		if val.Address().EqualsTo(addr) {
			return true
		}
	}
	return false
}
