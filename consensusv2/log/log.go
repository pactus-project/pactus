package log

import (
	"github.com/pactus-project/pactus/consensusv2/voteset"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
)

type Log struct {
	validators    map[crypto.Address]*validator.Validator
	totalPower    int64
	roundMessages map[int16]*Messages
}

func NewLog() *Log {
	return &Log{
		roundMessages: make(map[int16]*Messages, 0),
	}
}

func (log *Log) RoundMessages(round int16) *Messages {
	return log.mustGetRoundMessages(round)
}

func (log *Log) HasVote(h hash.Hash) bool {
	for _, m := range log.roundMessages {
		if m.HasVote(h) {
			return true
		}
	}

	return false
}

func (log *Log) mustGetRoundMessages(round int16) *Messages {
	msgs, ok := log.roundMessages[round]
	if !ok {
		msgs = &Messages{
			precommitVotes: voteset.NewPrecommitVoteSet(round, log.totalPower, log.validators),
			cpPreVotes:     voteset.NewCPPreVoteVoteSet(round, log.totalPower, log.validators),
			cpMainVotes:    voteset.NewCPMainVoteVoteSet(round, log.totalPower, log.validators),
			cpDecidedVotes: voteset.NewCPDecidedVoteSet(round, log.totalPower, log.validators),
		}
		log.roundMessages[round] = msgs
	}

	return msgs
}

func (log *Log) AddVote(v *vote.Vote) (bool, error) {
	msgs := log.mustGetRoundMessages(v.Round())

	return msgs.addVote(v)
}

func (log *Log) PrecommitVoteSet(round int16) *voteset.BlockVoteSet {
	msgs := log.mustGetRoundMessages(round)

	return msgs.precommitVotes
}

func (log *Log) CPPreVoteVoteSet(round int16) *voteset.BinaryVoteSet {
	msgs := log.mustGetRoundMessages(round)

	return msgs.cpPreVotes
}

func (log *Log) CPMainVoteVoteSet(round int16) *voteset.BinaryVoteSet {
	msgs := log.mustGetRoundMessages(round)

	return msgs.cpMainVotes
}

func (log *Log) CPDecidedVoteSet(round int16) *voteset.BinaryVoteSet {
	msgs := log.mustGetRoundMessages(round)

	return msgs.cpDecidedVotes
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

func (log *Log) SetRoundProposal(round int16, prop *proposal.Proposal) {
	msgs := log.mustGetRoundMessages(round)
	msgs.proposal = prop
}

func (log *Log) MoveToNewHeight(validators []*validator.Validator) {
	log.roundMessages = make(map[int16]*Messages)
	log.validators = make(map[crypto.Address]*validator.Validator)
	log.totalPower = 0
	for _, val := range validators {
		log.totalPower += val.Power()
		log.validators[val.Address()] = val
	}
}

func (log *Log) CanVote(addr crypto.Address) bool {
	for _, val := range log.validators {
		if val.Address() == addr {
			return true
		}
	}

	return false
}
