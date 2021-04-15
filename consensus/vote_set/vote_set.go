package vote_set

import (
	"fmt"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/validator"
)

type VoteSet struct {
	height     int
	round      int
	voteType   vote.VoteType
	validators []*validator.Validator
	blockVotes map[crypto.Hash]*blockVotes
	allVotes   map[crypto.Hash]*vote.Vote
	totalPower int64
	quorumHash *crypto.Hash
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
		allVotes:   make(map[crypto.Hash]*vote.Vote),
	}
}

func (vs *VoteSet) VoteType() vote.VoteType { return vs.voteType }
func (vs *VoteSet) Height() int             { return vs.height }
func (vs *VoteSet) Round() int              { return vs.round }

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

func (vs *VoteSet) mustGetBlockVotes(blockhash crypto.Hash) *blockVotes {
	bv, exists := vs.blockVotes[blockhash]
	if !exists {
		bv = newBlockVotes()
		vs.blockVotes[blockhash] = bv
	}
	return bv
}

func (vs *VoteSet) AddVote(v *vote.Vote) error {
	if (v.Height() != vs.Height()) ||
		(v.Round() != vs.Round()) ||
		(v.VoteType() != vs.VoteType()) {
		return errors.Errorf(errors.ErrInvalidVote, "expected %d/%d/%s, but got %d/%d/%s",
			vs.Height(), vs.Round(), vs.VoteType(),
			v.Height(), v.Round(), v.VoteType())
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

func (vs *VoteSet) QuorumHash() *crypto.Hash {
	return vs.quorumHash
}

func (vs *VoteSet) ToCertificate() *block.Certificate {
	if vs.VoteType() != vote.VoteTypePrecommit {
		return nil
	}
	blockHash := vs.quorumHash
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

	return block.NewCertificate(*blockHash, vs.Round(), committers, absences, sig)
}

func (vs *VoteSet) Fingerprint() string {
	return fmt.Sprintf("{%v/%v/%s SUM:%v}", vs.height, vs.round, vs.voteType, vs.Len())
}
