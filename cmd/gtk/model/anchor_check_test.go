//go:build gtk

package model

import (
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func pac(v float64) amount.Amount {
	return amount.Amount(v * amount.NanoPACPerPAC)
}

func TestCheckAnchorSet(t *testing.T) {
	none := CurrentAnchor{}
	existing := CurrentAnchor{Exists: true, Type: 1, Locked: pac(1)}
	fee := pac(0.01)

	t.Run("new anchor with the minimum", func(t *testing.T) {
		check := CheckAnchorSet(none, 1, pac(1), fee, pac(5))
		require.NoError(t, check.Err)
		assert.False(t, check.Replaces)
		assert.False(t, check.LargeDeposit)
		assert.Equal(t, pac(1), check.LockedAfter)
	})

	t.Run("new anchor below the minimum", func(t *testing.T) {
		check := CheckAnchorSet(none, 1, pac(0.5), fee, pac(5))
		assert.ErrorIs(t, check.Err, ErrAnchorDepositTooSmall)
	})

	t.Run("update adds nothing", func(t *testing.T) {
		check := CheckAnchorSet(existing, 1, 0, fee, pac(1))
		require.NoError(t, check.Err)
		assert.True(t, check.Replaces)
		assert.False(t, check.TypeChanges)
		assert.Equal(t, pac(1), check.LockedAfter)
	})

	t.Run("update with another type", func(t *testing.T) {
		check := CheckAnchorSet(existing, 2, 0, fee, pac(1))
		assert.True(t, check.TypeChanges)
	})

	t.Run("large deposit is highlighted", func(t *testing.T) {
		check := CheckAnchorSet(existing, 1, pac(99), fee, pac(200))
		require.NoError(t, check.Err)
		assert.True(t, check.LargeDeposit)
		assert.Equal(t, pac(100), check.LockedAfter)
	})

	t.Run("balance does not cover deposit and fee", func(t *testing.T) {
		check := CheckAnchorSet(none, 1, pac(1), fee, pac(1))
		assert.ErrorIs(t, check.Err, ErrAnchorInsufficientFunds)
	})

	t.Run("negative deposit", func(t *testing.T) {
		check := CheckAnchorSet(existing, 1, -1, fee, pac(1))
		assert.ErrorIs(t, check.Err, ErrAnchorNegativeDeposit)
	})
}

func TestCheckAnchorDelete(t *testing.T) {
	fee := pac(0.01)

	require.ErrorIs(t, CheckAnchorDelete(CurrentAnchor{}, fee).Err, ErrAnchorNotFound)

	low := CurrentAnchor{Exists: true, Locked: pac(0.005)}
	require.ErrorIs(t, CheckAnchorDelete(low, fee).Err, ErrAnchorLockedBelowFee)

	check := CheckAnchorDelete(CurrentAnchor{Exists: true, Locked: pac(100)}, fee)
	require.NoError(t, check.Err)
	assert.Equal(t, pac(100)-fee, check.Refund)

	zero := CheckAnchorDelete(CurrentAnchor{Exists: true, Locked: fee}, fee)
	require.NoError(t, zero.Err)
	assert.Equal(t, amount.Amount(0), zero.Refund)
}

func TestDefaultAnchorDeposit(t *testing.T) {
	none := CurrentAnchor{}
	existing := CurrentAnchor{Exists: true, Locked: pac(1)}
	target := pac(100)

	assert.Equal(t, MinAnchorDeposit, DefaultAnchorDeposit(none, nil))
	assert.Equal(t, amount.Amount(0), DefaultAnchorDeposit(existing, nil))
	assert.Equal(t, pac(100), DefaultAnchorDeposit(none, &target))
	assert.Equal(t, pac(99), DefaultAnchorDeposit(existing, &target))
}
