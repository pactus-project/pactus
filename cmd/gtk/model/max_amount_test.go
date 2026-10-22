//go:build gtk

package model

import (
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaxTransferAmount(t *testing.T) {
	maxAmount, err := MaxTransferAmount(amount.Amount(5_000_000_000), amount.Amount(10_000_000))
	require.NoError(t, err)
	assert.Equal(t, amount.Amount(4_990_000_000), maxAmount)

	_, err = MaxTransferAmount(amount.Amount(10_000_000), amount.Amount(10_000_000))
	require.ErrorIs(t, err, ErrInsufficientBalance)

	_, err = MaxTransferAmount(0, amount.Amount(10_000_000))
	require.ErrorIs(t, err, ErrInsufficientBalance)
}

func TestMaxBondAmount(t *testing.T) {
	fee := amount.Amount(10_000_000)             // 0.01 PAC
	maxStake := amount.Amount(1_000_000_000_000) // 1000 PAC

	// Limited by the balance.
	maxAmount, err := MaxBondAmount(amount.Amount(500_000_000_000), fee, 0, maxStake)
	require.NoError(t, err)
	assert.Equal(t, amount.Amount(499_990_000_000), maxAmount)

	// Limited by the stake the validator can still receive.
	maxAmount, err = MaxBondAmount(amount.Amount(500_000_000_000), fee, amount.Amount(900_000_000_000), maxStake)
	require.NoError(t, err)
	assert.Equal(t, amount.Amount(100_000_000_000), maxAmount)

	// The validator is full.
	_, err = MaxBondAmount(amount.Amount(500_000_000_000), fee, maxStake, maxStake)
	require.ErrorIs(t, err, ErrValidatorFull)

	// The balance does not cover the fee.
	_, err = MaxBondAmount(amount.Amount(5_000_000), fee, 0, maxStake)
	require.ErrorIs(t, err, ErrInsufficientBalance)
}

func TestMaxAmountRoundTrip(t *testing.T) {
	// The Max button writes the amount with Format, and the dialog reads it
	// back with FromString, so the value must not change on the way.
	for _, amt := range []amount.Amount{499_990_000_000, 1, 123_456_789_012_345} {
		parsed, err := amount.FromString(amt.Format())
		require.NoError(t, err)
		assert.Equal(t, amt, parsed, "amount: %d", amt)
	}
}
