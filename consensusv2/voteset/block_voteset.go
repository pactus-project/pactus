package voteset

import (
	"maps"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
)

type BlockVoteSet struct {
	*voteSet
	blockVotes map[hash.Hash]*voteBox
	allVotes   map[crypto.Address]*vote.Vote
	votedPower int64
}

func NewPrecommitVoteSet(round int16, totalPower int64,
	validators map[crypto.Address]*validator.Validator,
) *BlockVoteSet {
	voteSet := newVoteSet(round, totalPower, validators)

	return newBlockVoteSet(voteSet)
}

func newBlockVoteSet(voteSet *voteSet) *BlockVoteSet {
	return &BlockVoteSet{
		voteSet:    voteSet,
		blockVotes: make(map[hash.Hash]*voteBox),
		allVotes:   make(map[crypto.Address]*vote.Vote),
		votedPower: 0,
	}
}

func (vs *BlockVoteSet) BlockVotes(blockHash hash.Hash) map[crypto.Address]*vote.Vote {
	votes := map[crypto.Address]*vote.Vote{}
	blockVotes := vs.mustGetBlockVotes(blockHash)
	maps.Copy(votes, blockVotes.votes)

	return votes
}

func (vs *BlockVoteSet) mustGetBlockVotes(blockHash hash.Hash) *voteBox {
	blockVotes, exists := vs.blockVotes[blockHash]
	if !exists {
		blockVotes = newVoteBox()
		vs.blockVotes[blockHash] = blockVotes
	}

	return blockVotes
}

// AllVotes returns a list of all votes in the VoteSet.
func (vs *BlockVoteSet) AllVotes() []*vote.Vote {
	votes := make([]*vote.Vote, 0, len(vs.allVotes))
	for _, v := range vs.allVotes {
		votes = append(votes, v)
	}

	return votes
}

// AddVote attempts to add a vote to the VoteSet.
// Returns an error if the vote is invalid or if it is a double vote.
func (vs *BlockVoteSet) AddVote(vote *vote.Vote) (bool, error) {
	power, err := vs.voteSet.verifyVote(vote)
	if err != nil {
		return false, err
	}

	existingVote, ok := vs.allVotes[vote.Signer()]
	if ok {
		if existingVote.Hash() == vote.Hash() {
			// The vote is already added
			return false, nil
		}

		// It is a double vote
		err = ErrDoubleVote
	} else {
		vs.allVotes[vote.Signer()] = vote
		vs.votedPower += power
	}

	blockVotes := vs.mustGetBlockVotes(vote.BlockHash())
	blockVotes.addVote(vote, power)

	return true, err
}

// Has2FP1Votes checks whether has received 2f+1 votes.
func (vs *BlockVoteSet) Has2FP1Votes() bool {
	return vs.has2FP1Power(vs.votedPower)
}

// Has3FP1VotesFor checks whether the given block has received 3f+1 votes.
func (vs *BlockVoteSet) Has3FP1VotesFor(blockHash hash.Hash) bool {
	blockVotes := vs.mustGetBlockVotes(blockHash)

	return vs.has3FP1Power(blockVotes.votedPower)
}

// Has2FP1VotesFor checks whether the given block has received 2f+1 votes.
func (vs *BlockVoteSet) Has2FP1VotesFor(blockHash hash.Hash) bool {
	blockVotes := vs.mustGetBlockVotes(blockHash)

	return vs.has2FP1Power(blockVotes.votedPower)
}

// Has1FP1VotesFor checks whether the given block has received f+1 votes.
func (vs *BlockVoteSet) Has1FP1VotesFor(blockHash hash.Hash) bool {
	blockVotes := vs.mustGetBlockVotes(blockHash)

	return vs.has1FP1Power(blockVotes.votedPower)
}

// VotedPower returns the total voting power of the votes.
func (vs *BlockVoteSet) VotedPower() int64 {
	return vs.votedPower
}
