package vote

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/validator"
)

type blockVotes struct {
	votes map[crypto.Address]*Vote //
	sum   int                      // vote sum
}

func newBlockVotes() *blockVotes {
	return &blockVotes{
		votes: make(map[crypto.Address]*Vote),
		sum:   0,
	}
}

func (vs *blockVotes) addVote(vote *Vote) bool {
	signer := vote.Signer()
	if existing, ok := vs.votes[signer]; ok {
		if !existing.data.Signature.EqualsTo(*vote.data.Signature) {
			// Signature malleability?
			logger.Panic("Invalid vote")
		} else {
			//
			return false
		}
	}

	vs.votes[signer] = vote
	vs.sum++

	return true
}

type VoteSet struct {
	height       int
	round        int
	voteType     VoteType
	valSet       validator.ValidatorSetReader
	votesByBlock map[crypto.Hash]*blockVotes
	sum          int
	quorum       *crypto.Hash
}

func NewVoteSet(height int, round int, voteType VoteType, valSet validator.ValidatorSetReader) *VoteSet {
	return &VoteSet{
		height:       height,
		round:        round,
		voteType:     voteType,
		valSet:       valSet,
		votesByBlock: make(map[crypto.Hash]*blockVotes),
	}
}

func (vs *VoteSet) Type() VoteType { return vs.voteType }
func (vs *VoteSet) Height() int    { return vs.height }
func (vs *VoteSet) Round() int     { return vs.round }
func (vs *VoteSet) Len() int       { return vs.sum }

func (vs *VoteSet) AllVotes() []*Vote {
	votes := make([]*Vote, 0)

	for _, blockVotes := range vs.votesByBlock {
		for _, vote := range blockVotes.votes {
			votes = append(votes, vote)
		}
	}
	return votes
}

func (vs *VoteSet) AddVote(vote *Vote) (bool, error) {
	signer := vote.Signer()
	blockHash := vote.BlockHash()

	if signer.SanityCheck() != nil {
		return false, errors.Errorf(errors.ErrInvalidVote, "Empty address")
	}

	if (vote.data.Height != vs.height) ||
		(vote.data.Round != vs.round) ||
		(vote.data.VoteType != vs.voteType) {
		return false, errors.Errorf(errors.ErrInvalidVote, "Expected %d/%d/%s, but got %d/%d/%s",
			vs.height, vs.round, vs.voteType,
			vote.Height(), vote.Round(), vote.VoteType())
	}

	val := vs.valSet.Validator(signer)
	if val == nil {
		return false, errors.Errorf(errors.ErrInvalidVote, "Cannot find validator %s in valSet", signer)
	}

	if err := vote.Verify(val.PublicKey()); err != nil {
		return false, errors.Errorf(errors.ErrInvalidVote, "Failed to verify vote")
	}

	blockVotes, exists := vs.votesByBlock[blockHash]
	if !exists {
		blockVotes = newBlockVotes()
		vs.votesByBlock[blockHash] = blockVotes
	}

	// check for conflict
	for id, v := range vs.votesByBlock {
		if id != vote.data.BlockHash {
			duplicated, ok := v.votes[signer]

			if ok {
				// A possible scenario:
				// A peer doesn't have a proposal, he votes for undef.
				// Later he receives the proposal, so he vote again.
				// We should ignore undef vote
				if duplicated.data.BlockHash.IsUndef() {
					// Remove undef vote and replace it with new vote
					v.sum--
					vs.sum--
					delete(v.votes, signer)
				} else if vote.data.BlockHash.IsUndef() {
					// Because of network latency, we might receive undef vote after block vote.
					// Ignore undef vote in this case.
					return false, nil
				} else if duplicated.data.BlockHash != blockHash {
					// Duplicated vote:
					// 1- Same signer
					// 2- Previous blockhash is not undef
					// 3- Block hashes are different
					return false, errors.Error(errors.ErrDuplicateVote)
				}
			}
		}
	}

	added := blockVotes.addVote(vote)
	if added {
		if vs.hasQuorum(blockVotes.sum) {
			vs.quorum = &blockHash
		}
		vs.sum++
	}

	return added, nil
}
func (vs *VoteSet) hasQuorum(sum int) bool {
	return sum > (vs.valSet.Power() * 2 / 3)
}

func (vs *VoteSet) HasQuorum() bool {
	return vs.hasQuorum(vs.sum)
}

func (vs *VoteSet) QuorumBlock() *crypto.Hash {
	return vs.quorum
}

func (vs *VoteSet) HasQuorumBlock(blockHash crypto.Hash) bool {
	blockVotes, exists := vs.votesByBlock[blockHash]
	if !exists {
		return false
	}

	return vs.hasQuorum(blockVotes.sum)
}

func (vs *VoteSet) ToCommit() *block.Commit {
	if vs.voteType != VoteTypePrecommit {
		return nil
	}
	blockHash := vs.quorum
	if blockHash == nil || blockHash.IsUndef() {
		return nil
	}

	votesMap := vs.votesByBlock[*blockHash].votes
	vals := vs.valSet.Validators()
	committers := make([]block.Committer, len(vals))
	sigs := make([]*crypto.Signature, 0)

	for i, addr := range vals {
		status := block.CommitNotSigned
		v := votesMap[addr]

		if v != nil {
			sigs = append(sigs, v.Signature())
			status = block.CommitSigned
		}

		committers[i].Address = addr
		committers[i].Status = status
	}

	sig := crypto.Aggregate(sigs)

	return block.NewCommit(vs.round, committers, sig)
}

func (vs *VoteSet) Fingerprint() string {
	return fmt.Sprintf("{%v/%v/%s SUM:%v}", vs.height, vs.round, vs.voteType, vs.sum)
}
