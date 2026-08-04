package txpool

import "github.com/pactus-project/pactus/util/testsuite"

type FakeTxPool struct {
	*MockTxPool
}

func NewFakeTxPool(ts *testsuite.TestSuite) *FakeTxPool {
	return &FakeTxPool{
		MockTxPool: NewMockTxPool(ts.MockController()),
	}
}
