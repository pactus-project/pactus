package store

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
)

var _ Store = &MockStore{}

type MockStore struct {
	Blocks       map[uint32]block.Block
	Accounts     map[crypto.Address]account.Account
	Validators   map[crypto.Address]validator.Validator
	Transactions map[tx.ID]tx.Tx
	LastCert     *block.Certificate
	LastHeight   uint32
}

func MockingStore() *MockStore {
	return &MockStore{
		Blocks:       make(map[uint32]block.Block),
		Accounts:     make(map[crypto.Address]account.Account),
		Validators:   make(map[crypto.Address]validator.Validator),
		Transactions: make(map[tx.ID]tx.Tx),
	}
}
func (m *MockStore) Block(hash hash.Hash) (*StoredBlock, error) {
	for h, b := range m.Blocks {
		d, _ := b.Bytes()
		if b.Hash().EqualsTo(hash) {
			return &StoredBlock{
				height: h,
				data:   d,
			}, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockStore) BlockHash(height uint32) hash.Hash {
	b, ok := m.Blocks[height]
	if ok {
		return b.Hash()
	}
	return hash.UndefHash
}
func (m *MockStore) Transaction(id tx.ID) (*tx.Tx, error) {
	trx, ok := m.Transactions[id]
	if ok {
		return &trx, nil
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
func (m *MockStore) TotalAccounts() int32 {
	return int32(len(m.Accounts))
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
func (m *MockStore) ValidatorByNumber(num int32) (*validator.Validator, error) {
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
func (m *MockStore) TotalValidators() int32 {
	return int32(len(m.Validators))
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

func (m *MockStore) SaveBlock(height uint32, b *block.Block, cert *block.Certificate) {
	m.Blocks[height] = *b
	for _, trx := range b.Transactions() {
		m.Transactions[trx.ID()] = *trx
	}
	m.LastHeight = height
	m.LastCert = cert
}

func (m *MockStore) LastCertificate() (uint32, *block.Certificate) {
	if m.LastHeight == 0 {
		return 0, nil
	}
	return m.LastHeight, m.LastCert
}
func (m *MockStore) FindBlockHashByStamp(stamp hash.Stamp) (hash.Hash, bool) {
	if stamp.EqualsTo(hash.UndefHash.Stamp()) {
		return hash.UndefHash, true
	}
	for _, b := range m.Blocks {
		if b.Stamp().EqualsTo(stamp) {
			return b.Hash(), true
		}
	}

	return hash.UndefHash, false
}
func (m *MockStore) FindBlockHeightByStamp(stamp hash.Stamp) (uint32, bool) {
	if stamp.EqualsTo(hash.UndefHash.Stamp()) {
		return 0, true
	}
	for i, b := range m.Blocks {
		if b.Stamp().EqualsTo(stamp) {
			return i, true
		}
	}

	return 0, false
}
func (m *MockStore) WriteBatch() error {
	return nil
}

func (m *MockStore) AddTestValidator() *validator.Validator {
	val, _ := validator.GenerateTestValidator(util.RandInt32(10000))
	val.SubtractFromStake(val.Stake())
	m.UpdateValidator(val)
	return val
}

func (m *MockStore) AddTestAccount() *account.Account {
	acc, _ := account.GenerateTestAccount(util.RandInt32(10000))
	acc.SubtractFromBalance(acc.Balance())
	m.UpdateAccount(acc)
	return acc
}

func (m *MockStore) AddTestBlock(height uint32) *block.Block {
	b := block.GenerateTestBlock(nil, nil)
	cert := block.GenerateTestCertificate(b.Hash())
	m.SaveBlock(height, b, cert)
	return b
}

func (m *MockStore) AddTestTransaction() *tx.Tx {
	tx, _ := tx.GenerateTestSendTx()
	m.Transactions[tx.ID()] = *tx
	return tx
}
func (m *MockStore) RandomTestAcc() *account.Account {
	for _, acc := range m.Accounts {
		return &acc
	}
	panic("no account in sandbox")
}

func (m *MockStore) RandomTestVal() *validator.Validator {
	for _, val := range m.Validators {
		return &val
	}
	panic("no validator in sandbox")
}
