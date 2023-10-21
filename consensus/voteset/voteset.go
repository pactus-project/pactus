package voteset

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/errors"
)

type voteSet struct {
	round      int16
	validators map[crypto.Address]*validator.Validator
	totalPower int64
}

func newVoteSet(round int16, totalPower int64,
	validators map[crypto.Address]*validator.Validator,
) *voteSet {
	return &voteSet{
		round:      round,
		validators: validators,
		totalPower: totalPower,
	}
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
