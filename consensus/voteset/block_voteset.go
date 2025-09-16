package voteset

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
)

type BlockVoteSet struct {
	*voteSet
	blockVotes map[hash.Hash]*voteBox
	allVotes   map[crypto.Address]*vote.Vote
	quorumHash *hash.Hash
}

func NewPrepareVoteSet(round int16, totalPower int64,
	validators map[crypto.Address]*validator.Validator,
) *BlockVoteSet {
	voteSet := newVoteSet(round, totalPower, validators)

	return newBlockVoteSet(voteSet)
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
	}
}

func (vs *BlockVoteSet) BlockVotes(blockHash hash.Hash) map[crypto.Address]*vote.Vote {
	votes := map[crypto.Address]*vote.Vote{}
	blockVotes := vs.mustGetBlockVotes(blockHash)
	for a, v := range blockVotes.votes {
		votes[a] = v
	}

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
	votes := make([]*vote.Vote, 0)
	for _, v := range vs.allVotes {
		votes = append(votes, v)
	}

	return votes
}

// AddVote attempts to add a vote to the VoteSet. Returns an error if the vote is invalid.
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

		// It is a duplicated vote
		err = ErrDuplicatedVote
	} else {
		vs.allVotes[vote.Signer()] = vote
	}

	blockVotes := vs.mustGetBlockVotes(vote.BlockHash())
	blockVotes.addVote(vote, power)
	if vs.isTwoThirdOfTotalPower(blockVotes.votedPower) {
		h := vote.BlockHash()
		vs.quorumHash = &h
	}

	return true, err
}

// HasQuorumHash checks if there is a block that has received quorum votes (2/3+ of total power).
func (vs *BlockVoteSet) HasQuorumHash() bool {
	return vs.quorumHash != nil
}

// QuorumHash returns the hash of the block that has received quorum votes (2/3+ of total power).
// If no block has received the quorum threshold (2/3+ of total voting power), it returns nil.
func (vs *BlockVoteSet) QuorumHash() *hash.Hash {
	return vs.quorumHash
}
