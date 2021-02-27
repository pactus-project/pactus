package state

import (
	"fmt"
	"time"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var _ State = &MockState{}

type MockState struct {
	LastBlockCommit  *block.Commit
	GenHash          crypto.Hash
	Store            *store.MockStore
	InvalidBlockHash crypto.Hash
	ValSet           *validator.ValidatorSet
}

func MockingState(valSet *validator.ValidatorSet) *MockState {
	return &MockState{
		GenHash: crypto.GenerateTestHash(),
		Store:   store.MockingStore(),
		ValSet:  valSet,
	}
}

func (m *MockState) StoreReader() store.StoreReader {
	return m.Store
}
func (m *MockState) ValidatorSet() validator.ValidatorSetReader {
	return m.ValSet
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
func (m *MockState) LastCommit() *block.Commit {
	return m.LastBlockCommit
}
func (m *MockState) BlockTime() time.Duration {
	return time.Second
}
func (m *MockState) UpdateLastCommit(commit *block.Commit) error {
	m.LastBlockCommit = commit
	return nil
}
func (m *MockState) Fingerprint() string {
	return ""
}
func (m *MockState) CommitBlock(height int, b block.Block, c block.Commit) error {
	if height != m.LastBlockHeight()+1 {
		return fmt.Errorf("Invalid height")
	}
	if b.Hash().EqualsTo(m.InvalidBlockHash) {
		return fmt.Errorf("Invalid block")
	}
	m.Store.Blocks[height] = &b
	m.LastBlockCommit = &c
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
