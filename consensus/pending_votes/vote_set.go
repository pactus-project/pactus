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

func (vs *VoteSet) Type() vote.VoteType     { return vs.voteType }
func (vs *VoteSet) Height() int             { return vs.height }
func (vs *VoteSet) Round() int              { return vs.round }
func (vs *VoteSet) AccumulatedPower() int64 { return vs.accumulatedPower }

func (vs *VoteSet) Len() int {
	sum := 0
	for _, bv := range vs.blockVotes {
		sum += len(bv.votes)
	}
	return sum
}

func (vs *VoteSet) AllVotes() []*vote.Vote {
	votes := make([]*vote.Vote, 0)

	for _, bv := range vs.blockVotes {
		for _, vote := range bv.votes {
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

func (vs *VoteSet) checkVote(vote *vote.Vote) error {
	if (vote.Height() != vs.height) ||
		(vote.Round() != vs.round) ||
		(vote.VoteType() != vs.voteType) {
		return errors.Errorf(errors.ErrInvalidVote, "Expected %d/%d/%s, but got %d/%d/%s",
			vs.height, vs.round, vs.voteType,
			vote.Height(), vote.Round(), vote.VoteType())
	}

	return nil
}

func (vs *VoteSet) mustGetBlockVotes(blockhash crypto.Hash) *blockVotes {
	bv, exists := vs.blockVotes[blockhash]
	if !exists {
		bv = newBlockVotes()
		vs.blockVotes[blockhash] = bv
	}
	return bv
}

func (vs *VoteSet) AddVote(vote *vote.Vote) (bool, error) {
	if err := vs.checkVote(vote); err != nil {
		return false, err
	}

	signer := vote.Signer()
	val := vs.getValidatorByAddress(signer)
	if val == nil {
		return false, errors.Errorf(errors.ErrInvalidVote, "Cannot find validator %s in committee", signer)
	}

	if err := vote.Verify(val.PublicKey()); err != nil {
		return false, errors.Errorf(errors.ErrInvalidVote, "Failed to verify vote")
	}

	bv := vs.mustGetBlockVotes(vote.BlockHash())

	var err error
	var duplicated bool
	// check for conflicts
	for h, bv := range vs.blockVotes {
		if !h.EqualsTo(vote.BlockHash()) {
			anotherVote, ok := bv.votes[signer]
			if ok {
				if anotherVote.BlockHash().IsUndef() {
					// A possible scenario:
					// A peer doesn't have a proposal, it votes for null.
					// Later it receives the proposal, so it votes again.
					duplicated = true
				} else if vote.BlockHash().IsUndef() {
					// A possible scenario:
					// Because of network latency, we might receive null_vote after block_vote.
					duplicated = true
				} else if anotherVote.BlockHash() != vote.BlockHash() {
					// Duplicated vote:
					// 1- Same signer
					// 2- Both votes are not null
					// 3- Both votes are different
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
		if !duplicated {
			vs.accumulatedPower += val.Power()
		}
		if vs.hasQuorum(bv.power) {
			if vs.quorumBlock == nil || vs.quorumBlock.IsUndef() {
				blockHash := vote.BlockHash()
				vs.quorumBlock = &blockHash
			}
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

func (vs *VoteSet) HasOneThirdOfTotalPower(hash crypto.Hash) bool {
	bv := vs.blockVotes[hash]
	if bv == nil {
		return false
	}
	return bv.power > (vs.totalPower * 1 / 3)
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
