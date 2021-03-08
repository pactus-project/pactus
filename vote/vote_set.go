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
	votes map[crypto.Address]*Vote
	power int64
}

func newBlockVotes() *blockVotes {
	return &blockVotes{
		votes: make(map[crypto.Address]*Vote),
		power: 0,
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

	return true
}

type VoteSet struct {
	height           int
	round            int
	voteType         VoteType
	validators       []*validator.Validator
	blockVotes       map[crypto.Hash]*blockVotes
	totalPower       int64
	accumulatedPower int64
	quorumBlock      *crypto.Hash
}

func NewVoteSet(height int, round int, voteType VoteType, validators []*validator.Validator) *VoteSet {

	totalPower := int64(0)
	for _, val := range validators {
		totalPower += val.Power()
	}

	return &VoteSet{
		height:     height,
		round:      round,
		voteType:   voteType,
		validators: validators,
		totalPower: totalPower,
		blockVotes: make(map[crypto.Hash]*blockVotes),
	}
}

func (vs *VoteSet) Type() VoteType { return vs.voteType }
func (vs *VoteSet) Height() int    { return vs.height }
func (vs *VoteSet) Round() int     { return vs.round }
func (vs *VoteSet) Power() int64   { return vs.accumulatedPower }

func (vs *VoteSet) Len() int {
	sum := 0
	for _, bv := range vs.blockVotes {
		sum += len(bv.votes)
	}
	return sum
}

func (vs *VoteSet) AllVotes() []*Vote {
	votes := make([]*Vote, 0)

	for _, blockVotes := range vs.blockVotes {
		for _, vote := range blockVotes.votes {
			votes = append(votes, vote)
		}
	}
	return votes
}

func (vs *VoteSet) getValidatorByAddress(addr crypto.Address) *validator.Validator {
	for _, val := range vs.validators {
		if val.Address().EqualsTo(addr) {
			return val
		}
	}
	return nil
}

func (vs *VoteSet) AddVote(vote *Vote) (bool, error) {
	signer := vote.Signer()
	blockHash := vote.BlockHash()

	if (vote.data.Height != vs.height) ||
		(vote.data.Round != vs.round) ||
		(vote.data.VoteType != vs.voteType) {
		return false, errors.Errorf(errors.ErrInvalidVote, "Expected %d/%d/%s, but got %d/%d/%s",
			vs.height, vs.round, vs.voteType,
			vote.Height(), vote.Round(), vote.VoteType())
	}

	val := vs.getValidatorByAddress(signer)
	if val == nil {
		return false, errors.Errorf(errors.ErrInvalidVote, "Cannot find validator %s in committee", signer)
	}

	if err := vote.Verify(val.PublicKey()); err != nil {
		return false, errors.Errorf(errors.ErrInvalidVote, "Failed to verify vote")
	}

	bv, exists := vs.blockVotes[blockHash]
	if !exists {
		bv = newBlockVotes()
		vs.blockVotes[blockHash] = bv
	}

	// check for conflict
	for id, bv := range vs.blockVotes {
		if id != vote.data.BlockHash {
			duplicated, ok := bv.votes[signer]

			if ok {
				// A possible scenario:
				// A peer doesn't have a proposal, he votes for undef.
				// Later he receives the proposal, so he vote again.
				// We should ignore undef vote
				if duplicated.BlockHash().IsUndef() {
					// Remove undef vote and replace it with new vote
					vs.accumulatedPower -= val.Power()
					bv.power -= val.Power()
					delete(bv.votes, signer)
					vs.quorumBlock = nil
				} else if vote.BlockHash().IsUndef() {
					// Because of network latency, we might receive undef vote after block vote.
					// Ignore undef vote in this case.
					return false, nil
				} else if duplicated.BlockHash() != blockHash {
					// Duplicated vote:
					// 1- Same signer
					// 2- Previous blockhash is not undef
					// 3- Block hashes are different
					//
					// We report an error and remove the previous vote
					//
					vs.accumulatedPower -= val.Power()
					bv.power -= val.Power()
					delete(bv.votes, signer)

					return false, errors.Error(errors.ErrDuplicateVote)
				}
			}
		}
	}

	added := bv.addVote(vote)
	if added {
		vs.accumulatedPower += val.Power()
		bv.power += val.Power()
		if vs.hasQuorum(bv.power) {
			vs.quorumBlock = &blockHash
		}

	}

	return added, nil
}
func (vs *VoteSet) hasQuorum(power int64) bool {
	return power > (vs.totalPower * 2 / 3)
}

func (vs *VoteSet) HasQuorum() bool {
	return vs.hasQuorum(vs.accumulatedPower)
}

func (vs *VoteSet) QuorumBlock() *crypto.Hash {
	return vs.quorumBlock
}

func (vs *VoteSet) ToCommit() *block.Commit {
	if vs.voteType != VoteTypePrecommit {
		return nil
	}
	blockHash := vs.quorumBlock
	if blockHash == nil || blockHash.IsUndef() {
		return nil
	}

	votesMap := vs.blockVotes[*blockHash].votes
	committers := make([]block.Committer, len(vs.validators))
	sigs := make([]crypto.Signature, 0)

	for i, val := range vs.validators {
		status := block.CommitNotSigned
		v := votesMap[val.Address()]

		if v != nil {
			sigs = append(sigs, *v.Signature())
			status = block.CommitSigned
		}

		committers[i].Number = val.Number()
		committers[i].Status = status
	}

	sig := crypto.Aggregate(sigs)

	return block.NewCommit(*blockHash, vs.round, committers, sig)
}

func (vs *VoteSet) Fingerprint() string {
	return fmt.Sprintf("{%v/%v/%s SUM:%v}", vs.height, vs.round, vs.voteType, vs.Len())
}
