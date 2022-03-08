package sandbox

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/validator"
)

var _ Sandbox = &MockSandbox{}

// MockSandbox is a testing mock for sandbox
type MockSandbox struct {
	TestAccounts         map[crypto.Address]*account.Account
	TestValidators       map[crypto.Address]*validator.Validator
	TestBlocks           map[int]*block.Block
	CurHeight            int
	Params               param.Params
	TotalAccount         int
	TotalValidator       int
	TestCommittee        committee.Committee
	TestCommitteeSigners []crypto.Signer
	AcceptTestSortition  bool
}

func MockingSandbox() *MockSandbox {
	committee, signers := committee.GenerateTestCommittee(7)

	sb := &MockSandbox{
		TestAccounts:         make(map[crypto.Address]*account.Account),
		TestValidators:       make(map[crypto.Address]*validator.Validator),
		TestBlocks:           make(map[int]*block.Block),
		Params:               param.DefaultParams(),
		TestCommittee:        committee,
		TestCommitteeSigners: signers,
	}

	treasuryAmt := int64(21000000 * 1e8)

	for i, val := range committee.Validators() {
		acc := account.NewAccount(val.Address(), i+1)
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
	acc, ok := m.TestAccounts[addr]
	if !ok {
		return nil
	}
	return acc
}
func (m *MockSandbox) MakeNewAccount(addr crypto.Address) *account.Account {
	a := account.NewAccount(addr, m.TotalAccount)
	m.TotalAccount++
	return a
}
func (m *MockSandbox) UpdateAccount(acc *account.Account) {
	m.TestAccounts[acc.Address()] = acc
}
func (m *MockSandbox) Validator(addr crypto.Address) *validator.Validator {
	val, ok := m.TestValidators[addr]
	if !ok {
		return nil
	}
	return val
}
func (m *MockSandbox) MakeNewValidator(pub *bls.PublicKey) *validator.Validator {
	v := validator.NewValidator(pub, m.TotalAccount)
	m.TotalValidator++
	return v
}
func (m *MockSandbox) UpdateValidator(val *validator.Validator) {
	m.TestValidators[val.Address()] = val

}
func (m *MockSandbox) CurrentHeight() int {
	return m.CurHeight
}
func (m *MockSandbox) TransactionToLiveInterval() int {
	return m.Params.TransactionToLiveInterval
}
func (m *MockSandbox) MaxMemoLength() int {
	return m.Params.MaximumMemoLength
}
func (m *MockSandbox) FeeFraction() float64 {
	return m.Params.FeeFraction
}
func (m *MockSandbox) MinFee() int64 {
	return m.Params.MinimumFee
}

func (m *MockSandbox) AddTestBlock(height int, b *block.Block) {
	m.TestBlocks[height] = b
	m.CurHeight = height + 1
}

func (m *MockSandbox) RandomTestAcc() *account.Account {
	for _, acc := range m.TestAccounts {
		return acc
	}
	panic("no account in sandbox")
}

func (m *MockSandbox) RandomTestVal() *validator.Validator {
	for _, val := range m.TestValidators {
		return val
	}
	panic("no validator in sandbox")
}

func (m *MockSandbox) CommitteeSize() int {
	return m.Params.CommitteeSize
}
func (m *MockSandbox) UnbondInterval() int {
	return m.Params.UnbondInterval
}
func (m *MockSandbox) BondInterval() int {
	return m.Params.CommitteeSize * 2
}
func (m *MockSandbox) BlockHeightByStamp(stamp hash.Stamp) int {
	for i, b := range m.TestBlocks {
		if b.Stamp().EqualsTo(stamp) {
			return i
		}
	}

	return -1
}
func (m *MockSandbox) IterateAccounts(consumer func(*AccountStatus)) {
	for _, acc := range m.TestAccounts {
		as := &AccountStatus{
			Account: *acc,
			Updated: true,
		}
		consumer(as)
	}
}
func (m *MockSandbox) IterateValidators(consumer func(*ValidatorStatus)) {
	for _, val := range m.TestValidators {
		vs := &ValidatorStatus{
			Validator: *val,
			Updated:   true,
		}
		consumer(vs)
	}
}

func (m *MockSandbox) Committee() committee.Reader {
	return m.TestCommittee
}

func (m *MockSandbox) VerifyProof(hash.Stamp, sortition.Proof, *validator.Validator) bool {
	return m.AcceptTestSortition
}
