package store

import (
	"fmt"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

var _ Store = &MockStore{}

type MockStore struct {
	Blocks       map[int]*block.Block
	Accounts     map[crypto.Address]account.Account
	Validators   map[crypto.Address]validator.Validator
	Transactions map[hash.Hash]*tx.Tx
	LastCert     lastInfo
}

func MockingStore() *MockStore {
	return &MockStore{
		Blocks:       make(map[int]*block.Block),
		Accounts:     make(map[crypto.Address]account.Account),
		Validators:   make(map[crypto.Address]validator.Validator),
		Transactions: make(map[hash.Hash]*tx.Tx),
	}
}
func (m *MockStore) Block(hash hash.Hash) (*StoreBlock, error) {
	for h, b := range m.Blocks {
		d, _ := b.Header().Encode()
		if b.Hash().EqualsTo(hash) {
			return &StoreBlock{
				Height:     h,
				Block:      b,
				HeaderData: d,
			}, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockStore) BlockHash(height int) hash.Hash {
	b, ok := m.Blocks[height]
	if ok {
		return b.Hash()
	}
	return hash.UndefHash
}
func (m *MockStore) Transaction(id tx.ID) (*tx.Tx, error) {
	trx, ok := m.Transactions[id]
	if ok {
		return trx, nil
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
func (m *MockStore) Close() error {
	return nil
}

func (m *MockStore) HasAnyBlock() bool {
	return len(m.Blocks) > 0
}

func (m *MockStore) IterateAccounts(consumer func(*account.Account) (stop bool)) {
	for _, a := range m.Accounts {
		acc := a
		stopped := consumer(&acc)
		if stopped {
			return
		}
	}
}

func (m *MockStore) IterateValidators(consumer func(*validator.Validator) (stop bool)) {
	for _, v := range m.Validators {
		val := v
		stopped := consumer(&val)
		if stopped {
			return
		}
	}
}

func (m *MockStore) SaveBlock(height int, b *block.Block, cert *block.Certificate) {
	m.Blocks[height] = b
	for _, trx := range b.Transactions() {
		m.Transactions[trx.ID()] = trx
	}
	m.LastCert.Height = height
	m.LastCert.Cert = cert
}

func (m *MockStore) LastCertificate() (int, *block.Certificate) {
	if m.LastCert.Cert == nil {
		return 0, nil
	}
	return m.LastCert.Height, m.LastCert.Cert
}
func (m *MockStore) BlockHashByStamp(stamp hash.Stamp) hash.Hash {
	for _, b := range m.Blocks {
		if b.Stamp().EqualsTo(stamp) {
			return b.Hash()
		}
	}

	return hash.UndefHash
}
func (m *MockStore) BlockHeightByStamp(stamp hash.Stamp) int {
	for i, b := range m.Blocks {
		if b.Stamp().EqualsTo(stamp) {
			return i
		}
	}

	return 0
}

func (m *MockStore) WriteBatch() error {
	return nil
}
