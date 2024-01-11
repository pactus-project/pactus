package voteset

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/errors"
)

type roundVotes struct {
	// Each vote can have one of 3 possible values: {0,1,Abstain}.
	voteBoxes  [3]*voteBox
	allVotes   map[crypto.Address]*vote.Vote
	votedPower int64
}

func newRoundVotes() *roundVotes {
	voteBoxes := [3]*voteBox{}
	voteBoxes[vote.CPValueZero] = newVoteBox()
	voteBoxes[vote.CPValueOne] = newVoteBox()
	voteBoxes[vote.CPValueAbstain] = newVoteBox()

	return &roundVotes{
		voteBoxes:  voteBoxes,
		allVotes:   make(map[crypto.Address]*vote.Vote),
		votedPower: 0,
	}
}

func (rv *roundVotes) addVote(v *vote.Vote, power int64) {
	vb := rv.voteBoxes[v.CPValue()]
	vb.addVote(v, power)
}

type BinaryVoteSet struct {
	*voteSet
	roundVotes []*roundVotes
}

func NewCPPreVoteVoteSet(round int16, totalPower int64,
	validators map[crypto.Address]*validator.Validator,
) *BinaryVoteSet {
	voteSet := newVoteSet(round, totalPower, validators)

	return newBinaryVoteSet(voteSet)
}

func NewCPMainVoteVoteSet(round int16, totalPower int64,
	validators map[crypto.Address]*validator.Validator,
) *BinaryVoteSet {
	voteSet := newVoteSet(round, totalPower, validators)

	return newBinaryVoteSet(voteSet)
}

func NewCPDecidedVoteVoteSet(round int16, totalPower int64,
	validators map[crypto.Address]*validator.Validator,
) *BinaryVoteSet {
	voteSet := newVoteSet(round, totalPower, validators)

	return newBinaryVoteSet(voteSet)
}

func newBinaryVoteSet(voteSet *voteSet) *BinaryVoteSet {
	return &BinaryVoteSet{
		voteSet:    voteSet,
		roundVotes: make([]*roundVotes, 0, 1),
	}
}

func (vs *BinaryVoteSet) mustGetRoundVotes(cpRound int16) *roundVotes {
	for i := len(vs.roundVotes); i <= int(cpRound); i++ {
		rv := newRoundVotes()
		vs.roundVotes = append(vs.roundVotes, rv)
	}

	return vs.roundVotes[cpRound]
}

// AllVotes returns a list of all votes in the VoteSet.
func (vs *BinaryVoteSet) AllVotes() []*vote.Vote {
	votes := make([]*vote.Vote, 0)
	for _, rv := range vs.roundVotes {
		for _, v := range rv.allVotes {
			votes = append(votes, v)
		}
	}

	return votes
}

// AddVote attempts to add a vote to the VoteSet. Returns an error if the vote is invalid.
func (vs *BinaryVoteSet) AddVote(v *vote.Vote) (bool, error) {
	power, err := vs.voteSet.verifyVote(v)
	if err != nil {
		return false, err
	}

	roundVotes := vs.mustGetRoundVotes(v.CPRound())
	existingVote, ok := roundVotes.allVotes[v.Signer()]
	if ok {
		if existingVote.Hash() != v.Hash() {
			err = errors.Error(errors.ErrDuplicateVote)
		} else {
			// The vote is already added
			return false, nil
		}
	} else {
		roundVotes.allVotes[v.Signer()] = v
		roundVotes.votedPower += power
	}

	roundVotes.addVote(v, power)

	return true, err
}

func (vs *BinaryVoteSet) HasOneThirdOfTotalPower(cpRound int16) bool {
	roundVotes := vs.mustGetRoundVotes(cpRound)

	return vs.isOneThirdOfTotalPower(roundVotes.votedPower)
}

func (vs *BinaryVoteSet) HasTwoThirdOfTotalPower(cpRound int16) bool {
	roundVotes := vs.mustGetRoundVotes(cpRound)

	return vs.isTwoThirdOfTotalPower(roundVotes.votedPower)
}

func (vs *BinaryVoteSet) HasAnyVoteFor(cpRound int16, cpValue vote.CPValue) bool {
	roundVotes := vs.mustGetRoundVotes(cpRound)

	return roundVotes.voteBoxes[cpValue].votedPower > 0
}

func (vs *BinaryVoteSet) HasAllVotesFor(cpRound int16, cpValue vote.CPValue) bool {
	roundVotes := vs.mustGetRoundVotes(cpRound)

	return roundVotes.voteBoxes[cpValue].votedPower == roundVotes.votedPower
}

func (vs *BinaryVoteSet) HasQuorumVotesFor(cpRound int16, cpValue vote.CPValue) bool {
	roundVotes := vs.mustGetRoundVotes(cpRound)

	return vs.isTwoThirdOfTotalPower(roundVotes.voteBoxes[cpValue].votedPower)
}

func (vs *BinaryVoteSet) BinaryVotes(cpRound int16, cpValue vote.CPValue) map[crypto.Address]*vote.Vote {
	votes := map[crypto.Address]*vote.Vote{}
	roundVotes := vs.mustGetRoundVotes(cpRound)
	voteBox := roundVotes.voteBoxes[cpValue]
	for a, v := range voteBox.votes {
		votes[a] = v
	}

	return votes
}

func (vs *BinaryVoteSet) GetRandomVote(cpRound int16, cpValue vote.CPValue) *vote.Vote {
	roundVotes := vs.mustGetRoundVotes(cpRound)
	for _, v := range roundVotes.voteBoxes[cpValue].votes {
		return v
	}

	return nil
}
