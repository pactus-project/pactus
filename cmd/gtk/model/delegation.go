//go:build gtk

package model

import (
	"errors"
	"math"
	"time"

	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/types/amount"
)

const (
	// BlocksPerDay is the number of blocks in a day, with 10-second blocks.
	BlocksPerDay = 8640

	// delegateShareStep is the unit the owner share must be a multiple of (PIP-49).
	delegateShareStep = 32

	// blockReward is the reward of each block that PIP-49 splits.
	blockReward = amount.Amount(1_000_000_000)
)

// MaxDelegateSharePercent is the largest owner share, as a percent of the block reward.
var MaxDelegateSharePercent = param.MaxDelegateOwnerRewardShare.ToPAC() * 100

var (
	// ErrDelegateShareRange is returned when the owner share is not 0 to 70 percent.
	ErrDelegateShareRange = errors.New("the owner share must be 0 to 70 percent of the block reward")
	// ErrDelegateExpiryPast is returned when the delegation would end at or before the current block.
	ErrDelegateExpiryPast = errors.New("the delegation must end after the current block")
	// ErrDelegateValidatorStaked is returned when a delegated bond targets a validator that has stake.
	ErrDelegateValidatorStaked = errors.New(
		"a delegated bond stakes exactly the maximum stake, so the validator must have no stake yet")
)

// DelegationTerm is a preset length of a delegation.
type DelegationTerm struct {
	Label string
	Days  uint32
}

// DelegationTerms are the presets the bond dialog offers. A custom block height
// is offered after them.
var DelegationTerms = []DelegationTerm{
	{Label: "3 months", Days: 91},
	{Label: "6 months", Days: 182},
	{Label: "12 months", Days: 365},
}

// DelegateShareFromPercent turns a percent of the block reward into the owner
// share in nano PAC, rounded down to a multiple of 32 as PIP-49 requires.
func DelegateShareFromPercent(percent float64) (amount.Amount, error) {
	if math.IsNaN(percent) || percent < 0 || percent > MaxDelegateSharePercent {
		return 0, ErrDelegateShareRange
	}

	share := int64(math.Round(percent * float64(blockReward) / 100))
	share -= share % delegateShareStep

	return amount.Amount(share), nil
}

// DelegateSharePercent returns the owner share as a percent of the block reward.
func DelegateSharePercent(share amount.Amount) float64 {
	return float64(share) * 100 / float64(blockReward)
}

// DelegationRewardSplit returns what the owner and the operator receive from
// each block reward. The protocol keeps the rest (0.3 PAC).
func DelegationRewardSplit(share amount.Amount) (owner, operator amount.Amount) {
	return share, param.MaxDelegateOwnerRewardShare - share
}

// DelegationExpiryAfterDays returns the height a delegation that starts at the
// given height ends at after the given number of days.
func DelegationExpiryAfterDays(height, days uint32) uint32 {
	return height + days*BlocksPerDay
}

// CheckDelegationExpiry checks that a delegation ends after the current height.
func CheckDelegationExpiry(currentHeight, expiry uint32) error {
	if expiry <= currentHeight {
		return ErrDelegateExpiryPast
	}

	return nil
}

// CheckDelegatedBondValidator checks that a validator can take a delegated bond:
// the bond stakes exactly the maximum stake, so the validator must have none.
func CheckDelegatedBondValidator(validatorStake amount.Amount) error {
	if validatorStake > 0 {
		return ErrDelegateValidatorStaked
	}

	return nil
}

// EstimateBlockTime estimates when a future block is produced, from the
// current height and time, with 10-second blocks. Past heights return now.
func EstimateBlockTime(currentHeight, height uint32, now time.Time) time.Time {
	if height <= currentHeight {
		return now
	}

	return now.Add(time.Duration(height-currentHeight) * 10 * time.Second)
}

// maxBlocksBehind is how far behind the node can be for the wallet to rely on
// its height, about one minute of blocks.
const maxBlocksBehind = 6

// ChainStatus is the height of the node and how many blocks it still has to sync.
type ChainStatus struct {
	Height     uint32
	BlocksLeft int64
}

// IsSynced reports whether the node is close enough to the tip of the chain
// for its height to be used, for example to set when a delegation ends.
func (s ChainStatus) IsSynced() bool {
	return s.BlocksLeft <= maxBlocksBehind
}
