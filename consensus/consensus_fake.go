package consensus

import (
	"github.com/pactus-project/pactus/util/testsuite"
)

type FakeConsensus struct {
	*MockConsensus
}

func NewFakeConsensus(ts *testsuite.TestSuite) *FakeConsensus {
	return &FakeConsensus{
		MockConsensus: NewMockConsensus(ts.MockController()),
	}
}
