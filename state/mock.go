package state

import (
	"time"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

type MockState struct {
	Store *store.MockStore
}

func NewMockStore() *MockState {
	return &MockState{
		Store: store.NewMockStore(),
	}
}

func (m *MockState) StoreReader() store.StoreReader {
	return m.Store
}
func (m *MockState) ValidatorSet() validator.ValidatorSetReader {
	return nil
}
func (m *MockState) LastBlockHeight() int {
	return -1
}
func (m *MockState) GenesisHash() crypto.Hash {
	return crypto.UndefHash
}
func (m *MockState) LastBlockHash() crypto.Hash {
	return crypto.UndefHash
}
func (m *MockState) LastBlockTime() time.Time {
	return util.Now()
}
func (m *MockState) LastCommit() *block.Commit {
	return nil
}
func (m *MockState) BlockTime() time.Duration {
	return time.Second
}
func (m *MockState) UpdateLastCommit(blockHash crypto.Hash, commit block.Commit) {

}
func (m *MockState) Fingerprint() string {
	return ""
}
