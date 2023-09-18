package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
)

var _ Facade = &MockState{}

type MockState struct {
	// This locks prevents the Data Race in tests
	lk sync.RWMutex
	ts *testsuite.TestSuite

	TestGenesis   *genesis.Genesis
	TestStore     *store.MockStore
	TestPool      *txpool.MockTxPool
	TestCommittee committee.Committee
	TestParams    param.Params
}

func MockingState(ts *testsuite.TestSuite) *MockState {
	committee, _ := ts.GenerateTestCommittee(21)
	genDoc := genesis.TestnetGenesis()
	return &MockState{
		ts:            ts,
		TestGenesis:   genDoc,
		TestStore:     store.MockingStore(ts),
		TestPool:      txpool.MockingTxPool(),
		TestCommittee: committee,
		TestParams:    genDoc.Params(),
	}
}

func (m *MockState) CommitTestBlocks(num int) {
	for i := 0; i < num; i++ {
		blk := m.ts.GenerateTestBlock(nil)
		cert := m.ts.GenerateTestCertificate()

		m.TestStore.SaveBlock(m.LastBlockHeight()+1, blk, cert)
	}
}

func (m *MockState) LastBlockHeight() uint32 {
	m.lk.RLock()
	defer m.lk.RUnlock()

	return m.TestStore.LastHeight
}

func (m *MockState) Genesis() *genesis.Genesis {
	return m.TestGenesis
}

func (m *MockState) LastBlockHash() hash.Hash {
	m.lk.RLock()
	defer m.lk.RUnlock()

	return m.TestStore.BlockHash(m.TestStore.LastHeight)
}

func (m *MockState) LastBlockTime() time.Time {
	if len(m.TestStore.Blocks) > 0 {
		return m.TestStore.Blocks[m.TestStore.LastHeight].Header().Time()
	}

	return m.Genesis().GenesisTime()
}

func (m *MockState) LastCertificate() *certificate.Certificate {
	m.lk.RLock()
	defer m.lk.RUnlock()

	return m.TestStore.LastCert
}

func (m *MockState) UpdateLastCertificate(cert *certificate.Certificate) error {
	m.TestStore.LastCert = cert
	return nil
}

func (m *MockState) CommitBlock(h uint32, b *block.Block, cert *certificate.Certificate) error {
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

func (m *MockState) ProposeBlock(_ crypto.Signer, _ crypto.Address, _ int16) (*block.Block, error) {
	b := m.ts.GenerateTestBlock(nil)
	return b, nil
}

func (m *MockState) ValidateBlock(_ *block.Block) error {
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

func (m *MockState) IsValidator(addr crypto.Address) bool {
	return m.TestStore.HasValidator(addr)
}

func (m *MockState) TotalValidators() int32 {
	return m.TestStore.TotalAccounts()
}

func (m *MockState) TotalAccounts() int32 {
	return m.TestStore.TotalAccounts()
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

func (m *MockState) MakeCommittedBlock(data []byte, height uint32, blockHash hash.Hash) *store.CommittedBlock {
	return &store.CommittedBlock{
		Store:     m.TestStore,
		Data:      data,
		BlockHash: blockHash,
		Height:    height,
	}
}

func (m *MockState) CommittedBlock(height uint32) *store.CommittedBlock {
	m.lk.RLock()
	defer m.lk.RUnlock()

	b, _ := m.TestStore.Block(height)
	return b
}

func (m *MockState) CommittedTx(id tx.ID) *store.CommittedTx {
	m.lk.RLock()
	defer m.lk.RUnlock()

	trx, _ := m.TestStore.Transaction(id)
	return trx
}

func (m *MockState) BlockHash(height uint32) hash.Hash {
	m.lk.RLock()
	defer m.lk.RUnlock()

	return m.TestStore.BlockHash(height)
}

func (m *MockState) BlockHeight(hash hash.Hash) uint32 {
	m.lk.RLock()
	defer m.lk.RUnlock()

	return m.TestStore.BlockHeight(hash)
}

func (m *MockState) AccountByAddress(addr crypto.Address) *account.Account {
	a, _ := m.TestStore.Account(addr)
	return a
}

func (m *MockState) AccountByNumber(number int32) *account.Account {
	a, _ := m.TestStore.AccountByNumber(number)
	return a
}

func (m *MockState) ValidatorAddresses() []crypto.Address {
	return m.TestStore.ValidatorAddresses()
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

func (m *MockState) CalculateFee(_ int64, payloadType payload.Type) (int64, error) {
	switch payloadType {
	case payload.TypeTransfer,
		payload.TypeBond,
		payload.TypeWithdraw:
		{
			return m.ts.RandInt64(1e9), nil
		}

	case payload.TypeUnbond,
		payload.TypeSortition:
		{
			return 0, nil
		}

	default:
		return 0, errors.Errorf(errors.ErrInvalidTx, "unexpected tx type: %v", payloadType)
	}
}
