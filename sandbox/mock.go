package sandbox

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

// MockSandbox is a testing mock
type MockSandbox struct {
	Stamps         map[crypto.Hash]int
	CurrentHeight_ int
	TTLInterval    int
	MaxMemoLenght_ int
	FeeFraction_   float64
	MinFee_        int64
}

func NewMockSandbox() *MockSandbox {
	return &MockSandbox{
		Stamps:         make(map[crypto.Hash]int),
		TTLInterval:    500,
		MaxMemoLenght_: 1024,
		FeeFraction_:   0.001,
		MinFee_:        1000,
	}
}
func (m *MockSandbox) HasAccount(addr crypto.Address) bool {
	return false
}
func (m *MockSandbox) Account(addr crypto.Address) *account.Account {
	return nil
}
func (m *MockSandbox) UpdateAccount(acc *account.Account) {}
func (m *MockSandbox) HasValidator(addr crypto.Address) bool {
	return false
}
func (m *MockSandbox) Validator(addr crypto.Address) *validator.Validator {
	return nil
}
func (m *MockSandbox) UpdateValidator(val *validator.Validator) {}
func (m *MockSandbox) CurrentHeight() int {
	return m.CurrentHeight_
}
func (m *MockSandbox) RecentBlockHeight(hash crypto.Hash) int {
	h, ok := m.Stamps[hash]
	if !ok {
		return -1
	}
	return h
}
func (m *MockSandbox) TransactionToLiveInterval() int {
	return m.TTLInterval
}
func (m *MockSandbox) MaxMemoLenght() int {
	return m.MaxMemoLenght_
}
func (m *MockSandbox) FeeFraction() float64 {
	return m.FeeFraction_
}
func (m *MockSandbox) MinFee() int64 {
	return m.MinFee_
}

func (m *MockSandbox) AppendStampAndUpdateHeight(height int, stamp crypto.Hash) {
	m.Stamps[stamp] = height
	m.CurrentHeight_ = height + 1
}
