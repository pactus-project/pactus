package manager

import (
	"github.com/pactus-project/pactus/util/testsuite"
)

type FakeWalletManager struct {
	*MockWalletManager
}

func NewFakeWalletManager(ts *testsuite.TestSuite) *FakeWalletManager {
	return &FakeWalletManager{
		MockWalletManager: NewMockWalletManager(ts.MockController()),
	}
}
