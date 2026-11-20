package voteset

import (
	"bytes"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
)

type voteSet struct {
	round      types.Round
	validators map[crypto.Address]*validator.Validator
	totalPower int64
}

func newVoteSet(round types.Round, totalPower int64,
	validators map[crypto.Address]*validator.Validator,
) *voteSet {
	return &voteSet{
		round:      round,
		validators: validators,
		totalPower: totalPower,
	}
}

// Round returns the round number for the VoteSet.
func (vs *voteSet) Round() types.Round {
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

func (vs *voteSet) isTwoThirdOfTotalPower(power int64) bool {
	return power > (vs.totalPower * 2 / 3)
}

func (vs *voteSet) isOneThirdOfTotalPower(power int64) bool {
	return power > (vs.totalPower * 1 / 3)
}

// sameSignedContent checks if two votes of the same signer say the same thing.
// The justification of a vote isn't signed: the same signed vote with another valid
// justification is not an equivocation.
func sameSignedContent(vote1, vote2 *vote.Vote) bool {
	return bytes.Equal(vote1.SignBytes(), vote2.SignBytes())
}
