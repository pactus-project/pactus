package sandbox

import (
	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"go.uber.org/mock/gomock"
)

type FakeSandbox struct {
	*MockSandbox
	*testsuite.TestSuite

	SbxParams    *param.Params
	SbxHeight    types.Height
	Accounts     map[crypto.Address]*account.Account
	Validators   map[crypto.Address]*validator.Validator
	TotalCoin    amount.Amount
	SbxCommittee *committee.MockCommittee
}

func NewFakeSandbox(ts *testsuite.TestSuite) *FakeSandbox {
	mock := NewMockSandbox(ts.MockController())
	params := param.FromGenesis(genesis.MainnetGenesis())
	params.BlockVersion = protocol.ProtocolVersion3

	committee := committee.NewMockCommittee(ts.MockController())

	accounts := make(map[crypto.Address]*account.Account)
	validators := make(map[crypto.Address]*validator.Validator)

	fake := &FakeSandbox{
		MockSandbox:  mock,
		TestSuite:    ts,
		Accounts:     accounts,
		Validators:   validators,
		SbxCommittee: committee,
		SbxParams:    params,
		SbxHeight:    ts.RandHeight(),
	}

	fake.EXPECT().Account(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) *account.Account {
			return accounts[addr]
		},
	).AnyTimes()

	fake.EXPECT().MakeNewAccount(gomock.Any()).DoAndReturn(
		func(crypto.Address) *account.Account {
			return account.NewAccount(ts.RandInt32())
		},
	).AnyTimes()

	fake.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).DoAndReturn(
		func(addr crypto.Address, acc *account.Account) {
			accounts[addr] = acc
		},
	).AnyTimes()

	fake.EXPECT().Validator(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) *validator.Validator {
			return validators[addr]
		},
	).AnyTimes()

	fake.EXPECT().CurrentHeight().DoAndReturn(
		func() types.Height {
			return fake.SbxHeight
		},
	).AnyTimes()

	fake.EXPECT().Validator(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) *validator.Validator {
			return validators[addr]
		},
	).AnyTimes()

	fake.EXPECT().MakeNewValidator(gomock.Any()).DoAndReturn(
		func(pub *bls.PublicKey) *validator.Validator {
			return validator.NewValidator(pub, ts.RandInt32())
		},
	).AnyTimes()

	fake.EXPECT().UpdateValidator(gomock.Any()).DoAndReturn(
		func(val *validator.Validator) {
			validators[val.Address()] = val
		},
	).AnyTimes()

	fake.EXPECT().CurrentHeight().Return(ts.RandHeight()).AnyTimes()
	fake.EXPECT().Params().Return(params).AnyTimes()
	fake.EXPECT().Committee().Return(committee).AnyTimes()

	return fake
}

func (f *FakeSandbox) AddValidator(val *validator.Validator) {
	f.Validators[val.Address()] = val
}

func (f *FakeSandbox) AddAccount(addr crypto.Address, acc *account.Account) {
	f.Accounts[addr] = acc
}
