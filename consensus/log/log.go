package log

import (
	"github.com/pactus-project/pactus/consensus/voteset"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
)

type Log struct {
	validators    []*validator.Validator
	roundMessages []*Messages
}

func NewLog() *Log {
	return &Log{
		roundMessages: make([]*Messages, 0),
	}
}

func (log *Log) RoundMessages(round int16) *Messages {
	if round < int16(len(log.roundMessages)) {
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

func (log *Log) MustGetRoundMessages(round int16) *Messages {
	for i := int16(len(log.roundMessages)); i <= round; i++ {
		rv := &Messages{
			prepareVotes:        voteset.NewVoteSet(i, vote.VoteTypePrepare, log.validators),
			precommitVotes:      voteset.NewVoteSet(i, vote.VoteTypePrecommit, log.validators),
			changeProposerVotes: voteset.NewVoteSet(i, vote.VoteTypeChangeProposer, log.validators),
		}

		// expending votes slice
		log.roundMessages = append(log.roundMessages, rv)
	}

	return log.RoundMessages(round)
}

func (log *Log) AddVote(v *vote.Vote) error {
	m := log.MustGetRoundMessages(v.Round())
	return m.addVote(v)
}

func (log *Log) PrepareVoteSet(round int16) *voteset.VoteSet {
	m := log.MustGetRoundMessages(round)
	return m.voteSet(vote.VoteTypePrepare)
}

func (log *Log) PrecommitVoteSet(round int16) *voteset.VoteSet {
	m := log.MustGetRoundMessages(round)
	return m.voteSet(vote.VoteTypePrecommit)
}

func (log *Log) ChangeProposerVoteSet(round int16) *voteset.VoteSet {
	m := log.MustGetRoundMessages(round)
	return m.voteSet(vote.VoteTypeChangeProposer)
}

func (log *Log) HasRoundProposal(round int16) bool {
	return log.RoundProposal(round) != nil
}

func (log *Log) RoundProposal(round int16) *proposal.Proposal {
	m := log.RoundMessages(round)
	if m == nil {
		return nil
	}
	return m.proposal
}

func (log *Log) SetRoundProposal(round int16, proposal *proposal.Proposal) {
	m := log.MustGetRoundMessages(round)
	m.proposal = proposal
}

func (log *Log) MoveToNewHeight(validators []*validator.Validator) {
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
