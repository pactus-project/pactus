//go:build gtk

package model

import (
	"errors"

	"github.com/pactus-project/pactus/types/amount"
)

var (
	// ErrInsufficientBalance is returned when the balance does not cover the fee.
	ErrInsufficientBalance = errors.New("the balance does not cover the fee")

	// ErrValidatorFull is returned when the validator has reached the maximum stake.
	ErrValidatorFull = errors.New("the validator has reached the maximum stake")
)

// MaxTransferAmount returns the largest amount that can be transferred from
// the balance, after paying the fee.
func MaxTransferAmount(balance, fee amount.Amount) (amount.Amount, error) {
	if balance <= fee {
		return 0, ErrInsufficientBalance
	}

	return balance - fee, nil
}

// MaxBondAmount returns the largest stake that can be bonded from the balance,
// after paying the fee, without going over the maximum stake of the validator.
func MaxBondAmount(balance, fee, validatorStake, maximumStake amount.Amount) (amount.Amount, error) {
	if validatorStake >= maximumStake {
		return 0, ErrValidatorFull
	}

	maxAmount, err := MaxTransferAmount(balance, fee)
	if err != nil {
		return 0, err
	}

	return min(maxAmount, maximumStake-validatorStake), nil
}
