package store

import "github.com/pactus-project/pactus/util/testsuite"

type FakeStore struct {
	*MockStore
}

func NewFakeStore(ts *testsuite.TestSuite) *FakeStore {
	fake := &FakeStore{
		MockStore: NewMockStore(ts.MockController()),
	}

	fake.EXPECT().TotalAccounts().Return(ts.RandInt32()).AnyTimes()
	fake.EXPECT().TotalValidators().Return(ts.RandInt32()).AnyTimes()

	return fake
}
