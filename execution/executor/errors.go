package executor

import (
	"errors"
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
)

// ErrInsufficientFunds indicates the balance is low for the transaction.
var ErrInsufficientFunds = errors.New("insufficient funds")

// ErrPublicKeyNotSet indicates that the public key is not set for the initial Bond transaction.
var ErrPublicKeyNotSet = errors.New("public key is not set")

// ErrPublicKeyAlreadySet indicates that the public key has already been set for the given validator.
var ErrPublicKeyAlreadySet = errors.New("public key is set before")

// ErrValidatorUnbonded indicates that the validator has unbonded.
var ErrValidatorUnbonded = errors.New("validator has unbonded")

// ErrBondingPeriod is returned when a validator is in bonding period.
var ErrBondingPeriod = errors.New("validator in in bonding period")

// ErrInvalidSortitionProof indicates the sortition proof is invalid.
var ErrInvalidSortitionProof = errors.New("invalid sortition proof")

// ErrValidatorInCommittee indicates that the validator is in the committee.
var ErrValidatorInCommittee = errors.New("validator is in the committee")

// SmallStakeError is returned when the stake amount is less than the minimum stake.
type SmallStakeError struct {
	Minimum amount.Amount
}

func (e SmallStakeError) Error() string {
	return fmt.Sprintf("stake amount can't be less than %v", e.Minimum.String())
}

// MaximumStakeError is returned when the validator's stake exceeds the maximum stake limit.
type MaximumStakeError struct {
	Maximum amount.Amount
}

func (e MaximumStakeError) Error() string {
	return fmt.Sprintf("validator's stake amount can't be more than %v", e.Maximum.String())
}

// InvalidPayloadTypeError is returned when transaction payload type
// is not valid.
type InvalidPayloadTypeError struct {
	PayloadType payload.Type
}

func (e InvalidPayloadTypeError) Error() string {
	return fmt.Sprintf("unknown payload type: %s",
		e.PayloadType.String())
}

// AccountNotFoundError is raised when the given address has no associated account.
type AccountNotFoundError struct {
	Address crypto.Address
}

func (e AccountNotFoundError) Error() string {
	return fmt.Sprintf("no account found for address: %s", e.Address.String())
}

// ValidatorNotFoundError is raised when the given address has no associated validator.
type ValidatorNotFoundError struct {
	Address crypto.Address
}

func (e ValidatorNotFoundError) Error() string {
	return fmt.Sprintf("no validator found for address: %s", e.Address.String())
}
