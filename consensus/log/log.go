package log

import (
	"github.com/pactus-project/pactus/consensus/voteset"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
)

type Log struct {
	validators    map[crypto.Address]*validator.Validator
	totalPower    int64
	roundMessages map[types.Round]*Messages
}

func NewLog() *Log {
	return &Log{
		roundMessages: make(map[types.Round]*Messages, 0),
	}
}

func (log *Log) RoundMessages(round types.Round) *Messages {
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

func (log *Log) mustGetRoundMessages(round types.Round) *Messages {
	msgs, ok := log.roundMessages[round]
	if !ok {
		msgs = &Messages{
			prepareVotes:   voteset.NewPrepareVoteSet(round, log.totalPower, log.validators),
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

func (log *Log) PrepareVoteSet(round types.Round) *voteset.BlockVoteSet {
	msgs := log.mustGetRoundMessages(round)

	return msgs.prepareVotes
}

func (log *Log) PrecommitVoteSet(round types.Round) *voteset.BlockVoteSet {
	msgs := log.mustGetRoundMessages(round)

	return msgs.precommitVotes
}

func (log *Log) CPPreVoteVoteSet(round types.Round) *voteset.BinaryVoteSet {
	msgs := log.mustGetRoundMessages(round)

	return msgs.cpPreVotes
}

func (log *Log) CPMainVoteVoteSet(round types.Round) *voteset.BinaryVoteSet {
	msgs := log.mustGetRoundMessages(round)

	return msgs.cpMainVotes
}

func (log *Log) CPDecidedVoteSet(round types.Round) *voteset.BinaryVoteSet {
	msgs := log.mustGetRoundMessages(round)

	return msgs.cpDecidedVotes
}

func (log *Log) HasRoundProposal(round types.Round) bool {
	return log.RoundProposal(round) != nil
}

func (log *Log) RoundProposal(round types.Round) *proposal.Proposal {
	m := log.RoundMessages(round)
	if m == nil {
		return nil
	}

	return m.proposal
}

func (log *Log) SetRoundProposal(round types.Round, prop *proposal.Proposal) {
	msgs := log.mustGetRoundMessages(round)
	msgs.proposal = prop
}

func (log *Log) MoveToNewHeight(validators []*validator.Validator) {
	log.roundMessages = make(map[types.Round]*Messages)
	log.validators = make(map[crypto.Address]*validator.Validator)
	log.totalPower = 0
	for _, val := range validators {
		log.totalPower += val.Power()
		log.validators[val.Address()] = val
	}
}

// VerifyVote checks that the signer of the vote is in the committee and the signature is valid.
func (log *Log) VerifyVote(vte *vote.Vote) error {
	val, ok := log.validators[vte.Signer()]
	if !ok {
		return voteset.IneligibleVoterError{Address: vte.Signer()}
	}

	return vte.Verify(val.PublicKey())
}

// FindVote returns the vote of the signer that takes the same place as the given vote:
// the same type, round and change-proposer round. It returns nil if there is no such a vote.
func (log *Log) FindVote(vte *vote.Vote) *vote.Vote {
	msgs, ok := log.roundMessages[vte.Round()]
	if !ok {
		return nil
	}

	switch vte.Type() {
	case vote.VoteTypePrepare:
		return msgs.prepareVotes.GetVote(vte.Signer())
	case vote.VoteTypePrecommit:
		return msgs.precommitVotes.GetVote(vte.Signer())
	case vote.VoteTypeCPPreVote:
		return msgs.cpPreVotes.GetVote(vte.CPRound(), vte.Signer())
	case vote.VoteTypeCPMainVote:
		return msgs.cpMainVotes.GetVote(vte.CPRound(), vte.Signer())
	case vote.VoteTypeCPDecided:
		return msgs.cpDecidedVotes.GetVote(vte.CPRound(), vte.Signer())
	}

	return nil
}

func (log *Log) CanVote(addr crypto.Address) bool {
	_, ok := log.validators[addr]

	return ok
}
