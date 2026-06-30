package committee

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"go.uber.org/mock/gomock"
)

type FakeCommittee struct {
	*MockCommittee

	FakeValidators []*validator.Validator
}

func NewFakeCommittee(ts *testsuite.TestSuite) *FakeCommittee {
	fake := &FakeCommittee{
		MockCommittee: NewMockCommittee(ts.MockController()),
	}

	fake.EXPECT().Validators().DoAndReturn(
		func() []*validator.Validator {
			return fake.FakeValidators
		}).AnyTimes()

	fake.EXPECT().Size().DoAndReturn(
		func() int {
			return len(fake.FakeValidators)
		}).AnyTimes()

	fake.EXPECT().Contains(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) bool {
			for _, val := range fake.FakeValidators {
				if val.Address() == addr {
					return true
				}
			}

			return false
		}).AnyTimes()

	return fake
}
