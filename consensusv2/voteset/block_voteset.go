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

// AddVote attempts to add a vote to the VoteSet.
// Returns an error if the vote is invalid or if it is a duplicate vote.
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
	if vs.hasTwoFPlusOnePower(blockVotes.votedPower) {
		quorumHash := vote.BlockHash()
		vs.quorumHash = &quorumHash
	}

	return true, err
}

// HasVoted checks whether the given address has voted.
func (vs *BlockVoteSet) HasVoted(addr crypto.Address) bool {
	return vs.allVotes[addr] != nil
}

// HasAbsoluteQuorum checks whether there is a block that has received
// an absolute quorum of votes (3f+1 of total voting power).
func (vs *BlockVoteSet) HasAbsoluteQuorum() bool {
	if vs.quorumHash == nil {
		return false
	}
	blockVotes := vs.mustGetBlockVotes(*vs.quorumHash)

	return vs.hasThreeFPlusOnePower(blockVotes.votedPower)
}

// HasMajorityQuorum checks whether the given block has received
// a majority quorum of votes (2f+1 of total voting power).
func (vs *BlockVoteSet) HasMajorityQuorum(blockHash hash.Hash) bool {
	blockVotes := vs.mustGetBlockVotes(blockHash)

	return vs.hasTwoFPlusOnePower(blockVotes.votedPower)
}

// HasQuorumHash reports whether there exists a block that has received
// a quorum of votes (2f+1 of total voting power).
func (vs *BlockVoteSet) HasQuorumHash() bool {
	return vs.quorumHash != nil
}

// QuorumHash returns the hash of the block that has received
// a quorum of votes (2f+1 of total voting power).
// If no block has reached the quorum threshold, it returns nil.
func (vs *BlockVoteSet) QuorumHash() *hash.Hash {
	return vs.quorumHash
}
