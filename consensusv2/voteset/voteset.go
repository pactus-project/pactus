package voteset

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
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
func (vs *voteSet) verifyVote(vote *vote.Vote) (int64, error) {
	signer := vote.Signer()
	val := vs.validators[signer]
	if val == nil {
		return 0, IneligibleVoterError{
			Address: signer,
		}
	}

	if err := vote.Verify(val.PublicKey()); err != nil {
		return 0, err
	}

	return val.Power(), nil
}

// has3FP1Power checks whether the given power is greater than or equal to 3f+1,
// where f is the maximum faulty power.
func (vs *BlockVoteSet) has3FP1Power(power int64) bool {
	return certificate.Has3FP1Power(vs.totalPower, power)
}

// has2FP1Power checks whether the given power is greater than or equal to 2f+1,
// where f is the maximum faulty power.
func (vs *voteSet) has2FP1Power(power int64) bool {
	return certificate.Has2FP1Power(vs.totalPower, power)
}

// has1FP1Power checks whether the given power is greater than or equal to f+1,
// where f is the maximum faulty power.
func (vs *voteSet) has1FP1Power(power int64) bool {
	return certificate.Has1FP1Power(vs.totalPower, power)
}
