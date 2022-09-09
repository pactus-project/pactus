package sandbox

import (
	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
)

var _ Sandbox = &MockSandbox{}

// MockSandbox is a testing mock for sandbox.
type MockSandbox struct {
	TestParams           param.Params
	TestStore            *store.MockStore
	TestCommittee        committee.Committee
	TestCommitteeSigners []crypto.Signer
	AcceptTestSortition  bool
}

func MockingSandbox() *MockSandbox {
	committee, signers := committee.GenerateTestCommittee(7)

	sb := &MockSandbox{
		TestParams:           param.DefaultParams(),
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
func (m *MockSandbox) CurrentHeight() uint32 {
	return m.TestStore.LastHeight + 1
}
func (m *MockSandbox) Params() param.Params {
	return m.TestParams
}
func (m *MockSandbox) FindBlockHashByStamp(stamp hash.Stamp) (hash.Hash, bool) {
	return m.TestStore.FindBlockHashByStamp(stamp)
}
func (m *MockSandbox) FindBlockHeightByStamp(stamp hash.Stamp) (uint32, bool) {
	return m.TestStore.FindBlockHeightByStamp(stamp)
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
