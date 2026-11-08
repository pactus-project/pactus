package genesis_test

import (
	"testing"

	"github.com/pactus-project/pactus/genesis"
)

func BenchmarkChainType(b *testing.B) {
	gen := genesis.MainnetGenesis()
	b.ReportAllocs()
	for b.Loop() {
		_ = gen.ChainType()
	}
}
