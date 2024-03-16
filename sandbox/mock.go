package sandbox

import (
	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
)

var _ Sandbox = &MockSandbox{}

// MockSandbox is a testing mock for sandbox.
type MockSandbox struct {
	ts *testsuite.TestSuite

	TestParams           *param.Params
	TestStore            *store.MockStore
	TestCommittee        committee.Committee
	TestAcceptSortition  bool
	TestJoinedValidators map[crypto.Address]bool
	TestCommittedTrxs    map[tx.ID]*tx.Tx
	TestPowerDelta       int64
}

func MockingSandbox(ts *testsuite.TestSuite) *MockSandbox {
	cmt, _ := ts.GenerateTestCommittee(7)

	sb := &MockSandbox{
		ts:                   ts,
		TestParams:           param.DefaultParams(),
		TestStore:            store.MockingStore(ts),
		TestCommittee:        cmt,
		TestJoinedValidators: make(map[crypto.Address]bool),
		TestCommittedTrxs:    make(map[tx.ID]*tx.Tx),
	}

	treasuryAmt := amount.Amount(21_000_000 * 1e9)

	for i, val := range cmt.Validators() {
		acc := account.NewAccount(int32(i + 1))
		acc.AddToBalance(100 * 1e9)
		sb.UpdateAccount(val.Address(), acc)
		sb.UpdateValidator(val)

		treasuryAmt -= val.Stake()
		treasuryAmt -= acc.Balance()
	}
	acc0 := account.NewAccount(0)
	acc0.AddToBalance(treasuryAmt)
	sb.UpdateAccount(crypto.TreasuryAddress, acc0)

	return sb
}

func (m *MockSandbox) Account(addr crypto.Address) *account.Account {
	acc, _ := m.TestStore.Account(addr)

	return acc
}

func (m *MockSandbox) MakeNewAccount(_ crypto.Address) *account.Account {
	acc := account.NewAccount(m.TestStore.TotalAccounts())

	return acc
}

func (m *MockSandbox) UpdateAccount(addr crypto.Address, acc *account.Account) {
	m.TestStore.UpdateAccount(addr, acc)
}

func (m *MockSandbox) AnyRecentTransaction(txID tx.ID) bool {
	if m.TestCommittedTrxs[txID] != nil {
		return true
	}

	return m.TestStore.AnyRecentTransaction(txID)
}

func (m *MockSandbox) Validator(addr crypto.Address) *validator.Validator {
	val, _ := m.TestStore.Validator(addr)

	return val
}

func (m *MockSandbox) JoinedToCommittee(addr crypto.Address) {
	m.TestJoinedValidators[addr] = true
}

func (m *MockSandbox) IsJoinedCommittee(addr crypto.Address) bool {
	return m.TestJoinedValidators[addr]
}

func (m *MockSandbox) MakeNewValidator(pub *bls.PublicKey) *validator.Validator {
	val := validator.NewValidator(pub, m.TestStore.TotalValidators())

	return val
}

func (m *MockSandbox) UpdateValidator(val *validator.Validator) {
	m.TestStore.UpdateValidator(val)
}

func (m *MockSandbox) CurrentHeight() uint32 {
	return m.TestStore.LastHeight + 1
}

func (m *MockSandbox) Params() *param.Params {
	return m.TestParams
}

func (m *MockSandbox) IterateAccounts(consumer func(crypto.Address, *account.Account, bool)) {
	m.TestStore.IterateAccounts(func(addr crypto.Address, acc *account.Account) bool {
		consumer(addr, acc, true)

		return false
	})
}

func (m *MockSandbox) IterateValidators(consumer func(*validator.Validator, bool, bool)) {
	m.TestStore.IterateValidators(func(val *validator.Validator) bool {
		consumer(val, true, m.TestJoinedValidators[val.Address()])

		return false
	})
}

func (m *MockSandbox) Committee() committee.Reader {
	return m.TestCommittee
}

func (m *MockSandbox) UpdatePowerDelta(delta int64) {
	m.TestPowerDelta += delta
}

func (m *MockSandbox) PowerDelta() int64 {
	return m.TestPowerDelta
}

func (m *MockSandbox) VerifyProof(uint32, sortition.Proof, *validator.Validator) bool {
	return m.TestAcceptSortition
}

func (m *MockSandbox) CommitTransaction(trx *tx.Tx) {
	m.TestCommittedTrxs[trx.ID()] = trx
}

func (m *MockSandbox) AccumulatedFee() amount.Amount {
	return m.ts.RandAmount()
}
