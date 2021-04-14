package store

import (
	"fmt"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

type MockStore struct {
	Blocks       map[int]block.Block
	Accounts     map[crypto.Address]account.Account
	Validators   map[crypto.Address]validator.Validator
	Transactions map[crypto.Hash]tx.CommittedTx
}

func MockingStore() *MockStore {
	return &MockStore{
		Blocks:       make(map[int]block.Block),
		Accounts:     make(map[crypto.Address]account.Account),
		Validators:   make(map[crypto.Address]validator.Validator),
		Transactions: make(map[crypto.Hash]tx.CommittedTx),
	}
}
func (m *MockStore) Block(height int) (*block.Block, error) {
	b, ok := m.Blocks[height]
	if ok {
		return &b, nil
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockStore) BlockHeight(hash crypto.Hash) (int, error) {
	for i, b := range m.Blocks {
		if b.Hash().EqualsTo(hash) {
			return i, nil
		}
	}
	return -1, nil
}
func (m *MockStore) Transaction(id tx.ID) (*tx.CommittedTx, error) {
	b, ok := m.Transactions[id]
	if ok {
		return &b, nil
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockStore) HasAccount(addr crypto.Address) bool {
	_, ok := m.Accounts[addr]
	return ok
}
func (m *MockStore) Account(addr crypto.Address) (*account.Account, error) {
	a, ok := m.Accounts[addr]
	if ok {
		return &a, nil
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockStore) UpdateAccount(acc *account.Account) {
	m.Accounts[acc.Address()] = *acc
}
func (m *MockStore) TotalAccounts() int {
	return len(m.Accounts)
}
func (m *MockStore) HasValidator(addr crypto.Address) bool {
	_, ok := m.Validators[addr]
	return ok
}
func (m *MockStore) Validator(addr crypto.Address) (*validator.Validator, error) {
	v, ok := m.Validators[addr]
	if ok {
		return &v, nil
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockStore) ValidatorByNumber(num int) (*validator.Validator, error) {
	for _, v := range m.Validators {
		if v.Number() == num {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockStore) UpdateValidator(val *validator.Validator) {
	m.Validators[val.Address()] = *val
}
func (m *MockStore) TotalValidators() int {
	return len(m.Validators)
}
func (m *MockStore) LastBlockHeight() int {
	return len(m.Blocks)
}

func (m *MockStore) Close() error {
	m.Blocks = make(map[int]block.Block)
	m.Accounts = make(map[crypto.Address]account.Account)
	m.Validators = make(map[crypto.Address]validator.Validator)
	m.Transactions = make(map[crypto.Hash]tx.CommittedTx)
	return nil
}

func (m *MockStore) HasAnyBlock() bool {
	return len(m.Blocks) > 0
}

func (m *MockStore) IterateAccounts(consumer func(*account.Account) (stop bool)) {
	for _, a := range m.Accounts {
		stopped := consumer(&a)
		if stopped {
			return
		}
	}
}

func (m *MockStore) IterateValidators(consumer func(*validator.Validator) (stop bool)) {
	for _, v := range m.Validators {
		stopped := consumer(&v)
		if stopped {
			return
		}
	}
}

func (m *MockStore) SaveBlock(height int, block *block.Block) {
	m.Blocks[height] = *block
}

func (m *MockStore) SaveTransaction(ctrx *tx.CommittedTx) {
	m.Transactions[ctrx.Tx.ID()] = *ctrx
}

func (m *MockStore) WriteBatch() error {
	return nil
}
