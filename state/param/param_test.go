package param

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFoundationAddress(t *testing.T) {
	params := FromGenesis(genesis.MainnetGenesis())

	assert.Equal(t, "pc1z0k5ctvn02hsxvl9t3d2efnkv2d5k46ayfddzxg", params.FoundationAddress(1000).String())
	assert.Equal(t, "pc1zrv84qnh96pmkg2sykedtz5mlu0q72st2l5c6gg", params.FoundationAddress(1001).String())
	assert.Equal(t, "pc1z4tdnddwmxeppa3pcaquxhq4rrc5adcx7t3qfj0", params.FoundationAddress(1051).String())
	assert.Equal(t, "pc1zgqq766nf3782gxvrncv8cvfszdda4ss20y4e7a", params.FoundationAddress(1099).String())
	assert.Equal(t, "pc1z0k5ctvn02hsxvl9t3d2efnkv2d5k46ayfddzxg", params.FoundationAddress(1100).String())
}

func TestRewardCoefficient(t *testing.T) {
	params := FromGenesis(genesis.MainnetGenesis())

	tests := []struct {
		height           types.Height
		coefficient      float64
		blockReward      amount.Amount
		foundationReward amount.Amount
	}{
		{8_000_000, 1.0, 1e9, 3e8},
		{8_000_001, 0.5, 5e8, 15e7},
		{24_000_000, 0.5, 5e8, 15e7},
		{24_000_001, 0.25, 25e7, 75e6},
		{56_000_000, 0.25, 25e7, 75e6},
		{56_000_001, 0.125, 125e6, 375e5},
		{100_000_001, 0.125, 125e6, 375e5},
	}

	for _, tt := range tests {
		assert.InEpsilon(t, tt.coefficient, params.RewardCoefficient(tt.height), 0.0,
			"height %d should have coefficient %f", tt.height, tt.coefficient)
		assert.Equal(t, tt.blockReward, params.BlockReward(tt.height))
		assert.Equal(t, tt.foundationReward, params.FoundationReward(tt.height))
	}
}

func TestFromGenesis(t *testing.T) {
	tests := []struct {
		name      string
		genDoc    *genesis.Genesis
		chainType genesis.ChainType
	}{
		{"mainnet", genesis.MainnetGenesis(), genesis.Mainnet},
		{"testnet", genesis.TestnetGenesis(), genesis.Testnet},
		{"localnet", genesis.MakeGenesis(time.Now(), nil, nil, genesis.DefaultGenesisParams()), genesis.Localnet},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.chainType == genesis.Testnet {
				crypto.ToTestnetHRP()
				t.Cleanup(crypto.ToMainnetHRP)
			}

			require.Equal(t, tt.chainType, tt.genDoc.ChainType())

			params := FromGenesis(tt.genDoc)
			genParams := tt.genDoc.Params()

			assert.Equal(t, genParams.BlockVersion, params.BlockVersion)
			assert.Equal(t, genParams.BlockIntervalInSecond, params.BlockIntervalInSecond)
			assert.Equal(t, genParams.CommitteeSize, params.CommitteeSize)
			assert.Equal(t, genParams.BondInterval, params.BondInterval)
			assert.Equal(t, genParams.UnbondInterval, params.UnbondInterval)
			assert.Equal(t, genParams.SortitionInterval, params.SortitionInterval)
			assert.Equal(t, genParams.TransactionToLiveInterval, params.TransactionToLiveInterval)
			assert.Equal(t, genParams.MinimumStake, params.MinimumStake)
			assert.Equal(t, genParams.MaximumStake, params.MaximumStake)
			assert.Equal(t, time.Duration(genParams.BlockIntervalInSecond)*time.Second, params.BlockInterval())

			// The foundation addresses are used as recipients of the subsidy transaction.
			// A batch transfer accepts only unique account addresses as recipients.
			assert.Len(t, params.FoundationAddresses, 100)
			seen := make(map[crypto.Address]bool)
			for _, addr := range params.FoundationAddresses {
				assert.True(t, addr.IsAccountAddress(), "%s is not an account address", addr)
				assert.False(t, seen[addr], "%s is duplicated", addr)
				seen[addr] = true
			}
		})
	}
}

// The foundation address rotates with the height, and every address is used once per cycle.
func TestFoundationAddressRotation(t *testing.T) {
	params := FromGenesis(genesis.MainnetGenesis())
	cycle := len(params.FoundationAddresses)

	used := make(map[crypto.Address]int)
	for height := types.Height(1); height <= types.Height(cycle); height++ {
		used[params.FoundationAddress(height)]++
	}

	assert.Len(t, used, cycle)
	for addr, count := range used {
		assert.Equal(t, 1, count, "%s is used %d times in one cycle", addr, count)
	}
}

// The validator reward is what remains of the block reward after the foundation part.
// It must never be negative, otherwise the subsidy transaction can't be created.
func TestRewardsAreConsistent(t *testing.T) {
	params := FromGenesis(genesis.MainnetGenesis())

	for _, height := range []types.Height{1, 8_000_000, 8_000_001, 24_000_001, 56_000_001, 1_000_000_000} {
		blockReward := params.BlockReward(height)
		foundationReward := params.FoundationReward(height)

		assert.Positive(t, blockReward, "height %d", height)
		assert.Positive(t, foundationReward, "height %d", height)
		assert.Greater(t, blockReward, foundationReward, "height %d", height)

		// The maximum share of a delegate owner is paid from the validator part of the reward.
		maxOwnerShare := MaxDelegateOwnerRewardShare.MulF64(params.RewardCoefficient(height))
		assert.Equal(t, blockReward-foundationReward, maxOwnerShare,
			"at height %d the validator part is fully given to the owner for the maximum share", height)
	}
}
