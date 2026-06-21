package sync

import (
	"github.com/pactus-project/pactus/util/testsuite"
)

type FakeSync struct {
	*MockSync
}

func NewFakeSync(ts *testsuite.TestSuite) *FakeSync {
	return &FakeSync{
		MockSync: MockingSync(ts),
	}
}
