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

// AddVote attempts to add a vote to the VoteSet. Returns an error if the vote is invalid.
func (vs *BlockVoteSet) AddVote(vte *vote.Vote) (bool, error) {
	power, err := vs.voteSet.verifyVote(vte)
	if err != nil {
		return false, err
	}

	existingVote, ok := vs.allVotes[vte.Signer()]
	if ok {
		if existingVote.Hash() == vte.Hash() {
			// The vote is already added
			return false, nil
		}

		// It is a duplicated vote
		err = ErrDuplicatedVote
	} else {
		vs.allVotes[vte.Signer()] = vte
	}

	blockVotes := vs.mustGetBlockVotes(vte.BlockHash())
	blockVotes.addVote(vte, power)
	if vs.hasFPlusOnePower(blockVotes.votedPower) {
		quorumHash := vte.BlockHash()
		vs.quorumHash = &quorumHash
	}

	return true, err
}

func (vs *BlockVoteSet) HasVoted(addr crypto.Address) bool {
	return vs.allVotes[addr] != nil
}

// HasAbsoluteQuorum checks if there is a block that has received an absolute quorum of votes (3f+1 of total power).
func (vs *BlockVoteSet) HasAbsoluteQuorum() bool {
	if vs.quorumHash == nil {
		return false
	}
	blockVotes := vs.mustGetBlockVotes(*vs.quorumHash)

	return vs.hasThreeFPlusOnePower(blockVotes.votedPower)
}

// HasMajorityQuorum checks if there is a block that has received an majority quorum of votes (2f+1 of total power).
func (vs *BlockVoteSet) HasMajorityQuorum(blockHash hash.Hash) bool {
	blockVotes := vs.mustGetBlockVotes(blockHash)

	return vs.hasTwoFPlusOnePower(blockVotes.votedPower)
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
