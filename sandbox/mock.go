package sandbox

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
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
	Accounts       map[crypto.Address]*account.Account
	Validators     map[crypto.Address]*validator.Validator
	Blocks         map[int]*block.Block
	CurHeight      int
	Params         param.Params
	TotalAccount   int
	TotalValidator int
	InCommittee    bool
}

func MockingSandbox() *MockSandbox {
	return &MockSandbox{
		Accounts:   make(map[crypto.Address]*account.Account),
		Validators: make(map[crypto.Address]*validator.Validator),
		Blocks:     make(map[int]*block.Block),
		Params:     param.DefaultParams(),
	}
}

func (m *MockSandbox) Account(addr crypto.Address) *account.Account {
	acc, ok := m.Accounts[addr]
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
	m.Accounts[acc.Address()] = acc
}
func (m *MockSandbox) Validator(addr crypto.Address) *validator.Validator {
	val, ok := m.Validators[addr]
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
	m.Validators[val.Address()] = val

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

func (m *MockSandbox) AppendTestBlock(height int, b *block.Block) {
	m.Blocks[height] = b
	m.CurHeight = height + 1
}

func (m *MockSandbox) TestAccSeq(a crypto.Address) int {
	if acc, ok := m.Accounts[a]; ok {
		return acc.Sequence()
	}

	panic("invalid account address")
}

func (m *MockSandbox) TestValSeq(a crypto.Address) int {
	if val, ok := m.Validators[a]; ok {
		return val.Sequence()
	}

	panic("invalid validator address")

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
func (m *MockSandbox) IsInCommittee(crypto.Address) bool {
	return m.InCommittee
}
func (m *MockSandbox) CommitteeAge() int {
	return 0
}
func (m *MockSandbox) CommitteePower() int64 {
	return 0
}
func (m *MockSandbox) JoinedPower() int64 {
	return 0
}
func (m *MockSandbox) CommitteeHasFreeSeats() bool {
	return false
}
func (m *MockSandbox) BlockHeightByStamp(stamp hash.Stamp) int {
	for i, b := range m.Blocks {
		if b.Stamp().EqualsTo(stamp) {
			return i
		}
	}

	return -1
}

func (m *MockSandbox) BlockSeedByStamp(stamp hash.Stamp) sortition.VerifiableSeed {
	for _, b := range m.Blocks {
		if b.Stamp().EqualsTo(stamp) {
			return b.Header().SortitionSeed()
		}
	}

	return sortition.UndefVerifiableSeed
}
func (m *MockSandbox) TotalPower() int64 {
	p := int64(0)
	for _, val := range m.Validators {
		p += val.Power()
	}
	return p
}
func (m *MockSandbox) IterateAccounts(consumer func(*AccountStatus)) {

}
func (m *MockSandbox) IterateValidators(consumer func(*ValidatorStatus)) {

}
