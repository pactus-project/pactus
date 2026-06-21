package manager

import (
	"github.com/pactus-project/pactus/util/testsuite"
)

type FakeConsensusManager struct {
	*MockConsensusManager
}

func NewFakeConsensusManager(ts *testsuite.TestSuite) *FakeConsensusManager {
	return &FakeConsensusManager{
		MockConsensusManager: NewMockConsensusManager(ts.MockController()),
	}
}
