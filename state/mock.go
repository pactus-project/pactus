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
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var _ Facade = &MockState{}

type MockState struct {
	GenHash          hash.Hash
	Store            *store.MockStore
	TxPool           *txpool.MockTxPool
	InvalidBlockHash hash.Hash
	Committee        committee.Committee
}

func MockingState() *MockState {
	committee, _ := committee.GenerateTestCommittee(21)
	return &MockState{
		GenHash:   hash.GenerateTestHash(),
		Store:     store.MockingStore(),
		TxPool:    txpool.MockingTxPool(),
		Committee: committee,
	}
}

func (m *MockState) CommitTestBlocks(num int) {
	for i := 0; i < num; i++ {
		b, txs := block.GenerateTestBlock(nil, nil)
		cert := block.GenerateTestCertificate(b.Hash())

		m.Store.SaveBlock(i+1, b, cert)
		for _, tx := range txs {
			m.Store.SaveTransaction(tx)
		}
	}
}
func (m *MockState) LastBlockHeight() int {
	return m.Store.LastBlockHeight()
}
func (m *MockState) GenesisHash() hash.Hash {
	return m.GenHash
}
func (m *MockState) LastBlockHash() hash.Hash {
	h := m.Store.LastBlockHeight()
	if h > 0 {
		b := m.Store.Blocks[m.Store.LastBlockHeight()]
		return b.Hash()
	}
	return hash.UndefHash
}
func (m *MockState) LastBlockTime() time.Time {
	return util.Now()
}
func (m *MockState) LastCertificate() *block.Certificate {
	return m.Store.LastCert.Cert
}
func (m *MockState) BlockTime() time.Duration {
	return time.Second
}
func (m *MockState) UpdateLastCertificate(cert *block.Certificate) error {
	m.Store.LastCert.Cert = cert
	return nil
}
func (m *MockState) Fingerprint() string {
	return ""
}
func (m *MockState) CommitBlock(h int, b *block.Block, cert *block.Certificate) error {
	if h != m.Store.LastBlockHeight()+1 {
		return fmt.Errorf("invalid height")
	}
	if b.Hash().EqualsTo(m.InvalidBlockHash) {
		return fmt.Errorf("invalid block")
	}
	m.Store.SaveBlock(h, b, cert)
	return nil
}

func (m *MockState) Close() error {
	return nil
}
func (m *MockState) ProposeBlock(round int) (*block.Block, error) {
	b, _ := block.GenerateTestBlock(nil, nil)
	return b, nil
}
func (m *MockState) ValidateBlock(block *block.Block) error {
	return nil
}
func (m *MockState) CommitteeValidators() []*validator.Validator {
	return m.Committee.Validators()
}
func (m *MockState) IsInCommittee(addr crypto.Address) bool {
	return m.Committee.Contains(addr)
}
func (m *MockState) Proposer(round int) *validator.Validator {
	return m.Committee.Proposer(round)
}
func (m *MockState) IsProposer(addr crypto.Address, round int) bool {
	return m.Committee.IsProposer(addr, round)
}

func (m *MockState) TotalPower() int64 {

	p := int64(0)
	m.Store.IterateValidators(func(val *validator.Validator) bool {
		p += val.Power()
		return false
	})
	return p
}
func (m *MockState) CommitteePower() int64 {

	return m.Committee.TotalPower()
}
func (m *MockState) Transaction(id tx.ID) *tx.Tx {
	tx, _ := m.Store.Transaction(id)
	return tx
}
func (m *MockState) Block(height int) *block.Block {
	b, _ := m.Store.Block(height)
	return b
}
func (m *MockState) BlockHeight(hash hash.Hash) int {
	h, _ := m.Store.BlockHeight(hash)
	return h
}
func (m *MockState) Account(addr crypto.Address) *account.Account {
	a, _ := m.Store.Account(addr)
	return a
}
func (m *MockState) Validator(addr crypto.Address) *validator.Validator {
	v, _ := m.Store.Validator(addr)
	return v
}
func (m *MockState) ValidatorByNumber(n int) *validator.Validator {
	v, _ := m.Store.ValidatorByNumber(n)
	return v
}
func (m *MockState) PendingTx(id tx.ID) *tx.Tx {
	return m.TxPool.PendingTx(id)
}
func (m *MockState) AddPendingTx(trx *tx.Tx) error {
	if m.TxPool.HasTx(trx.ID()) {
		return errors.Error(errors.ErrGeneric)
	}
	return m.TxPool.AppendTx(trx)
}
func (m *MockState) AddPendingTxAndBroadcast(trx *tx.Tx) error {

	if m.TxPool.HasTx(trx.ID()) {
		return errors.Error(errors.ErrGeneric)
	}
	return m.TxPool.AppendTxAndBroadcast(trx)
}
