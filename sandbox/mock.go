package sandbox

import (
	"fmt"

	"github.com/zarbchain/zarb-go/account"
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
	Accounts           map[crypto.Address]*account.Account
	Validators         map[crypto.Address]*validator.Validator
	HashToHeight       map[hash.Hash]int
	CurHeight          int
	Params             param.Params
	TotalAccount       int
	TotalValidator     int
	AcceptSortition    bool
	WelcomeToCommittee bool
	InCommittee        bool
}

func MockingSandbox() *MockSandbox {
	return &MockSandbox{
		Accounts:        make(map[crypto.Address]*account.Account),
		Validators:      make(map[crypto.Address]*validator.Validator),
		HashToHeight:    make(map[hash.Hash]int),
		Params:          param.DefaultParams(),
		AcceptSortition: true,
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
func (m *MockSandbox) EnterCommittee(hash hash.Hash, addr crypto.Address) error {
	if !m.WelcomeToCommittee {
		return fmt.Errorf("cannot enter to the committee")
	}
	return nil
}
func (m *MockSandbox) VerifySortition(blockHash hash.Hash, proof sortition.Proof, val *validator.Validator) bool {
	return m.AcceptSortition
}
func (m *MockSandbox) CurrentHeight() int {
	return m.CurHeight
}
func (m *MockSandbox) BlockHeight(hash hash.Hash) int {
	h, ok := m.HashToHeight[hash]
	if !ok {
		return -1
	}
	return h
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

func (m *MockSandbox) AppendNewBlock(height int, hash hash.Hash) {
	m.HashToHeight[hash] = height
	m.CurHeight = height + 1
}

func (m *MockSandbox) AccSeq(a crypto.Address) int {
	if acc, ok := m.Accounts[a]; ok {
		return acc.Sequence()
	}

	panic("invalid account address")
}

func (m *MockSandbox) ValSeq(a crypto.Address) int {
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

func (m *MockSandbox) IsInCommittee(crypto.Address) bool {
	return m.InCommittee
}

func (m *MockSandbox) FindBlockInfoByStamp(stamp hash.Stamp) (int, hash.Hash) {
	for h, i := range m.HashToHeight {
		if h.Stamp().EqualsTo(stamp) {
			return i, h
		}
	}

	return -1, hash.UndefHash
}
