package voteset

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/errors"
)

type VoteSet struct {
	round      int16
	voteType   vote.Type
	validators []*validator.Validator
	blockVotes map[hash.Hash]*blockVotes
	allVotes   map[hash.Hash]*vote.Vote
	totalPower int64
	quorumHash *hash.Hash
}

func NewVoteSet(round int16, voteType vote.Type, validators []*validator.Validator) *VoteSet {
	totalPower := int64(0)
	for _, val := range validators {
		totalPower += val.Power()
	}

	return &VoteSet{
		round:      round,
		voteType:   voteType,
		validators: validators,
		totalPower: totalPower,
		blockVotes: make(map[hash.Hash]*blockVotes),
		allVotes:   make(map[hash.Hash]*vote.Vote),
	}
}

func (vs *VoteSet) Type() vote.Type { return vs.voteType }
func (vs *VoteSet) Round() int16    { return vs.round }

func (vs *VoteSet) Len() int {
	return len(vs.allVotes)
}

func (vs *VoteSet) AllVotes() []*vote.Vote {
	votes := make([]*vote.Vote, 0)
	for _, v := range vs.allVotes {
		votes = append(votes, v)
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

func (vs *VoteSet) mustGetBlockVotes(blockhash hash.Hash) *blockVotes {
	bv, exists := vs.blockVotes[blockhash]
	if !exists {
		bv = newBlockVotes()
		vs.blockVotes[blockhash] = bv
	}
	return bv
}

func (vs *VoteSet) AddVote(v *vote.Vote) error {
	if (v.Round() != vs.Round()) ||
		(v.Type() != vs.Type()) {
		return errors.Errorf(errors.ErrInvalidVote, "expected %d/%s, but got %d/%s",
			vs.Round(), vs.Type(),
			v.Round(), v.Type())
	}

	signer := v.Signer()
	val := vs.getValidatorByAddress(signer)
	if val == nil {
		return errors.Errorf(errors.ErrInvalidVote, "cannot find validator %s in committee", signer)
	}

	if err := v.Verify(val.PublicKey()); err != nil {
		return errors.Errorf(errors.ErrInvalidVote, "failed to verify vote")
	}

	_, exists := vs.allVotes[v.Hash()]
	if exists {
		return errors.Errorf(errors.ErrInvalidVote, "existing vote")
	}

	// Alright! We don't have this vote yet
	vs.allVotes[v.Hash()] = v

	// Now check for duplicity
	for h, bv := range vs.blockVotes {
		if !h.EqualsTo(v.BlockHash()) {
			_, ok := bv.votes[signer]
			if ok {
				// Duplicated vote:
				// 1- Same signer
				// 2- Both votes are different
				//
				// We report an error
				//
				return errors.Error(errors.ErrDuplicateVote)
			}
		}
	}

	blockVotes := vs.mustGetBlockVotes(v.BlockHash())
	blockVotes.addVote(v)
	blockVotes.power += val.Power()
	if vs.hasTwoThirdOfTotalPower(blockVotes.power) {
		hash := v.BlockHash()
		vs.quorumHash = &hash
	}

	return nil
}
func (vs *VoteSet) hasTwoThirdOfTotalPower(power int64) bool {
	return power > (vs.totalPower * 2 / 3)
}

func (vs *VoteSet) hasOneThirdOfTotalPower(power int64) bool {
	return power > (vs.totalPower * 1 / 3)
}

func (vs *VoteSet) BlockHashHasOneThirdOfTotalPower(hash hash.Hash) bool {
	blockVotes := vs.mustGetBlockVotes(hash)
	return vs.hasOneThirdOfTotalPower(blockVotes.power)
}

func (vs *VoteSet) QuorumHash() *hash.Hash {
	return vs.quorumHash
}

func (vs *VoteSet) ToCertificate() *block.Certificate {
	if vs.Type() != vote.VoteTypePrecommit {
		return nil
	}
	blockHash := vs.quorumHash
	if blockHash == nil || blockHash.IsUndef() {
		return nil
	}

	votesMap := vs.blockVotes[*blockHash].votes
	committers := make([]int32, len(vs.validators))
	absentees := make([]int32, 0)
	sigs := make([]*bls.Signature, 0)

	for i, val := range vs.validators {
		v := votesMap[val.Address()]

		if v != nil {
			sigs = append(sigs, v.Signature())
		} else {
			absentees = append(absentees, val.Number())
		}

		committers[i] = val.Number()
	}

	sig := bls.SignatureAggregate(sigs)

	return block.NewCertificate(vs.Round(), committers, absentees, sig)
}

func (vs *VoteSet) Fingerprint() string {
	return fmt.Sprintf("{%v/%s SUM:%v}", vs.round, vs.voteType, vs.Len())
}
