package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var _ Facade = &MockState{}

type MockState struct {
	// This locks prevents the Data Race in tests
	lk sync.RWMutex

	TestGenHash   hash.Hash
	TestStore     *store.MockStore
	TestPool      *txpool.MockTxPool
	TestCommittee committee.Committee
	TestParams    param.Params
}

func MockingState() *MockState {
	committee, _ := committee.GenerateTestCommittee(21)
	return &MockState{
		TestGenHash:   hash.GenerateTestHash(),
		TestStore:     store.MockingStore(),
		TestPool:      txpool.MockingTxPool(),
		TestCommittee: committee,
		TestParams:    param.DefaultParams(),
	}
}

func (m *MockState) CommitTestBlocks(num int) {
	for i := 0; i < num; i++ {
		b := block.GenerateTestBlock(nil, nil)
		cert := block.GenerateTestCertificate(b.Hash())

		m.TestStore.SaveBlock(int32(i+1), b, cert)
	}
}
func (m *MockState) LastBlockHeight() int32 {
	m.lk.RLock()
	defer m.lk.RUnlock()

	return m.TestStore.LastHeight
}
func (m *MockState) GenesisHash() hash.Hash {
	return m.TestGenHash
}
func (m *MockState) LastBlockHash() hash.Hash {
	m.lk.RLock()
	defer m.lk.RUnlock()

	return m.TestStore.BlockHash(m.TestStore.LastHeight)
}
func (m *MockState) LastBlockTime() time.Time {
	return util.Now()
}
func (m *MockState) LastCertificate() *block.Certificate {
	m.lk.RLock()
	defer m.lk.RUnlock()

	return m.TestStore.LastCert
}
func (m *MockState) BlockTime() time.Duration {
	return time.Second
}
func (m *MockState) UpdateLastCertificate(cert *block.Certificate) error {
	m.TestStore.LastCert = cert
	return nil
}
func (m *MockState) Fingerprint() string {
	return ""
}
func (m *MockState) CommitBlock(h int32, b *block.Block, cert *block.Certificate) error {
	m.lk.Lock()
	defer m.lk.Unlock()

	if h != m.TestStore.LastHeight+1 {
		return fmt.Errorf("invalid height")
	}
	m.TestStore.SaveBlock(h, b, cert)
	return nil
}

func (m *MockState) Close() error {
	return nil
}
func (m *MockState) ProposeBlock(round int16) (*block.Block, error) {
	b := block.GenerateTestBlock(nil, nil)
	return b, nil
}
func (m *MockState) ValidateBlock(block *block.Block) error {
	return nil
}
func (m *MockState) CommitteeValidators() []*validator.Validator {
	return m.TestCommittee.Validators()
}
func (m *MockState) IsInCommittee(addr crypto.Address) bool {
	return m.TestCommittee.Contains(addr)
}
func (m *MockState) Proposer(round int16) *validator.Validator {
	return m.TestCommittee.Proposer(round)
}
func (m *MockState) IsProposer(addr crypto.Address, round int16) bool {
	return m.TestCommittee.IsProposer(addr, round)
}

func (m *MockState) TotalPower() int64 {

	p := int64(0)
	m.TestStore.IterateValidators(func(val *validator.Validator) bool {
		p += val.Power()
		return false
	})
	return p
}
func (m *MockState) CommitteePower() int64 {

	return m.TestCommittee.TotalPower()
}
func (m *MockState) Transaction(id tx.ID) *tx.Tx {
	m.lk.RLock()
	defer m.lk.RUnlock()

	tx, _ := m.TestStore.Transaction(id)
	return tx
}
func (m *MockState) Block(hash hash.Hash) *block.Block {
	m.lk.RLock()
	defer m.lk.RUnlock()

	bi, _ := m.TestStore.Block(hash)
	if bi != nil {
		b, _ := bi.ToFullBlock()
		return b
	}
	return nil
}
func (m *MockState) BlockHash(height int32) hash.Hash {
	m.lk.RLock()
	defer m.lk.RUnlock()

	return m.TestStore.BlockHash(height)
}
func (m *MockState) AccountByAddress(addr crypto.Address) *account.Account {
	a, _ := m.TestStore.Account(addr)
	return a
}
func (m *MockState) ValidatorByAddress(addr crypto.Address) *validator.Validator {
	v, _ := m.TestStore.Validator(addr)
	return v
}
func (m *MockState) ValidatorByNumber(n int32) *validator.Validator {
	v, _ := m.TestStore.ValidatorByNumber(n)
	return v
}
func (m *MockState) PendingTx(id tx.ID) *tx.Tx {
	return m.TestPool.PendingTx(id)
}
func (m *MockState) AddPendingTx(trx *tx.Tx) error {
	if m.TestPool.HasTx(trx.ID()) {
		return errors.Error(errors.ErrGeneric)
	}
	return m.TestPool.AppendTx(trx)
}
func (m *MockState) AddPendingTxAndBroadcast(trx *tx.Tx) error {

	if m.TestPool.HasTx(trx.ID()) {
		return errors.Error(errors.ErrGeneric)
	}
	return m.TestPool.AppendTxAndBroadcast(trx)
}
func (m *MockState) Params() param.Params {
	return m.TestParams
}
func (m *MockState) MintbaseAddress() crypto.Address {
	return crypto.GenerateTestAddress()
}
func (m *MockState) ValidatorAddress() crypto.Address {
	return crypto.GenerateTestAddress()
}
