package pending_votes

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

type VoteSet struct {
	height           int
	round            int
	voteType         vote.VoteType
	validators       []*validator.Validator
	blockVotes       map[crypto.Hash]*blockVotes
	totalPower       int64
	accumulatedPower int64
	quorumBlock      *crypto.Hash
}

func NewVoteSet(height int, round int, voteType vote.VoteType, validators []*validator.Validator) *VoteSet {

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

func (vs *VoteSet) Type() vote.VoteType { return vs.voteType }
func (vs *VoteSet) Height() int         { return vs.height }
func (vs *VoteSet) Round() int          { return vs.round }
func (vs *VoteSet) Power() int64        { return vs.accumulatedPower }

func (vs *VoteSet) Len() int {
	sum := 0
	for _, bv := range vs.blockVotes {
		sum += len(bv.votes)
	}
	return sum
}

func (vs *VoteSet) AllVotes() []*vote.Vote {
	votes := make([]*vote.Vote, 0)

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

func (vs *VoteSet) AddVote(vote *vote.Vote) (bool, error) {
	signer := vote.Signer()
	blockHash := vote.BlockHash()

	if (vote.Height() != vs.height) ||
		(vote.Round() != vs.round) ||
		(vote.VoteType() != vs.voteType) {
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

	var err error
	var duplicated bool
	// check for conflict
	for id, bv := range vs.blockVotes {
		if !id.EqualsTo(vote.BlockHash()) {
			anotherVote, ok := bv.votes[signer]

			if ok {
				// A possible scenario:
				// A peer doesn't have a proposal, he votes for undef.
				// Later he receives the proposal, so he vote again.
				// We should ignore undef vote
				if anotherVote.BlockHash().IsUndef() {
					// Remove undef vote and replace it with new vote
					duplicated = true
				} else if vote.BlockHash().IsUndef() {
					// Because of network latency, we might receive undef vote after block vote.
					// Ignore undef vote in this case.
					return false, nil
				} else if anotherVote.BlockHash() != blockHash {
					// Duplicated vote:
					// 1- Same signer
					// 2- Block hashes are not undef and different
					//
					// We report an error but keep both votes
					//
					duplicated = true
					err = errors.Error(errors.ErrDuplicateVote)
				}
			}
		}
	}

	added := bv.addVote(vote)
	if added {
		bv.power += val.Power()
		if vs.hasQuorum(bv.power) {
			vs.quorumBlock = &blockHash
		}
		if !duplicated {
			vs.accumulatedPower += val.Power()
		}
	}

	return added, err
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

func (vs *VoteSet) ToCertificate() *block.Certificate {
	if vs.voteType != vote.VoteTypePrecommit {
		return nil
	}
	blockHash := vs.quorumBlock
	if blockHash == nil || blockHash.IsUndef() {
		return nil
	}

	votesMap := vs.blockVotes[*blockHash].votes
	committers := make([]int, len(vs.validators))
	absences := make([]int, 0)
	sigs := make([]crypto.Signature, 0)

	for i, val := range vs.validators {
		v := votesMap[val.Address()]

		if v != nil {
			sigs = append(sigs, *v.Signature())
		} else {
			absences = append(absences, val.Number())
		}

		committers[i] = val.Number()
	}

	sig := crypto.Aggregate(sigs)

	return block.NewCertificate(*blockHash, vs.round, committers, absences, sig)
}

func (vs *VoteSet) Fingerprint() string {
	return fmt.Sprintf("{%v/%v/%s SUM:%v}", vs.height, vs.round, vs.voteType, vs.Len())
}
