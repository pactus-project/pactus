package state

import (
	"fmt"
	"time"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var _ StateFacade = &MockState{}

type MockState struct {
	LastBlockCertificate *block.Certificate
	GenHash              crypto.Hash
	Store                *store.MockStore
	TxPool               *txpool.MockTxPool
	InvalidBlockHash     crypto.Hash
	Committee            *committee.Committee
}

func MockingState(committee *committee.Committee) *MockState {
	return &MockState{
		GenHash:   crypto.GenerateTestHash(),
		Store:     store.MockingStore(),
		TxPool:    txpool.MockingTxPool(),
		Committee: committee,
	}
}

func (m *MockState) LastBlockHeight() int {
	return m.Store.LastBlockHeight()
}
func (m *MockState) GenesisHash() crypto.Hash {
	return m.GenHash
}
func (m *MockState) LastBlockHash() crypto.Hash {
	h := m.Store.LastBlockHeight()
	if h > 0 {
		return m.Store.Blocks[m.Store.LastBlockHeight()].Hash()
	}
	return crypto.UndefHash
}
func (m *MockState) LastBlockTime() time.Time {
	return util.Now()
}
func (m *MockState) LastCertificate() *block.Certificate {
	return m.LastBlockCertificate
}
func (m *MockState) BlockTime() time.Duration {
	return time.Second
}
func (m *MockState) UpdateLastCertificate(cert *block.Certificate) error {
	m.LastBlockCertificate = cert
	return nil
}
func (m *MockState) Fingerprint() string {
	return ""
}
func (m *MockState) CommitBlock(height int, b block.Block, cert block.Certificate) error {
	if height != m.LastBlockHeight()+1 {
		return fmt.Errorf("Invalid height")
	}
	if b.Hash().EqualsTo(m.InvalidBlockHash) {
		return fmt.Errorf("Invalid block")
	}
	m.Store.Blocks[height] = &b
	m.LastBlockCertificate = &cert
	return nil
}

func (m *MockState) Close() error {
	return nil
}
func (m *MockState) ProposeBlock(round int) (*block.Block, error) {
	b, _ := block.GenerateTestBlock(nil, nil)
	return b, nil
}
func (m *MockState) ValidateBlock(block block.Block) error {
	return nil
}

func (m *MockState) AddBlock(h int, b *block.Block, trxs []*tx.Tx) {
	m.Store.Blocks[h] = b
	for _, t := range trxs {
		m.Store.Transactions[t.ID()] = &tx.CommittedTx{
			Tx: t, Receipt: t.GenerateReceipt(0, b.Hash()),
		}
	}
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
func (m *MockState) Transaction(id tx.ID) *tx.CommittedTx {
	tx, _ := m.Store.Transaction(id)
	return tx
}
func (m *MockState) Block(height int) *block.Block {
	b, _ := m.Store.Block(height)
	return b
}
func (m *MockState) BlockHeight(hash crypto.Hash) int {
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
func (m *MockState) PendingTx(id tx.ID) *tx.Tx {
	return m.TxPool.PendingTx(id)
}
func (m *MockState) AddPendingTx(trx *tx.Tx) error {
	return m.TxPool.AppendTx(trx)
}
func (m *MockState) AddPendingTxAndBroadcast(trx *tx.Tx) error {
	return m.TxPool.AppendTxAndBroadcast(trx)
}
