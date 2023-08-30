package voteset

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/errors"
)

type voteSet struct {
	voteType   vote.Type
	round      int16
	validators map[crypto.Address]*validator.Validator
	totalPower int64
}

func newVoteSet(voteType vote.Type, round int16, totalPower int64,
	validators map[crypto.Address]*validator.Validator,
) *voteSet {
	return &voteSet{
		round:      round,
		voteType:   voteType,
		validators: validators,
		totalPower: totalPower,
	}
}

func (vs *voteSet) Type() vote.Type {
	return vs.voteType
}

// Round returns the round number for the VoteSet.
func (vs *voteSet) Round() int16 {
	return vs.round
}

// verifyVote checks if the given vote is valid.
// It returns the voting power of if valid, or an error if not.
func (vs *voteSet) verifyVote(v *vote.Vote) (int64, error) {
	signer := v.Signer()
	val := vs.validators[signer]
	if val == nil {
		return 0, errors.Errorf(errors.ErrInvalidAddress,
			"cannot find validator %s in committee", signer)
	}

	if err := v.Verify(val.PublicKey()); err != nil {
		return 0, errors.Errorf(errors.ErrInvalidSignature,
			"failed to verify vote")
	}

	return val.Power(), nil
}

func (vs *voteSet) isTwoThirdOfTotalPower(power int64) bool {
	return power > (vs.totalPower * 2 / 3)
}

func (vs *voteSet) isOneThirdOfTotalPower(power int64) bool {
	return power > (vs.totalPower * 1 / 3)
}

func (vs *voteSet) String() string {
	return fmt.Sprintf("{%v/%s TOTAL:%v}", vs.round, vs.voteType, vs.totalPower)
}
