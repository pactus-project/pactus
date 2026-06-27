package param

import (
	"testing"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/assert"
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
