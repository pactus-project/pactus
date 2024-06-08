package voteset

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/errors"
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
	bv, exists := vs.blockVotes[blockHash]
	if !exists {
		bv = newVoteBox()
		vs.blockVotes[blockHash] = bv
	}

	return bv
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
func (vs *BlockVoteSet) AddVote(v *vote.Vote) (bool, error) {
	power, err := vs.voteSet.verifyVote(v)
	if err != nil {
		return false, err
	}

	existingVote, ok := vs.allVotes[v.Signer()]
	if ok {
		if existingVote.Hash() == v.Hash() {
			// The vote is already added
			return false, nil
		}

		// It is a duplicated vote
		err = errors.Error(errors.ErrDuplicateVote)
	} else {
		vs.allVotes[v.Signer()] = v
	}

	blockVotes := vs.mustGetBlockVotes(v.BlockHash())
	blockVotes.addVote(v, power)
	if vs.hasThreeTPlusOnePower(blockVotes.votedPower) {
		quorumHash := v.BlockHash()
		vs.quorumHash = &quorumHash
	}

	return true, err
}

func (vs *BlockVoteSet) HasVoted(addr crypto.Address) bool {
	return vs.allVotes[addr] != nil
}

// HasAbsoluteQuorum checks if there is a block that has received an absolute quorum of votes (4t+1 of total power).
func (vs *BlockVoteSet) HasAbsoluteQuorum(blockHash hash.Hash) bool {
	blockVotes := vs.mustGetBlockVotes(blockHash)

	return vs.hasFourTPlusOnePower(blockVotes.votedPower)
}

// HasQuorumHash checks if there is a block that has received a quorum of votes (3t+1 of total power).
func (vs *BlockVoteSet) HasQuorumHash() bool {
	return vs.quorumHash != nil
}

// QuorumHash returns the hash of the block that has received a quorum of votes (3t+1 of total power).
// If no block has received the quorum threshold, it returns nil.
func (vs *BlockVoteSet) QuorumHash() *hash.Hash {
	return vs.quorumHash
}

// thresholdPower calculates the threshold power based on the total power.
// The formula used is: t = (n - 1) / 5, where n is the total power.
func (vs *BlockVoteSet) thresholdPower() int64 {
	return (vs.totalPower - 1) / 5
}

// hasFourTPlusOnePower checks whether the given power is greater than or equal to 4t+1,
// where t is the threshold power.
func (vs *BlockVoteSet) hasFourTPlusOnePower(power int64) bool {
	return power >= (4*vs.thresholdPower() + 1)
}

// hasThreeTPlusOnePower checks whether the given power is greater than or equal to 3t+1,
// where t is the threshold power.
func (vs *BlockVoteSet) hasThreeTPlusOnePower(power int64) bool {
	return power >= (3*vs.thresholdPower() + 1)
}
