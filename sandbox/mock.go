package sandbox

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/validator"
)

var _ Sandbox = &MockSandbox{}

// MockSandbox is a testing mock for sandbox
type MockSandbox struct {
	Params               param.Params
	TestStore            *store.MockStore
	TestCommittee        committee.Committee
	TestCommitteeSigners []crypto.Signer
	AcceptTestSortition  bool
}

func MockingSandbox() *MockSandbox {
	committee, signers := committee.GenerateTestCommittee(7)

	sb := &MockSandbox{
		Params:               param.DefaultParams(),
		TestStore:            store.MockingStore(),
		TestCommittee:        committee,
		TestCommitteeSigners: signers,
	}

	treasuryAmt := int64(21000000 * 1e8)

	for i, val := range committee.Validators() {
		acc := account.NewAccount(val.Address(), int32(i+1))
		acc.AddToBalance(100 * 1e8)
		sb.UpdateAccount(acc)
		sb.UpdateValidator(val)

		treasuryAmt -= val.Stake()
		treasuryAmt -= acc.Balance()
	}
	acc0 := account.NewAccount(crypto.TreasuryAddress, 0)
	acc0.AddToBalance(treasuryAmt)
	sb.UpdateAccount(acc0)

	return sb
}

func (m *MockSandbox) Account(addr crypto.Address) *account.Account {
	acc, _ := m.TestStore.Account(addr)
	return acc
}
func (m *MockSandbox) MakeNewAccount(addr crypto.Address) *account.Account {
	return account.NewAccount(addr, m.TestStore.TotalAccounts())
}
func (m *MockSandbox) UpdateAccount(acc *account.Account) {
	m.TestStore.UpdateAccount(acc)
}
func (m *MockSandbox) Validator(addr crypto.Address) *validator.Validator {
	val, _ := m.TestStore.Validator(addr)
	return val
}
func (m *MockSandbox) MakeNewValidator(pub *bls.PublicKey) *validator.Validator {
	return validator.NewValidator(pub, m.TestStore.TotalValidators())
}
func (m *MockSandbox) UpdateValidator(val *validator.Validator) {
	m.TestStore.UpdateValidator(val)

}
func (m *MockSandbox) CurrentHeight() int32 {
	return m.TestStore.LastHeight + 1
}
func (m *MockSandbox) TransactionToLiveInterval() int32 {
	return m.Params.TransactionToLiveInterval
}
func (m *MockSandbox) FeeFraction() float64 {
	return m.Params.FeeFraction
}
func (m *MockSandbox) MinFee() int64 {
	return m.Params.MinimumFee
}

func (m *MockSandbox) CommitteeSize() int {
	return m.Params.CommitteeSize
}
func (m *MockSandbox) UnbondInterval() int32 {
	return m.Params.UnbondInterval
}
func (m *MockSandbox) BondInterval() int32 {
	return m.Params.BondInterval
}
func (m *MockSandbox) BlockHashByStamp(stamp hash.Stamp) hash.Hash {
	return m.TestStore.BlockHashByStamp(stamp)
}
func (m *MockSandbox) BlockHeightByStamp(stamp hash.Stamp) int32 {
	return m.TestStore.BlockHeightByStamp(stamp)
}
func (m *MockSandbox) IterateAccounts(consumer func(*AccountStatus)) {
	m.TestStore.IterateAccounts(func(acc *account.Account) bool {
		consumer(&AccountStatus{
			Account: *acc,
			Updated: true,
		})
		return false
	})
}
func (m *MockSandbox) IterateValidators(consumer func(*ValidatorStatus)) {
	m.TestStore.IterateValidators(func(val *validator.Validator) bool {
		consumer(&ValidatorStatus{
			Validator: *val,
			Updated:   true,
		})
		return false
	})
}

func (m *MockSandbox) Committee() committee.Reader {
	return m.TestCommittee
}

func (m *MockSandbox) VerifyProof(hash.Stamp, sortition.Proof, *validator.Validator) bool {
	return m.AcceptTestSortition
}
