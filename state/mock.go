package state

import (
	"sync"
	"time"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
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
	TestValKeys   []*bls.ValidatorKey
	TestParams    *param.Params
}

func MockingState(ts *testsuite.TestSuite) *MockState {
	cmt, valKeys := ts.GenerateTestCommittee(21)
	genDoc := genesis.MainnetGenesis()

	return &MockState{
		ts:            ts,
		TestGenesis:   genDoc,
		TestStore:     store.MockingStore(ts),
		TestPool:      txpool.MockingTxPool(),
		TestCommittee: cmt,
		TestValKeys:   valKeys,
		TestParams:    param.FromGenesis(genDoc.Params()),
	}
}

func (m *MockState) CommitTestBlocks(num int) {
	for i := 0; i < num; i++ {
		height := m.LastBlockHeight() + 1
		blk, cert := m.ts.GenerateTestBlock(height)

		m.TestStore.SaveBlock(blk, cert)
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
	m.lk.RLock()
	defer m.lk.RUnlock()

	if len(m.TestStore.Blocks) > 0 {
		return m.TestStore.Blocks[m.TestStore.LastHeight].Header().Time()
	}

	return m.Genesis().GenesisTime()
}

func (m *MockState) LastCertificate() *certificate.BlockCertificate {
	m.lk.RLock()
	defer m.lk.RUnlock()

	return m.TestStore.LastCert
}

func (*MockState) UpdateLastCertificate(_ *vote.Vote) error {
	return nil
}

func (m *MockState) CommitBlock(blk *block.Block, cert *certificate.BlockCertificate) error {
	m.lk.Lock()
	defer m.lk.Unlock()

	m.TestStore.SaveBlock(blk, cert)

	return nil
}

func (*MockState) Close() {}

func (m *MockState) ProposeBlock(valKey *bls.ValidatorKey, _ crypto.Address) (*block.Block, error) {
	blk, _ := m.ts.GenerateTestBlock(m.TestStore.LastHeight, testsuite.BlockWithProposer(valKey.Address()))

	return blk, nil
}

func (*MockState) ValidateBlock(_ *block.Block, _ int16) error {
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
	power := int64(0)
	m.TestStore.IterateValidators(func(val *validator.Validator) bool {
		power += val.Power()

		return false
	})

	return power
}

func (m *MockState) CommitteePower() int64 {
	return m.TestCommittee.TotalPower()
}

func (m *MockState) CommittedBlock(height uint32) *store.CommittedBlock {
	m.lk.RLock()
	defer m.lk.RUnlock()

	b, _ := m.TestStore.Block(height)

	return b
}

func (m *MockState) CommittedTx(txID tx.ID) *store.CommittedTx {
	m.lk.RLock()
	defer m.lk.RUnlock()

	trx, _ := m.TestStore.Transaction(txID)

	return trx
}

func (m *MockState) BlockHash(height uint32) hash.Hash {
	m.lk.RLock()
	defer m.lk.RUnlock()

	return m.TestStore.BlockHash(height)
}

func (m *MockState) BlockHeight(h hash.Hash) uint32 {
	m.lk.RLock()
	defer m.lk.RUnlock()

	return m.TestStore.BlockHeight(h)
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

func (m *MockState) PendingTx(txID tx.ID) *tx.Tx {
	return m.TestPool.PendingTx(txID)
}

func (m *MockState) AddPendingTx(trx *tx.Tx) error {
	return m.TestPool.AppendTx(trx)
}

func (m *MockState) AddPendingTxAndBroadcast(trx *tx.Tx) error {
	return m.TestPool.AppendTxAndBroadcast(trx)
}

func (m *MockState) Params() *param.Params {
	return m.TestParams
}

func (m *MockState) CalculateFee(amt amount.Amount, payloadType payload.Type) amount.Amount {
	return m.TestPool.EstimatedFee(amt, payloadType)
}

func (m *MockState) PublicKey(addr crypto.Address) (crypto.PublicKey, error) {
	return m.TestStore.PublicKey(addr)
}

func (*MockState) AvailabilityScore(_ int32) float64 {
	return 0.987
}

func (m *MockState) AllPendingTxs() []*tx.Tx {
	return m.TestPool.Txs
}

func (m *MockState) IsPruned() bool {
	return m.TestStore.IsPruned()
}

func (m *MockState) PruningHeight() uint32 {
	return m.TestStore.PruningHeight()
}
