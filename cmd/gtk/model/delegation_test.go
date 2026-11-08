//go:build gtk

package model

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelegateShareFromPercent(t *testing.T) {
	tests := []struct {
		percent float64
		want    amount.Amount
	}{
		{0, 0},
		{70, 700_000_000},
		{50, 500_000_000},
		{35, 350_000_000},
		// 12.345% is 123,450,000 nano PAC, rounded down to a multiple of 32.
		{12.345, 123_449_984},
	}

	for _, tt := range tests {
		share, err := DelegateShareFromPercent(tt.percent)
		require.NoError(t, err)
		assert.Equal(t, tt.want, share, "percent %v", tt.percent)
		assert.Zero(t, int64(share)%32)
	}

	for _, bad := range []float64{-1, 70.0001, 100} {
		_, err := DelegateShareFromPercent(bad)
		assert.ErrorIs(t, err, ErrDelegateShareRange, "percent %v", bad)
	}
}

func TestDelegateSharePercent(t *testing.T) {
	assert.InDelta(t, 70.0, DelegateSharePercent(700_000_000), 1e-9)
	assert.InDelta(t, 0.0, DelegateSharePercent(0), 1e-9)
}

func TestDelegationRewardSplit(t *testing.T) {
	owner, operator := DelegationRewardSplit(500_000_000)
	assert.Equal(t, amount.Amount(500_000_000), owner)
	assert.Equal(t, amount.Amount(200_000_000), operator)
}

func TestDelegationExpiry(t *testing.T) {
	assert.Equal(t, uint32(1000+91*BlocksPerDay), DelegationExpiryAfterDays(1000, 91))

	require.NoError(t, CheckDelegationExpiry(1000, 1001))
	require.ErrorIs(t, CheckDelegationExpiry(1000, 1000), ErrDelegateExpiryPast)
	assert.ErrorIs(t, CheckDelegationExpiry(1000, 999), ErrDelegateExpiryPast)
}

func TestCheckDelegatedBondValidator(t *testing.T) {
	require.NoError(t, CheckDelegatedBondValidator(0))
	assert.ErrorIs(t, CheckDelegatedBondValidator(1), ErrDelegateValidatorStaked)
}

func TestEstimateBlockTime(t *testing.T) {
	now := time.Date(2026, 11, 8, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, now, EstimateBlockTime(100, 50, now))
	assert.Equal(t, now.Add(24*time.Hour), EstimateBlockTime(100, 100+BlocksPerDay, now))
}

func TestChainStatusIsSynced(t *testing.T) {
	assert.True(t, ChainStatus{Height: 100, BlocksLeft: 0}.IsSynced())
	assert.True(t, ChainStatus{Height: 100, BlocksLeft: 6}.IsSynced())
	assert.False(t, ChainStatus{Height: 100, BlocksLeft: 7}.IsSynced())
}
