package store

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
)

var _ Store = &MockStore{}

type MockStore struct {
	ts *testsuite.TestSuite

	Blocks     map[uint32]*block.Block
	Accounts   map[crypto.Address]*account.Account
	Validators map[crypto.Address]*validator.Validator
	LastCert   *certificate.BlockCertificate
	LastHeight uint32
	TestConfig *Config
}

func MockingStore(ts *testsuite.TestSuite) *MockStore {
	return &MockStore{
		ts:         ts,
		Blocks:     make(map[uint32]*block.Block),
		Accounts:   make(map[crypto.Address]*account.Account),
		Validators: make(map[crypto.Address]*validator.Validator),
		TestConfig: DefaultConfig(),
	}
}

func (m *MockStore) Block(height uint32) (*CommittedBlock, error) {
	b, ok := m.Blocks[height]
	if ok {
		d, _ := b.Bytes()

		return &CommittedBlock{
			BlockHash: b.Hash(),
			Height:    height,
			Data:      d,
		}, nil
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

func (m *MockStore) BlockHeight(h hash.Hash) uint32 {
	for height, b := range m.Blocks {
		if b.Hash() == h {
			return height
		}
	}

	return 0
}

func (m *MockStore) SortitionSeed(blockHeight uint32) *sortition.VerifiableSeed {
	if blk, ok := m.Blocks[blockHeight]; ok {
		sortitionSeed := blk.Header().SortitionSeed()

		return &sortitionSeed
	}

	return nil
}

func (m *MockStore) PublicKey(addr crypto.Address) (crypto.PublicKey, error) {
	for _, blk := range m.Blocks {
		for _, trx := range blk.Transactions() {
			if trx.Payload().Signer() == addr {
				return trx.PublicKey(), nil
			}
		}
	}
	for _, val := range m.Validators {
		if val.Address() == addr {
			return val.PublicKey(), nil
		}
	}

	return nil, ErrNotFound
}

func (m *MockStore) HasPublicKey(addr crypto.Address) bool {
	pub, _ := m.PublicKey(addr)

	return pub != nil
}

func (m *MockStore) Transaction(txID tx.ID) (*CommittedTx, error) {
	for height, blk := range m.Blocks {
		for _, trx := range blk.Transactions() {
			if trx.ID() == txID {
				data, _ := trx.Bytes()

				return &CommittedTx{
					TxID:      txID,
					Height:    height,
					BlockTime: blk.Header().UnixTime(),
					Data:      data,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("not found")
}

func (m *MockStore) RecentTransaction(txID tx.ID) bool {
	for _, blk := range m.Blocks {
		for _, trx := range blk.Transactions() {
			if trx.ID() == txID {
				return true
			}
		}
	}

	return false
}

func (m *MockStore) HasAccount(addr crypto.Address) bool {
	_, ok := m.Accounts[addr]

	return ok
}

func (m *MockStore) Account(addr crypto.Address) (*account.Account, error) {
	a, ok := m.Accounts[addr]
	if ok {
		return a.Clone(), nil
	}

	return nil, fmt.Errorf("not found")
}

func (m *MockStore) AccountByNumber(number int32) (*account.Account, error) {
	for _, v := range m.Accounts {
		if v.Number() == number {
			return v.Clone(), nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (m *MockStore) UpdateAccount(addr crypto.Address, acc *account.Account) {
	m.Accounts[addr] = acc
}

func (m *MockStore) TotalAccounts() int32 {
	return int32(len(m.Accounts))
}

func (m *MockStore) HasValidator(addr crypto.Address) bool {
	_, ok := m.Validators[addr]

	return ok
}

func (m *MockStore) ValidatorAddresses() []crypto.Address {
	addrs := make([]crypto.Address, 0, len(m.Validators))
	for addr := range m.Validators {
		addrs = append(addrs, addr)
	}

	return addrs
}

func (m *MockStore) Validator(addr crypto.Address) (*validator.Validator, error) {
	v, ok := m.Validators[addr]
	if ok {
		return v.Clone(), nil
	}

	return nil, ErrNotFound
}

func (m *MockStore) ValidatorByNumber(num int32) (*validator.Validator, error) {
	for _, v := range m.Validators {
		if v.Number() == num {
			return v.Clone(), nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (m *MockStore) UpdateValidator(val *validator.Validator) {
	m.Validators[val.Address()] = val
}

func (m *MockStore) TotalValidators() int32 {
	return int32(len(m.Validators))
}

func (*MockStore) Close() {}

func (m *MockStore) HasAnyBlock() bool {
	return len(m.Blocks) > 0
}

func (m *MockStore) IterateAccounts(consumer func(crypto.Address, *account.Account) (stop bool)) {
	for addr, acc := range m.Accounts {
		stopped := consumer(addr, acc.Clone())
		if stopped {
			return
		}
	}
}

func (m *MockStore) IterateValidators(consumer func(*validator.Validator) (stop bool)) {
	for _, val := range m.Validators {
		stopped := consumer(val.Clone())
		if stopped {
			return
		}
	}
}

func (m *MockStore) SaveBlock(blk *block.Block, cert *certificate.BlockCertificate) {
	m.Blocks[cert.Height()] = blk
	m.LastHeight = cert.Height()
	m.LastCert = cert
}

func (m *MockStore) LastCertificate() *certificate.BlockCertificate {
	if m.LastHeight == 0 {
		return nil
	}

	return m.LastCert
}

func (*MockStore) WriteBatch() error {
	return nil
}

func (m *MockStore) AddTestValidator() *validator.Validator {
	val := m.ts.GenerateTestValidator()
	m.UpdateValidator(val)

	return val
}

func (m *MockStore) AddTestAccount() (*account.Account, crypto.Address) {
	acc, addr := m.ts.GenerateTestAccount()
	m.UpdateAccount(addr, acc)

	return acc, addr
}

func (m *MockStore) AddTestBlock(height uint32) *block.Block {
	blk, cert := m.ts.GenerateTestBlock(height)
	m.SaveBlock(blk, cert)

	return blk
}

func (m *MockStore) RandomTestAcc() (crypto.Address, *account.Account) {
	for addr, acc := range m.Accounts {
		// Do not return the Treasury address for tests,
		// as it may cause some tests to randomly fail.
		if addr == crypto.TreasuryAddress {
			continue
		}

		return addr, acc
	}
	panic("no account in sandbox")
}

func (m *MockStore) RandomTestVal() *validator.Validator {
	for _, val := range m.Validators {
		return val
	}
	panic("no validator in sandbox")
}

func (m *MockStore) IsBanned(addr crypto.Address) bool {
	return m.TestConfig.BannedAddrs[addr]
}

func (*MockStore) Prune(_ func(_ bool, _ uint32) bool) error {
	return nil
}

func (*MockStore) IsPruned() bool {
	return false
}

func (*MockStore) PruningHeight() uint32 {
	return 0
}

func (m *MockStore) XeggexAccount() *XeggexAccount {
	return &m.TestConfig.XeggexAccount
}
