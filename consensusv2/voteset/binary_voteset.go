package voteset

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
)

type roundVotes struct {
	// Each vote can have one of 3 possible values: {0,1,Abstain}.
	voteBoxes  [3]*voteBox
	allVotes   map[crypto.Address]*vote.Vote
	votedPower int64
}

func newRoundVotes() *roundVotes {
	voteBoxes := [3]*voteBox{}
	voteBoxes[vote.CPValueNo] = newVoteBox()
	voteBoxes[vote.CPValueYes] = newVoteBox()
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

func NewCPDecidedVoteSet(round int16, totalPower int64,
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
func (vs *BinaryVoteSet) AddVote(vote *vote.Vote) (bool, error) {
	power, err := vs.voteSet.verifyVote(vote)
	if err != nil {
		return false, err
	}

	roundVotes := vs.mustGetRoundVotes(vote.CPRound())
	existingVote, ok := roundVotes.allVotes[vote.Signer()]
	if ok {
		if existingVote.Hash() == vote.Hash() {
			// The vote is already added
			return false, nil
		}

		// It is a duplicated vote
		err = ErrDuplicatedVote
	} else {
		roundVotes.allVotes[vote.Signer()] = vote
		roundVotes.votedPower += power
	}

	roundVotes.addVote(vote, power)

	return true, err
}

func (vs *BinaryVoteSet) HasTwoFPlusOneVotes(cpRound int16) bool {
	roundVotes := vs.mustGetRoundVotes(cpRound)

	return vs.hasTwoFPlusOnePower(roundVotes.votedPower)
}

func (vs *BinaryVoteSet) HasAnyVoteFor(cpRound int16, cpValue vote.CPValue) bool {
	roundVotes := vs.mustGetRoundVotes(cpRound)

	return roundVotes.voteBoxes[cpValue].votedPower > 0
}

func (vs *BinaryVoteSet) HasAllVotesFor(cpRound int16, cpValue vote.CPValue) bool {
	roundVotes := vs.mustGetRoundVotes(cpRound)

	return roundVotes.voteBoxes[cpValue].votedPower == roundVotes.votedPower
}

func (vs *BinaryVoteSet) HasFPlusOneVotesFor(cpRound int16, cpValue vote.CPValue) bool {
	roundVotes := vs.mustGetRoundVotes(cpRound)

	return vs.hasFPlusOnePower(roundVotes.voteBoxes[cpValue].votedPower)
}

func (vs *BinaryVoteSet) HasTwoFPlusOneVotesFor(cpRound int16, cpValue vote.CPValue) bool {
	roundVotes := vs.mustGetRoundVotes(cpRound)

	return vs.hasTwoFPlusOnePower(roundVotes.voteBoxes[cpValue].votedPower)
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
