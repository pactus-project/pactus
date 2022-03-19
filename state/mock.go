package state

import (
	"fmt"
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

		m.TestStore.SaveBlock(i+1, b, cert)
	}
}
func (m *MockState) LastBlockHeight() int {
	return m.TestStore.LastCert.Height
}
func (m *MockState) GenesisHash() hash.Hash {
	return m.TestGenHash
}
func (m *MockState) LastBlockHash() hash.Hash {
	if m.TestStore.LastCert.Cert == nil {
		return hash.UndefHash
	}
	return m.TestStore.LastCert.Cert.BlockHash()
}
func (m *MockState) LastBlockTime() time.Time {
	return util.Now()
}
func (m *MockState) LastCertificate() *block.Certificate {
	return m.TestStore.LastCert.Cert
}
func (m *MockState) BlockTime() time.Duration {
	return time.Second
}
func (m *MockState) UpdateLastCertificate(cert *block.Certificate) error {
	m.TestStore.LastCert.Cert = cert
	return nil
}
func (m *MockState) Fingerprint() string {
	return ""
}
func (m *MockState) CommitBlock(h int, b *block.Block, cert *block.Certificate) error {
	if h != m.TestStore.LastCert.Height+1 {
		return fmt.Errorf("invalid height")
	}
	m.TestStore.SaveBlock(h, b, cert)
	return nil
}

func (m *MockState) Close() error {
	return nil
}
func (m *MockState) ProposeBlock(round int) (*block.Block, error) {
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
func (m *MockState) Proposer(round int) *validator.Validator {
	return m.TestCommittee.Proposer(round)
}
func (m *MockState) IsProposer(addr crypto.Address, round int) bool {
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
	tx, _ := m.TestStore.Transaction(id)
	return tx
}
func (m *MockState) Block(hash hash.Hash) *block.Block {
	bi, _ := m.TestStore.Block(hash)
	if bi != nil {
		return bi.Block
	}
	return nil
}
func (m *MockState) BlockHash(height int) hash.Hash {
	return m.TestStore.BlockHash(height)
}
func (m *MockState) Account(addr crypto.Address) *account.Account {
	a, _ := m.TestStore.Account(addr)
	return a
}
func (m *MockState) Validator(addr crypto.Address) *validator.Validator {
	v, _ := m.TestStore.Validator(addr)
	return v
}
func (m *MockState) ValidatorByNumber(n int) *validator.Validator {
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
