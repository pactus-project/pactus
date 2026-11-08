package genesis_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/stretchr/testify/assert"
)

func BenchmarkChainType(b *testing.B) {
	gen := genesis.MainnetGenesis()
	b.ReportAllocs()
	for b.Loop() {
		_ = gen.ChainType()
	}
}

// The hashes of the published documents are computed when the package loads,
// before a testnet node switches the address prefix; they must not depend on it.
func TestChainTypeWithTestnetPrefix(t *testing.T) {
	crypto.ToTestnetHRP()
	defer crypto.ToMainnetHRP()

	assert.Equal(t, genesis.Testnet, genesis.TestnetGenesis().ChainType())
	assert.Equal(t, genesis.Mainnet, genesis.MainnetGenesis().ChainType())
}
